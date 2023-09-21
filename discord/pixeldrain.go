package discord

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"regexp"
	"strings"
	"time"

	"github.com/tywil04/slavartdl/privatebin"
)

const (
	pixeldrainInviteId = "WsJu2zh8"

	pixeldrainRequestChannel = "1151600423220822147"
	pixeldrainBotId          = "1152954682877149234"
)

var (
	PixeldrainAllowedHosts = []string{
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

	pixeldrainRegex = regexp.MustCompile(`(?ms).*https:\/\/pixeldrain\.com\/u\/(.{8}).*`)

	//slavartRequestedUrlAndDownloadLinkRegex = regexp.MustCompile(`(?m)Your requested link\, (.*)\, is now available for download:\n \*\*Download Link\*\*\n (.*)`)
)

func (s *Session) PixeldrainTryInviteUser() error {
	s.AuthenticatedRequest(
		http.MethodPost,
		"/invites/"+pixeldrainInviteId,
		nil,
		nil,
	)

	return nil
}

func (s *Session) PixeldrainSendDownloadCommand(url string, quality int) (*Message, error) {
	var command string
	if quality == -1 {
		command = fmt.Sprintf(`$dl %s`, url)
	} else {
		command = fmt.Sprintf(`$dl %s %d`, url, quality)
	}

	message, err := s.SendMessage(pixeldrainRequestChannel, command)
	if err != nil {
		return nil, err
	}

	return message, nil
}

func (s *Session) PixeldrainGetUploadUrl(downloadMessageId, requestUrl string, timeout time.Duration) (string, error) {
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
			response := <-messageChannel

			if response.Op == 0 {
				switch response.T {
				case "MESSAGE_CREATE":
					rawData, _ := json.Marshal(response.D)
					message := SocketMessageCreate{}
					json.Unmarshal(rawData, &message)

					if message.Author.Id == pixeldrainBotId {
						// dm message has no contents, only embeds
						if message.Content == "" {
							url := strings.Trim(strings.TrimSpace(message.Embeds[0].Fields[0].Value), "\n")
							url = strings.Split(strings.TrimPrefix(url, "https://links.gamesdrive.net/#/link/"), ".")[0]

							rawUrl, err := base64.RawURLEncoding.DecodeString(url)
							if err != nil {
								return "", err
							}
							url = string(rawUrl)

							if url != "" {
								url = strings.TrimPrefix(url, "https://paste.gamesdrive.net/?")
								split := strings.Split(url, "#")

								plaintext, err := privatebin.GetPaste(
									"https://paste.gamesdrive.net",
									split[0],
									split[1],
								)
								if err != nil {
									return "", err
								}

								pixeldrainDownloadId := pixeldrainRegex.FindAllStringSubmatch(plaintext, -1)[0][1]
								downloadUrl := fmt.Sprintf("https://pixeldrain.com/api/file/%s?download", pixeldrainDownloadId)

								return downloadUrl, nil
							}
						}
					}

					// message is a reply
					if message.Author.Id == pixeldrainBotId && message.Type == 19 {
						var mentionsUser bool
						for _, mention := range message.Mentions {
							if mention.Id == s.userId {
								mentionsUser = true
								break
							}
						}

						repliedToDownloadMessage := message.ReferencedMessage.Author.Id == s.userId
						containsError := strings.Contains(strings.ToLower(message.Content), "error")

						if mentionsUser && message.ChannelId == pixeldrainRequestChannel && message.Author.Id == pixeldrainBotId && repliedToDownloadMessage && containsError {
							return "", errors.New("there was an error processing your request. maybe you are trying to use a disabled service?")
						}
					}
				case "MESSAGE_UPDATE":
					rawData, _ := json.Marshal(response.D)
					message := SocketMessageUpdate{}
					json.Unmarshal(rawData, &message)

					if message.Id == downloadMessageId {
						if strings.Contains(strings.ToLower(message.Content), "error") {
							return "", errors.New("there was an error with your download request, therefore no upload message will be sent")
						}
					}
				}
			}
		}
	}
}
