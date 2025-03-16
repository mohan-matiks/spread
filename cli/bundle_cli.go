package cli

import (
	"bytes"
	"crypto/sha256"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"mime/multipart"
	"net/http"
	"net/url"
	"os"
	"os/exec"
	"runtime"
	"strings"

	"github.com/SwishHQ/spread/types"
	"github.com/SwishHQ/spread/utils"
	"github.com/google/uuid"
)

// Configuration struct to hold bundle parameters
type BundleConfig struct {
	RemoteURL           string
	AuthKey             string
	AppName             string
	Environment         string
	Description         string
	TargetVersion       string
	ProjectDir          string
	OSName              string
	IsTypescriptProject bool
	DisableMinify       bool
	Hermes              bool
}

// PushBundle uploads a new bundle to the server
func PushBundle(config BundleConfig) error {
	if err := validateConfig(config); err != nil {
		return fmt.Errorf("invalid configuration: %w", err)
	}

	if err := prepareBuildDirectory(config.ProjectDir); err != nil {
		return fmt.Errorf("failed to prepare build directory: %w", err)
	}

	if err := buildBundle(config); err != nil {
		return fmt.Errorf("failed to build bundle: %w", err)
	}

	if config.Hermes {
		if err := processHermesBundle(config); err != nil {
			return fmt.Errorf("failed to process hermes bundle: %w", err)
		}
	}

	hash, err := getHash(config.ProjectDir + "build")
	if err != nil {
		return fmt.Errorf("failed to generate hash: %w", err)
	}

	fileName := uuid.New().String() + ".zip"
	if err := createAndUploadBundle(config, fileName, hash); err != nil {
		return fmt.Errorf("failed to create and upload bundle: %w", err)
	}

	return nil
}

func validateConfig(config BundleConfig) error {
	if config.TargetVersion == "" || config.AppName == "" || config.Environment == "" {
		return fmt.Errorf("missing required fields: target version, app name, or environment")
	}
	return nil
}

func prepareBuildDirectory(projectDir string) error {
	buildPath := projectDir + "build"
	if exist, _ := utils.PathExists(buildPath); exist {
		if err := os.RemoveAll(buildPath); err != nil {
			return err
		}
	}
	return os.MkdirAll(projectDir+"build/CodePush", os.ModePerm)
}

func buildBundle(config BundleConfig) error {
	jsName := "main.jsbundle"
	if config.OSName == "android" {
		jsName = "index.android.bundle"
	}

	indexFile := "index.js"
	if config.IsTypescriptProject {
		indexFile = "index.tsx"
	}

	minify := "true"
	if config.DisableMinify {
		minify = "false"
	}

	buildPath := config.ProjectDir + "build/CodePush"
	bundleURL := buildPath + "/" + jsName

	return executeCommand(createBundleCommand(config, buildPath, bundleURL, indexFile, minify))
}

func getHash(path string) (string, error) {
	// First, we need to find all the files in the given path.
	files, err := getAllFiles(path)
	if err != nil {
		return "", err
	}
	// We create a list to hold the names and special codes (hashes) of the files.
	var fileHash []string
	// Now, we go through each file we found.
	for _, table1 := range files {
		fb, _ := os.ReadFile(table1)
		// We use a special tool (sha256) to create a unique code (hash) for the file content.
		h := sha256.New()
		h.Write(fb)
		hash := h.Sum(nil)
		fileName := strings.TrimPrefix(table1, path)
		fileHash = append(fileHash, fileName+":"+fmt.Sprintf("%x", hash))
	}
	j, _ := json.Marshal(fileHash)
	jStr := strings.ReplaceAll(string(j), "\\/", "/")
	h := sha256.New()
	h.Write([]byte(jStr))
	hash := h.Sum(nil)
	return fmt.Sprintf("%x", hash), nil
}

// This function helps us find all the files in a given directory and its subdirectories.
func getAllFiles(dirPth string) (files []string, err error) {
	var dirs []string
	dir, err := os.ReadDir(dirPth)
	if err != nil {
		return nil, err
	}
	pathSeparator := string(os.PathSeparator)

	// Now, we go through each item in the directory.
	for _, fi := range dir {
		// If the item is a directory, we add it to our list of directories to check later.
		if fi.IsDir() {
			dirs = append(dirs, dirPth+pathSeparator+fi.Name())
			// We also call this function again to check inside this subdirectory.
			temp, _ := getAllFiles(dirPth + pathSeparator + fi.Name())
			files = append(files, temp...)
		} else {
			// If the item is a file, we add its path to our list of files.
			files = append(files, dirPth+pathSeparator+fi.Name())
		}
	}

	// Now, we go through each subdirectory we found and get all the files inside them.
	for _, table := range dirs {
		temp, _ := getAllFiles(table)
		files = append(files, temp...)
	}
	return files, nil
}

// This function creates a new HTTP request for file upload.
// It takes in the URI of the server, the authentication key, parameters, the name of the parameter, and the path of the file.
// It returns the HTTP request and any error that occurred.
func newfileUploadRequest(uri string, authKey string, params map[string]string, paramName, path string) (*http.Request, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)
	part, err := writer.CreateFormFile(paramName, path)
	if err != nil {
		return nil, err
	}
	_, err = io.Copy(part, file)
	if err != nil {
		return nil, err
	}
	for key, val := range params {
		_ = writer.WriteField(key, val)
	}
	err = writer.Close()
	if err != nil {
		return nil, err
	}
	request, err := http.NewRequest("POST", uri, body)
	request.Header.Set("Content-Type", writer.FormDataContentType())
	if authKey != "" {
		request.Header.Set("x-auth-key", authKey)
	}
	return request, err
}

