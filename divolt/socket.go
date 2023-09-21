package divolt

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"github.com/gorilla/websocket"
)

const (
	divoltWebsocket = "wss://ws.divolt.xyz"
)

func (s *Socket) StartListening() {
	for {
		_, message, _ := s.connection.ReadMessage()
		response := struct {
			Type string `json:"type"`
		}{}
		json.Unmarshal(message, &response)

		for _, channel := range s.messageChannels {
			*channel <- SocketResponse{
				Type: response.Type,
				Data: message,
			}
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
	return s.connection.Close()
}

func (s *Session) OpenAuthenticatedSocket() (*Socket, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Minute)
	defer cancel()

	dialer := websocket.Dialer{}

	connection, _, err := dialer.Dial(divoltWebsocket, nil)
	if err != nil {
		return nil, err
	}

	payload := fmt.Sprintf(`{"type":"Authenticate", "token":"%s"}`, s.sessionToken)
	if err := connection.WriteMessage(websocket.TextMessage, []byte(payload)); err != nil {
		return nil, err
	}

	for {
		select {
		case <-ctx.Done():
			return nil, errors.New("timed out before authenticated")
		default:
			message := map[string]string{}
			err := connection.ReadJSON(&message)
			if err != nil {
				return nil, err
			}

			if message["type"] == "Authenticated" {
				socket := &Socket{
					connection: connection,
					open:       true,
				}
				go socket.StartListening()
				return socket, nil
			}
		}
	}
}
