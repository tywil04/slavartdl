package helpers

import (
	"archive/zip"
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/schollz/progressbar/v3"
)

func JsonApiRequest(method, url string, responseWriter any, data, headers map[string]string) error {
	rawData, err := json.Marshal(data)
	if err != nil {
		return err
	}
	dataBuffer := bytes.NewBuffer(rawData)

	request, err := http.NewRequest(method, url, dataBuffer)
	if err != nil {
		return err
	}

	for key, value := range headers {
		request.Header.Set(key, value)
	}

	client := http.Client{}

	response, err := client.Do(request)
	if err != nil {
		return err
	}

	// if the response is not a successful http status
	if !(response.StatusCode >= 200 && response.StatusCode <= 299) {
		return fmt.Errorf("unsuccessful request got http status code: %d", response.StatusCode)
	}

	if err := json.NewDecoder(response.Body).Decode(&responseWriter); err != nil {
		return err
	}

	return nil
}

func DownloadFile(url string, outputFilePath string) error {
	response, err := http.Get(url)
	if err != nil {
		return err
	}
	defer response.Body.Close()

	file, err := os.Create(outputFilePath)
	if err != nil {
		return err
	}
	defer file.Close()

	fmt.Println()
	bar := progressbar.DefaultBytes(response.ContentLength)
	if _, err = io.Copy(io.MultiWriter(file, bar), response.Body); err != nil {
		return err
	}

	return nil
}

func unzipFile(file *zip.File, outputFolderPath string, ignoreSubdirectories, ignoreCover bool) error {
	filePath := filepath.Join(outputFolderPath, file.Name)
	if !strings.HasPrefix(filePath, filepath.Clean(outputFolderPath)+string(os.PathSeparator)) {
		return errors.New("invalid file path")
	}

	fileNameOnly := filepath.Base(filePath)

	if !ignoreSubdirectories {
		if file.FileInfo().IsDir() {
			if err := os.MkdirAll(filePath, os.ModePerm); err != nil {
				return err
			}

			return nil
		}
	} else {
		if file.FileInfo().IsDir() {
			return nil
		}

		filePath = filepath.Join(outputFolderPath, fileNameOnly)
	}

	if !(fileNameOnly == "cover.jpg" && ignoreCover) {
		destinationFile, err := os.Create(filePath)
		if err != nil {
			return err
		}
		defer destinationFile.Close()

		zippedFile, err := file.Open()
		if err != nil {
			return err
		}
		defer zippedFile.Close()

		fmt.Println("\n" + fileNameOnly)
		bar := progressbar.DefaultBytes(file.FileInfo().Size())
		if _, err := io.Copy(io.MultiWriter(destinationFile, bar), zippedFile); err != nil {
			return err
		}
	}

	return nil
}

func Unzip(inputFilePath, outputFolderPath string, ignoreSubdirectories, ignoreCover bool) error {
	reader, err := zip.OpenReader(inputFilePath)
	if err != nil {
		return err
	}
	defer reader.Close()

	for _, file := range reader.File {
		err := unzipFile(file, outputFolderPath, ignoreSubdirectories, ignoreCover)
		if err != nil {
			return err
		}
	}

	return nil
}

func GetZipName(inputFilePath string) (string, error) {
	zip, err := zip.OpenReader(inputFilePath)
	if err != nil {
		return "", err
	}
	defer zip.Close()
	return zip.File[0].Name, nil
}

func CopyFile(sourcePath, destinationPath string) error {
	source, err := os.Open(sourcePath)
	if err != nil {
		return err
	}
	defer source.Close()

	destination, err := os.Create(destinationPath)
	if err != nil {
		return err
	}
	defer destination.Close()

	if _, err := io.Copy(destination, source); err != nil {
		return err
	}

	return nil
}
