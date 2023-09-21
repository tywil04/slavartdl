package divolt

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
)

const (
	divoltApi = "https://api.divolt.xyz"
)

func (s *Session) divoltRequest(method, endpoint string, body *bytes.Buffer, headers map[string]string, result any) error {
	if body == nil {
		body = bytes.NewBuffer([]byte{})
	}

	request, err := http.NewRequest(
		method,
		divoltApi+endpoint,
		body,
	)
	if err != nil {
		return err
	}
	defer request.Body.Close()

	request.Header.Set("Content-type", "application/json")
	request.Header.Set("Accept", "application/json")
	request.Header.Set("User-Agent", "Slavartdl (contact at github.26ac3d4d@alias.tylerw.co.uk)")

	for key, value := range headers {
		request.Header.Set(key, value)
	}

	client := http.Client{}
	response, err := client.Do(request)
	if err != nil {
		return err
	}
	defer response.Body.Close()

	if response.StatusCode < 200 || response.StatusCode > 300 {
		return errors.New("request returned unsuccessful http status code")
	}

	if err := json.NewDecoder(response.Body).Decode(&result); err != nil {
		return err
	}

	return nil
}

func (s *Session) UnauthenticatedRequest(method, endpoint string, body *bytes.Buffer, result any) error {
	return s.divoltRequest(method, endpoint, body, nil, &result)
}

func (s *Session) AuthenticatedRequest(method, endpoint string, body *bytes.Buffer, result any) error {
	if s.sessionToken == "" {
		return errors.New("unable to make an authenticated request because there is no sessionToken set")
	}

	headers := map[string]string{
		"X-Session-Token": s.sessionToken,
	}

	return s.divoltRequest(method, endpoint, body, headers, &result)
}
