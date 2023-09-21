package divolt

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"regexp"
	"strings"
	"time"
)

const (
	slavartRequestChannel = "01G9AZ9AMWDV227YA7FQ5RV8WB"
	slavartUploadChannel  = "01G9AZ9Q2R5VEGVPQ4H99C01YP"
	slavartBotId          = "01G9824MQPGD7GVYR0F6A6GJ2Q"

	SlavartBotStatusOnline  = "online"
	SlavartBotStatusOffline = "offline"
)

var (
	SlavartAllowedHosts = []string{
		"tidal.com",
		"www.qobuz.com",
		"play.qobuz.com",
		"open.qobuz.com",
		"soundcloud.com",
		"www.deezer.com",
		"open.spotify.com",
		"music.youtube.com",
		"www.jiosaavn.com",
	}

	slavartRequestedUrlAndDownloadLinkRegex = regexp.MustCompile(`(?m)Your requested link\, (.*)\, is now available for download:\n \*\*Download Link\*\*\n (.*)`)
)

func (s *Session) SlavartTryInviteUser() error {
	request, err := http.NewRequest(http.MethodGet, "https://slavart.divolt.xyz", nil)
	if err != nil {
		return err
	}

	client := http.Client{}
	response, err := client.Do(request)
	if err != nil {
		return err
	}

	inviteId := response.Request.URL.String()

	s.AuthenticatedRequest(
		http.MethodPost,
		"/invites/"+inviteId,
		nil,
		nil,
	)

	return nil
}

func (s *Session) SlavartGetBotStatus() (string, error) {
	response := User{}

	err := s.AuthenticatedRequest(
		http.MethodGet,
		"/users/"+slavartBotId,
		nil,
		&response,
	)
	if err != nil {
		return "", err
	}

	if response.Online {
		return SlavartBotStatusOnline, nil
	} else {
		return SlavartBotStatusOffline, nil
	}
}

func (s *Session) SlavartSendDownloadCommand(url string, quality int) (*Message, error) {
	var command string
	if quality == -1 {
		command = fmt.Sprintf(`!dl %s`, url)
	} else {
		command = fmt.Sprintf(`!dl %s %d`, url, quality)
	}

	message, err := s.SendMessage(slavartRequestChannel, command)
	if err != nil {
		return nil, err
	}

	return message, nil
}

func (s *Session) SlavartGetUploadUrl(downloadMessageId, requestUrl string, timeout time.Duration) (string, error) {
	if !s.socket.open {
		return "", errors.New("socket is not open")
	}

	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	messageChannel := make(chan SocketResponse)
	s.socket.RegisterMessageChannel(&messageChannel)
	defer s.socket.DeregisterMessageChannel(&messageChannel)

	for {
		select {
		case <-ctx.Done():
			return "", errors.New("timed out before message could be found")
		default:
			rawResponse := <-messageChannel

			switch rawResponse.Type {
			case "Message":
				message := SocketMessage{}
				json.Unmarshal(rawResponse.Data, &message)

				// check to see if message is an upload message
				var mentionsUser bool
				for _, mention := range message.Mentions {
					if mention == s.userId {
						mentionsUser = true
						break
					}
				}

				if message.Channel == slavartUploadChannel && message.Author == slavartBotId && mentionsUser {
					matches := slavartRequestedUrlAndDownloadLinkRegex.FindAllStringSubmatch(message.Embeds[0].Description, -1)[0]

					if matches[1] == requestUrl {
						return matches[2], nil
					}
				}

				// check to see if message is an error replying to the download message
				var repliedToDownloadMessage bool
				for _, reply := range message.Replies {
					if reply == downloadMessageId {
						repliedToDownloadMessage = true
						break
					}
				}

				containsError := strings.Contains(strings.ToLower(message.Content), "error")

				if message.Channel == slavartRequestChannel && message.Author == slavartBotId && repliedToDownloadMessage && containsError {
					return "", errors.New("there was an error processing your request. maybe you are trying to use a disabled service?")
				}
			case "MessageUpdate":
				message := SocketMessageUpdate{}
				json.Unmarshal(rawResponse.Data, &message)

				if message.Id == downloadMessageId {
					if strings.Contains(strings.ToLower(message.Data.Content), "error") {
						return "", errors.New("there was an error with your download request, therefore no upload message will be sent")
					}
				}
			}
		}
	}
}
