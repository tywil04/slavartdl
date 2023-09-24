package discord

import (
	"bytes"
	"errors"
	"fmt"
	"net/http"
)

func (s *Session) AuthenticateWithCredentials(email, password string) error {
	if s.authorizationToken != "" {
		return errors.New("you have already authenticated. if you want to change the authorizationToken, please logout then login again")
	}

	loginResponse := Login{}

	payload := fmt.Sprintf(`{"login":"%s", "password":"%s", "undelete": false, "login_source": null, "gift_code_sku_id": null }`, email, password)

	err := s.UnauthenticatedRequest(
		http.MethodPost,
		"/auth/login",
		bytes.NewBufferString(payload),
		&loginResponse,
	)
	if err != nil {
		return err
	}

	s.authorizationToken = loginResponse.Token
	s.userId = loginResponse.UserId

	s.socket, err = s.OpenAuthenticatedSocket()
	if err != nil {
		return err
	}

	return nil
}

func (s *Session) AuthenticateWithAuthorizationToken(authorizationToken string) error {
	if s.authorizationToken != "" {
		return errors.New("you have already authenticated. if you want to change the sessionToken, please logout then login again")
	}

	if authorizationToken == "" {
		return errors.New("the sessionToken provided is empty")
	}

	s.authorizationToken = authorizationToken

	user, err := s.GetAuthenticatedUserInfo()
	if err != nil {
		return err
	}

	s.userId = user.Id

	s.socket, err = s.OpenAuthenticatedSocket()
	if err != nil {
		return err
	}

	return nil
}

func (s *Session) Logout() error {
	if s.authorizationToken == "" {
		return errors.New("you are not authenticated, so you cannot logout")
	}

	s.authorizationToken = ""
	s.userId = ""
	go s.socket.Close()

	return nil
}

func (s *Session) GetAuthenticatedUserInfo() (*User, error) {
	userResponse := &User{}

	err := s.AuthenticatedRequest(
		http.MethodGet,
		"/users/@me",
		nil,
		&userResponse,
	)
	if err != nil {
		return nil, err
	}

	return userResponse, nil
}

func (s *Session) SendMessage(channelId, message string) (*Message, error) {
	payload := fmt.Sprintf(`{"content": "%s"}`, message)
	response := &Message{}

	err := s.AuthenticatedRequest(
		http.MethodPost,
		"/channels/"+channelId+"/messages",
		bytes.NewBufferString(payload),
		&response,
	)
	if err != nil {
		return nil, err
	}

	return response, nil
}
