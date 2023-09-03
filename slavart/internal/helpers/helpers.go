package helpers

import (
	"archive/zip"
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/schollz/progressbar/v3"
)

func DownloadFile(url string, outputFilePath string, disableLogs bool) error {
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

	if !disableLogs {
		fmt.Println()
		bar := progressbar.DefaultBytes(response.ContentLength)
		if _, err = io.Copy(io.MultiWriter(file, bar), response.Body); err != nil {
			return err
		}
	} else {
		if _, err = io.Copy(file, response.Body); err != nil {
			return err
		}
	}

	return nil
}

func unzipFile(file *zip.File, outputFolderPath string, ignoreSubdirectories, ignoreCover, disableLogs bool) error {
	filePath := filepath.Join(outputFolderPath, file.Name)
	if !strings.HasPrefix(filePath, filepath.Clean(outputFolderPath)+string(os.PathSeparator)) {
		return errors.New("invalid file path")
	}

	// protect against zip slip
	if strings.Contains(filePath, "..") {
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

		if !disableLogs {
			fmt.Println("\n" + fileNameOnly)
			bar := progressbar.DefaultBytes(file.FileInfo().Size())
			if _, err := io.Copy(io.MultiWriter(destinationFile, bar), zippedFile); err != nil {
				return err
			}
		} else {
			if _, err := io.Copy(destinationFile, zippedFile); err != nil {
				return err
			}
		}
	}

	return nil
}

func Unzip(inputFilePath, outputFolderPath string, ignoreSubdirectories, ignoreCover, disableLogs bool) error {
	reader, err := zip.OpenReader(inputFilePath)
	if err != nil {
		return err
	}
	defer reader.Close()

	for _, file := range reader.File {
		err := unzipFile(file, outputFolderPath, ignoreSubdirectories, ignoreCover, disableLogs)
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

	// protect against zip slip
	fileName := zip.File[0].Name
	if strings.Contains(fileName, "..") {
		return "", errors.New("invalid file path")
	}

	return fileName, nil
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