func executeCommand(cmd *exec.Cmd) error {
	stdout, err := cmd.StdoutPipe()
	if err != nil {
		log.Panic("cmd.StdoutPipe() failed with ", err)
	}
	stderr, err := cmd.StderrPipe()
	if err != nil {
		log.Panic("cmd.StderrPipe() failed with ", err)
	}
	if err := cmd.Start(); err != nil {
		log.Panic("cmd.Start() failed with ", err)
	}

	go func() {
		if _, err := io.Copy(os.Stdout, stdout); err != nil {
			log.Panic("failed to copy stdout: ", err)
		}
	}()
	go func() {
		if _, err := io.Copy(os.Stderr, stderr); err != nil {
			log.Panic("failed to copy stderr: ", err)
		}
	}()
	if err := cmd.Wait(); err != nil {
		log.Panic("cmd.Run() failed with ", err)
	}
	return nil
}

func createBundleCommand(config BundleConfig, buildPath, bundleURL, indexFile, minify string) *exec.Cmd {
	cmd := exec.Command(
		"npx",
		"react-native",
		"bundle",
		"--assets-dest",
		buildPath,
		"--bundle-output",
		bundleURL,
		"--dev",
		"false",
		"--entry-file",
		indexFile,
		"--platform",
		config.OSName,
		"--minify",
		minify)
	cmd.Dir = config.ProjectDir
	return cmd
}

func processHermesBundle(config BundleConfig) error {
	sysType := runtime.GOOS
	exc := "/osx-bin/hermesc"
	if sysType == "linux" {
		exc = "/linux64-bin/hermesc"
	}
	if sysType == "windows" {
		exc = "/win64-bin/hermesc.exe"
	}
	hbcUrl := config.ProjectDir + "build/CodePush/" + "main.jsbundle" + ".hbc"
	cmd := exec.Command(
		config.ProjectDir+"node_modules/react-native/sdks/hermesc"+exc,
		"-emit-binary",
		"-out",
		hbcUrl,
		config.ProjectDir+"build/CodePush/"+"main.jsbundle",
		// "-output-source-map",
	)

	cmd.Dir = config.ProjectDir
	out, err := cmd.CombinedOutput()
	if err != nil {
		fmt.Printf("combined out:\n%s\n", string(out))
		log.Panic("cmd.Run() failed with ", err)
	}
	err = os.Remove(config.ProjectDir + "build/CodePush/" + "main.jsbundle")
	if err != nil {
		panic(err.Error())
	}
	data, err := os.ReadFile(hbcUrl)
	if err != nil {
		panic(err.Error())
	}
	err = os.WriteFile(config.ProjectDir+"build/CodePush/"+"main.jsbundle", data, 0755)
	if err != nil {
		panic(err.Error())
	}
	err = os.Remove(hbcUrl)
	if err != nil {
		panic(err.Error())
	}
	return nil
}

func createAndUploadBundle(config BundleConfig, fileName string, hash string) error {
	log.Println("✦ Zipping bundle"+" "+config.AppName+" "+config.TargetVersion, fileName)
	utils.Zip(config.ProjectDir+"build", fileName)

	// exec.Command("open", config.ProjectDir+"build").Run()
	os.RemoveAll(config.ProjectDir + "build")
	log.Println("✦ Uploading bundle")

	Url, err := url.Parse(config.RemoteURL + "/bundle/upload")
	if err != nil {
		log.Panic(err.Error())
	}
	pathName := fileName
	req, err := newfileUploadRequest(Url.String(), config.AuthKey, map[string]string{"filename": fileName}, "file", pathName)
	if err != nil {
		log.Panic(err.Error())
	}
	uploadClient := &http.Client{}
	resp, err := uploadClient.Do(req)
	if err != nil {
		log.Panic(err.Error())
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Panicf("Failed to read response body: %v", err)
	}
	if resp.StatusCode != 200 {
		log.Println("✦ Upload fail", string(body))
		return fmt.Errorf("upload failed: %s", resp.Status)
	}
	log.Println("✦ Bundle has been uploaded successfully.")
	log.Println("✦ Creating a new bundle")

	Url, err = url.Parse(config.RemoteURL + "/bundle/create")
	if err != nil {
		log.Panic("Server url error :", err)
	}
	fileInfo, _ := os.Stat(fileName)
	size := fileInfo.Size()
	key := uuid.New().String() + ".zip"
	createBundleReq := types.CreateNewBundleRequest{
		AppName:      config.AppName,
		Environment:  config.Environment,
		DownloadFile: key,
		Description:  config.Description,
		AppVersion:   config.TargetVersion,
		Size:         size,
		Hash:         hash,
	}
	jsonByte, _ := json.Marshal(createBundleReq)
	req, _ = http.NewRequest("POST", Url.String(), bytes.NewBuffer(jsonByte))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("x-auth-key", config.AuthKey)
	createBundleClient := &http.Client{}
	resp, err = createBundleClient.Do(req)
	if err != nil {
		log.Panic(err.Error())
	}
	defer resp.Body.Close()
	if resp.StatusCode != 200 {
		body, err := io.ReadAll(resp.Body)
		if err != nil {
			log.Panicf("Failed to read error response: %v", err)
		}
		log.Printf("✦ Create bundle failed. Status: %s, Response: %s", resp.Status, string(body))
		return fmt.Errorf("failed to create bundle: %s", resp.Status)
	}
	log.Println("✦ Bundle has been created successfully.")
	os.RemoveAll(fileName)
	return nil
}
