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

func PushBundle(
	remoteUrl string,
	authKey string,
	appName string,
	environment string,
	description string,
	targetVersion string,
	projectDir string,
	osName string, // ios or android
	isTypescriptProject bool,
	disableMinify bool,
	hermes bool,
) {

	if targetVersion == "" || appName == "" || environment == "" {
		fmt.Println("Usage: spread bundle push --target-version <TargetVersion> --app-name <AppName> --environment <environment> --project-dir <*Optional React native project dir> --os-name <OSName> --description <*Optional Description> --is-typescript (*Optional) --disable-minify (*Optional) --hermes (*Optional)")
		return
	}

	exist, _ := utils.PathExists(projectDir + "build")
	// check if the build folder exists
	// if it exists, remove it
	if exist {
		os.RemoveAll(projectDir + "build")
	}

	// create the build folder in the given bundle path
	if err := os.MkdirAll(projectDir+"build/CodePush", os.ModePerm); err != nil {
		fmt.Println("Create folder error :" + err.Error())
		return
	}

	// create the CodePush folder in the build folder
	jsName := "main.jsbundle"
	if osName == "android" {
		jsName = "index.android.bundle"
	}

	minify := "true"
	if disableMinify {
		minify = "false"
	}

	indexFile := "index.js"
	if isTypescriptProject {
		indexFile = "index.tsx"
	}

	buildUrl := projectDir + "build/CodePush"

	bundelUrl := buildUrl + "/" + jsName
	cmd := exec.Command(
		"npx",
		"react-native",
		"bundle",
		"--assets-dest",
		buildUrl,
		"--bundle-output",
		bundelUrl,
		"--dev",
		"false",
		"--entry-file",
		indexFile,
		"--platform",
		osName,
		"--minify",
		minify)
	cmd.Dir = projectDir
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

	// if hermes is true, then we need to build the hbc file
	// if hermes is false, then we need to build the js bundle file
	// how to build the hbc file:
	// 1. get the hermesc binary path
	// 2. build the hbc file
	// 3. replace the js bundle file with the hbc file
	if hermes {
		sysType := runtime.GOOS
		exc := "/osx-bin/hermesc"
		if sysType == "linux" {
			exc = "/linux64-bin/hermesc"
		}
		if sysType == "windows" {
			exc = "/win64-bin/hermesc.exe"
		}
		hbcUrl := projectDir + "build/CodePush/" + jsName + ".hbc"
		cmd := exec.Command(
			projectDir+"node_modules/react-native/sdks/hermesc"+exc,
			"-emit-binary",
			"-out",
			hbcUrl,
			bundelUrl,
			// "-output-source-map",
		)

		cmd.Dir = projectDir
		out, err := cmd.CombinedOutput()
		if err != nil {
			fmt.Printf("combined out:\n%s\n", string(out))
			log.Panic("cmd.Run() failed with ", err)
		}
		err = os.Remove(bundelUrl)
		if err != nil {
			panic(err.Error())
		}
		data, err := os.ReadFile(hbcUrl)
		if err != nil {
			panic(err.Error())
		}
		err = os.WriteFile(bundelUrl, data, 0755)
		if err != nil {
			panic(err.Error())
		}
		err = os.Remove(hbcUrl)
		if err != nil {
			panic(err.Error())
		}
	}

	hash, error := getHash(projectDir + "build")
	if error != nil {
		log.Panic("hash error", error.Error())
	}
	log.Println("✦ Hash: ", hash)
	uuidStr, _ := uuid.NewUUID()
	fileName := uuidStr.String() + ".zip"
	log.Println("✦ Zipping bundle", fileName)
	utils.Zip(projectDir+"build", fileName)

	// exec.Command("open", projectDir+"build").Run()
	os.RemoveAll(projectDir + "build")
	log.Println("✦ Uploading bundle")

	Url, err := url.Parse(remoteUrl + "/bundle/upload")
	if err != nil {
		log.Panic(err.Error())
	}
	pathName := fileName
	req, err := newfileUploadRequest(Url.String(), authKey, map[string]string{"filename": fileName}, "file", pathName)
	if err != nil {
		log.Panic(err.Error())
	}
	uploadClient := &http.Client{}
	resp, err := uploadClient.Do(req)
	if err != nil {
		log.Panic(err.Error())
	}
	defer resp.Body.Close()
	if resp.StatusCode != 200 {
		log.Println("✦ Upload fail", resp)
		return
	}

	log.Println("✦ Bundle has been uploaded successfully.")
	log.Println("✦ Creating a new bundle")

	Url, err = url.Parse(remoteUrl + "/bundle/create")
	if err != nil {
		log.Panic("Server url error :", err)
	}
	fileInfo, _ := os.Stat(fileName)
	size := fileInfo.Size()
	key := uuidStr.String() + ".zip"
	createBundleReq := types.CreateNewBundleRequest{
		AppName:      appName,
		Environment:  environment,
		DownloadFile: key,
		Description:  description,
		AppVersion:   targetVersion,
		Size:         size,
		Hash:         hash,
	}
	jsonByte, _ := json.Marshal(createBundleReq)
	req, _ = http.NewRequest("POST", Url.String(), bytes.NewBuffer(jsonByte))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("x-auth-key", authKey)
	createBundleClient := &http.Client{}
	resp, err = createBundleClient.Do(req)
	if err != nil {
		log.Panic(err.Error())
	}
	defer resp.Body.Close()
	if resp.StatusCode != 200 {
		log.Println("✦ Create bundle fail", resp)
		return
	}
	log.Println("✦ Bundle has been created successfully.")
	os.RemoveAll(fileName)
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
	PthSep := string(os.PathSeparator)

	// Now, we go through each item in the directory.
	for _, fi := range dir {
		// If the item is a directory, we add it to our list of directories to check later.
		if fi.IsDir() {
			dirs = append(dirs, dirPth+PthSep+fi.Name())
			// We also call this function again to check inside this subdirectory.
			getAllFiles(dirPth + PthSep + fi.Name())
		} else {
			// If the item is a file, we add its path to our list of files.
			files = append(files, dirPth+PthSep+fi.Name())
		}
	}

	// Now, we go through each subdirectory we found and get all the files inside them.
	for _, table := range dirs {
		temp, _ := getAllFiles(table)
		for _, temp1 := range temp {
			files = append(files, temp1)
		}
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
