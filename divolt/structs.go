package divolt

import (
	"time"

	"github.com/gorilla/websocket"
)

type Session struct {
	sessionToken string
	userId       string
	socket       *Socket
}

type Socket struct {
	open            bool
	connection      *websocket.Conn
	messageChannels []*chan SocketResponse
}

type Message struct {
	Id      string `json:"_id"`
	Nonce   string `json:"nonce"`
	Channel string `json:"channel"`
	Author  string `json:"author"`
	Webhook struct {
		Name   string `json:"name"`
		Avatar string `json:"avatar"`
	} `json:"webhook"`
	Content string `json:"content"`
	System  struct {
		Type    string `json:"type"`
		Content string `json:"content"`
	} `json:"system"`
	Attachments []struct {
		Id       string `json:"_id"`
		Tag      string `json:"tag"`
		Filename string `json:"filename"`
		Metadata struct {
			Type string `json:"type"`
		} `json:"metadata"`
		ContentType string `json:"content_type"`
		Size        int    `json:"size"`
		Deleted     bool   `json:"deleted"`
		Reported    bool   `json:"reported"`
		MessageId   string `json:"message_id"`
		UserId      string `json:"user_id"`
		ServerId    string `json:"server_id"`
		ObjectId    string `json:"object_id"`
	} `json:"attachments"`
	Edited time.Time `json:"edited"`
	Embeds []struct {
		Type        string `json:"type"`
		URL         string `json:"url"`
		OriginalUrl string `json:"original_url"`
		Special     struct {
			Type string `json:"type"`
		} `json:"special"`
		Title       string `json:"title"`
		Description string `json:"description"`
		Image       struct {
			Url    string `json:"url"`
			Width  int    `json:"width"`
			Height int    `json:"height"`
			Size   string `json:"size"`
		} `json:"image"`
		Video struct {
			Url    string `json:"url"`
			Width  int    `json:"width"`
			Height int    `json:"height"`
		} `json:"video"`
		SiteName string `json:"site_name"`
		IconUrl  string `json:"icon_url"`
		Colour   string `json:"colour"`
	} `json:"embeds"`
	Mentions     []string            `json:"mentions"`
	Replies      []string            `json:"replies"`
	Reactions    map[string][]string `json:"reactions"`
	Interactions struct {
		Reactions         []string `json:"reactions"`
		RestrictReactions bool     `json:"restrict_reactions"`
	} `json:"interactions"`
	Masquerade struct {
		Name   string `json:"name"`
		Avatar string `json:"avatar"`
		Colour string `json:"colour"`
	} `json:"masquerade"`
}

type User struct {
	Id            string `json:"_id"`
	Username      string `json:"username"`
	Discriminator string `json:"discriminator"`
	DisplayName   string `json:"display_name"`
	Avatar        struct {
		Id       string `json:"_id"`
		Tag      string `json:"tag"`
		Filename string `json:"filename"`
		Metadata struct {
			Type string `json:"type"`
		} `json:"metadata"`
		ContentType string `json:"content_type"`
		Size        int    `json:"size"`
		Deleted     bool   `json:"deleted"`
		Reported    bool   `json:"reported"`
		MessageId   string `json:"message_id"`
		UserId      string `json:"user_id"`
		ServerId    string `json:"server_id"`
		ObjectId    string `json:"object_id"`
	} `json:"avatar"`
	Relations []struct {
		Id     string `json:"_id"`
		Status string `json:"status"`
	} `json:"relations"`
	Badges int `json:"badges"`
	Status struct {
		Text     string `json:"text"`
		Presence string `json:"presence"`
	} `json:"status"`
	Profile struct {
		Content    string `json:"content"`
		Background struct {
			Id       string `json:"_id"`
			Tag      string `json:"tag"`
			Filename string `json:"filename"`
			Metadata struct {
				Type string `json:"type"`
			} `json:"metadata"`
			ContentType string `json:"content_type"`
			Size        int    `json:"size"`
			Deleted     bool   `json:"deleted"`
			Reported    bool   `json:"reported"`
			MessageId   string `json:"message_id"`
			UserId      string `json:"user_id"`
			ServerId    string `json:"server_id"`
			ObjectId    string `json:"object_id"`
		} `json:"background"`
	} `json:"profile"`
	Flags      int  `json:"flags"`
	Privileged bool `json:"privileged"`
	Bot        struct {
		Owner string `json:"owner"`
	} `json:"bot"`
	Relationship string `json:"relationship"`
	Online       bool   `json:"online"`
}

type Login struct {
	Result string `json:"result"`
	Id     string `json:"_id"`
	UserId string `json:"user_id"`
	Token  string `json:"token"`
	Name   string `json:"name"`
}

type SocketResponse struct {
	Type string `json:"type"`
	Data []byte `json:"data"`
}

type SocketMessage struct {
	Type string `json:"type"`
	Message
}

type SocketMessageUpdate struct {
	Type    string  `json:"type"`
	Id      string  `json:"id"`
	Channel string  `json:"channel"`
	Data    Message `json:"data"`
}
