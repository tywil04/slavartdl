package downloader

import (
	"archive/zip"
	"bytes"
	"errors"
	"io"
	"os"
	"path/filepath"
	"strings"
)

type UnzipOptions struct {
	IgnoreSubdirs bool
	IgnoreCover   bool
}

func unzip(file *zip.File, outputFolderPath string, ignoreSubdirs, ignoreCover bool) error {
	filePath := filepath.Join(outputFolderPath, file.Name)
	if !strings.HasPrefix(filePath, filepath.Clean(outputFolderPath)+string(os.PathSeparator)) {
		return errors.New("invalid file path")
	}

	// protect against zip slip
	if strings.Contains(filePath, "..") {
		return errors.New("invalid file path")
	}

	fileNameOnly := filepath.Base(filePath)

	if !ignoreSubdirs {
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

		if _, err := io.Copy(destinationFile, zippedFile); err != nil {
			return err
		}
	}

	return nil
}

func Unzip(inputFile *bytes.Reader, size int64, outputPath string, ignoreSubdirs, ignoreCover bool) error {
	reader, err := zip.NewReader(inputFile, size)
	if err != nil {
		return err
	}

	for _, file := range reader.File {
		err := unzip(file, outputPath, ignoreSubdirs, ignoreCover)
		if err != nil {
			return err
		}
	}

	return nil
}
