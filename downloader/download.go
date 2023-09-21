package downloader

import (
	"bytes"
	"io"
	"net/http"
)

func DownloadFile(url string) (*bytes.Reader, int64, error) {
	response, err := http.Get(url)
	if err != nil {
		return nil, 0, err
	}
	defer response.Body.Close()

	buffer := bytes.NewBuffer([]byte{})
	bytesWritten, err := io.Copy(buffer, response.Body)
	if err != nil {
		return nil, 0, err
	}

	reader := bytes.NewReader(buffer.Bytes())

	return reader, bytesWritten, nil
}
