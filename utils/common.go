package utils

import (
	"archive/zip"
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"io"
	"log"
	"math/rand"
	"os"
	"path/filepath"
	"sort"
	"strconv"
	"strings"
	"time"
)

// generate a random auth key
// how it works:
// 1. seed the random number generator
// 2. create a character array
// 3. create a byte array
// 4. for each element in the byte array, generate a random character
// 5. return the byte array as a string
func GenerateAuthKey() string {
	rand.Seed(time.Now().UnixNano())
	chars := "ABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	b := make([]byte, 32)
	for i := range b {
		b[i] = chars[rand.Intn(len(chars))]
	}
	return string(b)
}

// Zip function creates a zip file from a given source directory.
func Zip(src_dir string, zip_file_name string) {
	prefix := string(os.PathSeparator) // Define the path separator

	// Remove the zip file if it already exists
	os.RemoveAll(zip_file_name)
	// Create a new zip file
	zipfile, _ := os.Create(zip_file_name)
	defer zipfile.Close() // Close the file when done

	// Create a new zip writer
	archive := zip.NewWriter(zipfile)
	defer archive.Close() // Close the zip writer when done

	// Normalize the source directory path
	nowSrc := strings.ReplaceAll(strings.Replace(src_dir, "./", "", 1), "\\\\", prefix)

	// Walk through the source directory
	filepath.Walk(src_dir, func(path string, info os.FileInfo, _ error) error {
		// Skip the source directory itself
		if path == src_dir {
			return nil
		}

		// Create a zip file header
		header, _ := zip.FileInfoHeader(info)
		// Set the header name to the relative path within the zip
		header.Name = strings.TrimPrefix(path, nowSrc+prefix)
		// If the file is a directory, add a path separator to the end
		if info.IsDir() {
			header.Name += prefix
		} else {
			// For files, use the Deflate compression method
			header.Method = zip.Deflate
		}

		// Create a writer for the zip file
		writer, _ := archive.CreateHeader(header)
		// If the file is not a directory, copy its content to the zip
		if !info.IsDir() {
			file, _ := os.Open(path)
			defer file.Close()    // Close the file when done
			io.Copy(writer, file) // Copy the file content to the zip
		}
		return nil
	})
}

// given a version string, return a number
// example: 1.2.3 -> 10203
// how it works:
// 1. split the version string by "."
// 2. reverse the array
// 3. for each element, convert it to an integer
// 4. multiply the integer by 100^index
// 5. add the integer to the total
// 6. return the total
func FormatVersionStr(v string) int64 {
	vs := strings.Split(v, ".")
	if len(vs) <= 0 {
		log.Panic("Version str error")
	}
	var vNum int64
	ReverseArr(vs)
	for index, v := range vs {
		num, err := strconv.ParseInt(v, 10, 64)
		if err != nil {
			log.Panic(err.Error())
		}
		for i := 0; i < index; i++ {
			num = num * 100
		}
		vNum += num
	}
	return vNum
}

func ReverseArr(s interface{}) {
	sort.SliceStable(s, func(i, j int) bool {
		return true
	})
}

// how to use:
// md5, err := FileMD5("path/to/file")
//
//	if err != nil {
//		log.Panic(err)
//	}
//
// fmt.Println(md5)
func FileMD5(filePath string) (string, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return "", err
	}
	hash := md5.New()
	_, _ = io.Copy(hash, file)
	return hex.EncodeToString(hash.Sum(nil)), nil
}

// given a string, return the md5 hash
// example: "hello" -> "5d41402abc4b2a76b9719d911017c592"
func MD5(str string) string {
	data := []byte(str)
	has := md5.Sum(data)
	md5str1 := fmt.Sprintf("%x", has)
	return md5str1
}

// check if a pth exists
func PathExists(path string) (bool, error) {
	_, err := os.Stat(path)
	if err == nil {
		return true, nil
	}
	if os.IsNotExist(err) {
		return false, nil
	}
	return false, err
}

func GetBaseBucketUrl(env string) string {
	if env == "production" {
		return PROD_BASE_BUCKET_URL
	}
	return DEV_BASE_BUCKET_URL
}
