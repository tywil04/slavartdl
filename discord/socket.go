package discord

import (
	"context"
	"encoding/json"
	"errors"
	"time"

	"github.com/gorilla/websocket"
)

const (
	discordWebsocket = "wss://gateway.discord.gg/?v=10&encoding=json"
)

func (s *Socket) StartListening() {
	for {
		message := SocketResponse{}
		s.connection.ReadJSON(&message)

		for _, channel := range s.messageChannels {
			*channel <- message
		}
	}
}

func (s *Socket) RegisterMessageChannel(channel *chan SocketResponse) {
	s.messageChannels = append(s.messageChannels, channel)
}

func (s *Socket) DeregisterMessageChannel(channel *chan SocketResponse) {
	var indexOf = -1
	for index, messageChannel := range s.messageChannels {
		if messageChannel == channel {
			indexOf = index
			break
		}
	}

	if indexOf != -1 {
		s.messageChannels = append(s.messageChannels[:indexOf], s.messageChannels[indexOf+1:]...)
	}
}

func (s *Socket) Close() error {
	s.heartbeatCancel()
	return s.connection.Close()
}

func (s *Session) OpenAuthenticatedSocket() (*Socket, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Minute)
	defer cancel()

	dialer := websocket.Dialer{}

	connection, _, err := dialer.Dial(discordWebsocket, nil)
	if err != nil {
		return nil, err
	}

	var socket *Socket

	for {
		select {
		case <-ctx.Done():
			return nil, errors.New("timed out before authenticated")
		default:
			message := SocketResponse{}
			err := connection.ReadJSON(&message)
			if err != nil {
				return nil, err
			}

			// opcode 10 is "hello"
			if message.Op == 10 && socket == nil {
				heartbeatCtx, heatbeatCancel := context.WithCancel(context.Background())
				go func() {
					helloData := struct {
						HeartbeatInterval int `json:"heartbeat_interval"`
					}{}
					raw, _ := json.Marshal(message.D)
					json.Unmarshal(raw, &helloData)

					interval := time.Millisecond * time.Duration(helloData.HeartbeatInterval)

					for {
						select {
						case <-heartbeatCtx.Done():
							return
						default:
							var rawHeartbeat []byte
							if socket != nil {
								heartbeat := SocketMessage{
									Op: 1,
									D:  socket.lastSequenceNumber,
								}
								rawHeartbeat, _ = json.Marshal(heartbeat)
							} else {
								heartbeat := SocketMessage{
									Op: 1,
									D:  nil,
								}
								rawHeartbeat, _ = json.Marshal(heartbeat)
							}

							connection.WriteMessage(websocket.TextMessage, rawHeartbeat)
							time.Sleep(interval)
						}
					}
				}()

				identify := map[string]any{
					"op": 2,
					"d": map[string]any{
						"token":   s.authorizationToken,
						"intents": 37376,
						"properties": map[string]string{
							"os":      "windows",
							"browser": "go",
							"device":  "go",
						},
					},
				}
				rawIdentify, _ := json.Marshal(identify)

				connection.WriteMessage(websocket.TextMessage, rawIdentify)

				socket = &Socket{
					connection:      connection,
					open:            false,
					heartbeatCtx:    heartbeatCtx,
					heartbeatCancel: heatbeatCancel,
				}
			}

			// opcode 0 is a dispatch event
			if message.Op == 0 && message.T == "READY" && socket != nil {
				socket.open = true
				go socket.StartListening()
				return socket, nil
			}
		}
	}
}
