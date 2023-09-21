package discord

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
)

const (
	discordApi = "https://discord.com/api/v9"

	fakeUserAgent = "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/116.0.0.0 Safari/537.36"
)

func (s *Session) discordRequest(method, endpoint string, body *bytes.Buffer, headers map[string]string, result any) error {
	if body == nil {
		body = bytes.NewBuffer([]byte{})
	}

	request, err := http.NewRequest(
		method,
		discordApi+endpoint,
		body,
	)
	if err != nil {
		return err
	}
	defer request.Body.Close()

	request.Header.Set("Content-type", "application/json")
	request.Header.Set("Accept", "application/json")
	request.Header.Set("User-Agent", fakeUserAgent)

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
	return s.discordRequest(method, endpoint, body, nil, &result)
}

func (s *Session) AuthenticatedRequest(method, endpoint string, body *bytes.Buffer, result any) error {
	if s.authorizationToken == "" {
		return errors.New("unable to make an authenticated request because there is no authorizationToken set")
	}

	headers := map[string]string{
		"Authorization": s.authorizationToken,
	}

	return s.discordRequest(method, endpoint, body, headers, &result)
}
