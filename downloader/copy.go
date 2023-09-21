package downloader

import (
	"bytes"
	"io"
	"os"
)

func CopyFile(source *bytes.Reader, destinationPath string) error {
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
