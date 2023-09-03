package common

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
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
		var body = bytes.NewBuffer([]byte{})
		io.Copy(body, response.Body)

		return fmt.Errorf("unsuccessful request got http status code: %d, with a body of: %s", response.StatusCode, body.String())
	}

	if err := json.NewDecoder(response.Body).Decode(&responseWriter); err != nil {
		return err
	}

	return nil
}
