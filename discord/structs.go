package discord

import (
	"context"
	"time"

	"github.com/gorilla/websocket"
)

type Session struct {
	authorizationToken string
	userId             string
	socket             *Socket
}

type Socket struct {
	open               bool
	connection         *websocket.Conn
	messageChannels    []*chan SocketResponse
	heartbeatCtx       context.Context
	heartbeatCancel    context.CancelFunc
	lastSequenceNumber int
}

type Message struct {
	Reactions []struct {
		Count        int `json:"count"`
		CountDetails struct {
			Burst  int `json:"burst"`
			Normal int `json:"normal"`
		} `json:"count_details"`
		Me    bool `json:"me"`
		Emoji struct {
			Id            string   `json:"id"`
			Name          string   `json:"name"`
			Roles         []string `json:"roles"`
			User          User     `json:"user"`
			RequireColons bool     `json:"require_colons"`
			Managed       bool     `json:"managed"`
			Animated      bool     `json:"animated"`
			Available     bool     `json:"available"`
		} `json:"emoji"`
		BurstColors []string `json:"burst_colors"`
	} `json:"reactions"`
	Attachments []struct {
		Id              string  `json:"id"`
		Filename        string  `json:"filename"`
		Description     string  `json:"description"`
		ContentType     string  `json:"content_type"`
		Size            int     `json:"size"`
		Url             string  `json:"url"`
		ProxyUrl        string  `json:"proxy_url"`
		Height          int     `json:"height"`
		Width           int     `json:"width"`
		Emphemeral      bool    `json:"ephemeral"`
		DurationSeconds float64 `json:"duration_secs"`
		Waveform        string  `json:"waveform"`
		Flags           int     `json:"flags"`
	} `json:"attachments"`
	Tts    bool `json:"tts"`
	Embeds []struct {
		Title       string `json:"title"`
		Type        string `json:"type"`
		Description string `json:"description"`
		Url         string `json:"url"`
		Timestamp   string `json:"timestamp"`
		Color       int    `json:"color"`
		Footer      struct {
			Text         string  `json:"text"`
			IconUrl      *string `json:"icon_url"`
			ProxyIconUrl *string `json:"proxy_icon_url"`
		} `json:"footer"`
		Image struct {
			Url      string `json:"url"`
			ProxyUrl string `json:"proxy_url"`
			Height   int    `json:"height"`
			Width    int    `json:"width"`
		} `json:"image"`
		Thumbnail struct {
			Url      string `json:"url"`
			ProxyUrl string `json:"proxy_url"`
			Height   int    `json:"height"`
			Width    int    `json:"width"`
		} `json:"thumbnail"`
		Video struct {
			Url      string `json:"url"`
			ProxyUrl string `json:"proxy_url"`
			Height   int    `json:"height"`
			Width    int    `json:"width"`
		} `json:"video"`
		Provider struct {
			Name string `json:"name"`
			Url  string `json:"url"`
		} `json:"provider"`
		Author struct {
			Name         string `json:"name"`
			Url          string `json:"url"`
			IconUrl      string `json:"icon_url"`
			ProxyIconUrl string `json:"proxy_icon_url"`
		} `json:"author"`
		Fields []struct {
			Name   string `json:"name"`
			Value  string `json:"value"`
			Inline bool   `json:"inline"`
		} `json:"fields"`
	} `json:"embeds"`
	Timestamp       string   `json:"timestamp"`
	MentionEveryone bool     `json:"mention_everyone"`
	Id              string   `json:"id"`
	Pinned          bool     `json:"pinned"`
	EditedTimestamp string   `json:"edited_timestamp"`
	Author          User     `json:"author"`
	MentionRoles    []string `json:"mention_roles"`
	Content         string   `json:"content"`
	ChannelId       string   `json:"channel_id"`
	Mentions        []User   `json:"mentions"`
	Type            int      `json:"type"`
}

type GuildMember struct {
	User                       User      `json:"user,omitempty"`
	Nick                       string    `json:"nick"`
	Avatar                     string    `json:"avatar"`
	Roles                      []string  `json:"roles"`
	JoinedAt                   time.Time `json:"joined_at"`
	PremiumSince               time.Time `json:"premium_since"`
	Deaf                       bool      `json:"deaf"`
	Mute                       bool      `json:"mute"`
	Flags                      int       `json:"flags"`
	Pending                    bool      `json:"pending"`
	Permissions                string    `json:"permissions"`
	CommunicationDisabledUntil time.Time `json:"communication_disabled_until"`
}

type Channel struct {
	Id                            string `json:"id"`
	Type                          int    `json:"type"`
	GuildId                       string `json:"guild_id"`
	Position                      int    `json:"position"`
	Name                          string `json:"name"`
	Topic                         string `json:"topic"`
	Nsfw                          bool   `json:"nsfw"`
	LastMessageId                 string `json:"last_message_id"`
	Bitrate                       int    `json:"bitrate"`
	UserLimit                     int    `json:"user_limit"`
	RateLimitPerUser              int    `json:"rate_limit_per_user"`
	Recipients                    []User `json:"recipients"`
	Icon                          string `json:"icon"`
	OwnerId                       string `json:"owner_id"`
	ApplicationId                 string `json:"application_id"`
	Managed                       bool   `json:"managed"`
	ParentId                      string `json:"parent_id"`
	LastPinTimestamp              string `json:"last_pin_timestamp"`
	RtcRegion                     string `json:"rtc_region"`
	VideoQualityMode              int    `json:"video_quality_mode,"`
	MessageCount                  int    `json:"message_count"`
	MemberCount                   int    `json:"member_count"`
	DefaultAutoArchiveDuration    int    `json:"default_auto_archive_duration"`
	Permissions                   string `json:"permissions"`
	Flags                         int    `json:"flags"`
	TotalMessageSent              int    `json:"total_message_sent"`
	DefaultThreadRateLimitPerUser int    `json:"default_thread_rate_limit_per_user"`
	DefaultSortOrder              int    `json:"default_sort_order"`
	DefaultForumLayout            int    `json:"default_forum_layout"`
}

type User struct {
	Id            string `json:"id"`
	Username      string `json:"username"`
	Discriminator string `json:"discriminator"`
	Avatar        string `json:"avatar"`
	Verified      bool   `json:"verified"`
	Email         string `json:"email"`
	Flags         int    `json:"flags"`
	Banner        string `json:"banner"`
	AccentColor   int    `json:"accent_color"`
	PremiumType   int    `json:"premium_type"`
	PublicFlags   int    `json:"public_flags"`
}

type SocketResponse struct {
	Op int    `json:"op"`
	D  any    `json:"d"`
	S  int    `json:"s"`
	T  string `json:"t"`
}

type SocketMessage SocketResponse

type SocketMessageCreate struct {
	Message
	GuildId           string  `json:"guild_id"`
	Member            any     `json:"member"`
	ReferencedMessage Message `json:"referenced_message"`
	Mentions          []struct {
		User
		Member GuildMember
	} `json:"mentions"`
}

type SocketMessageUpdate SocketMessageCreate

type Login struct {
	UserId       string `json:"user_id"`
	Token        string `json:"token"`
	UserSettings struct {
		Locale string `json:"locale"`
		Theme  string `json:"theme"`
	} `json:"user_settings"`
}
