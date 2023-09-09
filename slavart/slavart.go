package slavart

import (
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/tywil04/slavartdl/common"
	"github.com/tywil04/slavartdl/slavart/internal/helpers"
)

const (
	Api            = "https://api.divolt.xyz"
	InviteRedirect = "https://slavart.divolt.xyz"

	RequestChannel = "01G9AZ9AMWDV227YA7FQ5RV8WB"
	UploadChannel  = "01G9AZ9Q2R5VEGVPQ4H99C01YP"
	BotId          = "01G9824MQPGD7GVYR0F6A6GJ2Q"
)

var (
	DownloadLinkInUploadMessageRegex = regexp.MustCompile(`(?m)Your requested link\, (.*)\, is now available for download:\n \*\*Download Link\*\*\n (.*)`)
	DivoltInviteIdFromInviteUrlRegex = regexp.MustCompile(`https://divolt\.xyz.invite/(.*)`)
)

var AllowedHosts = []string{
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

// interact with the slavart divolt server (its a self-hosted instance of revolt) and use its api

// this only has what I need
type RevoltMessage struct {
	Mentions []string `json:"mentions"`
	Embeds   []struct {
		Description string `json:"description"`
	} `json:"embeds"`
	Content string   `json:"content"`
	Replies []string `json:"replies"`
}

// check if bot is down
func GetBotOnlineStatus(sessionToken string) (bool, error) {
	serverMemberResponse := struct {
		Online bool `json:"online"`
	}{}

	err := common.JsonApiRequest(
		http.MethodGet,
		Api+"/users/"+BotId,
		&serverMemberResponse,
		map[string]string{},
		map[string]string{
			"X-Session-Token": sessionToken,
		},
	)

	return serverMemberResponse.Online, err
}

// send message in request channel
func SendDownloadMessage(sessionToken, link string, quality int) (string, error) {
	downloadRequestResponse := struct {
		MessageId string `json:"_id"`
	}{}

	content := "!dl " + link
	if quality != -1 {
		content += " " + strconv.Itoa(quality)
	}

	err := common.JsonApiRequest(
		http.MethodPost,
		Api+"/channels/"+RequestChannel+"/messages",
		&downloadRequestResponse,
		map[string]string{
			"content": content,
		},
		map[string]string{
			"X-Session-Token": sessionToken,
		},
	)

	return downloadRequestResponse.MessageId, err
}

// get last 100 messages in upload channel
func GetUploadMessages(sessionToken string) ([]RevoltMessage, error) {
	downloadRequestFinishedTestResponse := []RevoltMessage{}

	err := common.JsonApiRequest(
		http.MethodGet,
		Api+"/channels/"+UploadChannel+"/messages",
		&downloadRequestFinishedTestResponse,
		map[string]string{
			"sort":          "Latest",
			"include_users": "false",
		},
		map[string]string{
			"X-Session-Token": sessionToken,
		},
	)

	return downloadRequestFinishedTestResponse, err
}

// get last 250 messages in request channel
func GetRequestMessages(sessionToken string) ([]RevoltMessage, error) {
	downloadRequestFinishedTestResponse := []RevoltMessage{}

	err := common.JsonApiRequest(
		http.MethodGet,
		Api+"/channels/"+RequestChannel+"/messages",
		&downloadRequestFinishedTestResponse,
		map[string]string{
			"sort":          "Latest",
			"include_users": "false",
			"limit":         "250",
		},
		map[string]string{
			"X-Session-Token": sessionToken,
		},
	)

	return downloadRequestFinishedTestResponse, err
}

func GetSessionTokenFromCredentials(email, password string) (string, error) {
	loginResponse := struct {
		Token string `json:"token"`
	}{}

	err := common.JsonApiRequest(
		http.MethodPost,
		Api+"/auth/session/login",
		&loginResponse,
		map[string]string{
			"email":         email,
			"password":      password,
			"friendly_name": "SlavartDL Tool (github.com/tywil04/slavartdl)",
		},
		nil,
	)

	return loginResponse.Token, err
}

func CheckForErrorMessageInRequestMessages(requestMessageId string, messages []RevoltMessage) (string, bool) {
	for _, revoltMessage := range messages {
		for _, reply := range revoltMessage.Replies {
			if reply == requestMessageId && strings.Contains(strings.ToLower(revoltMessage.Content), "error") {
				return revoltMessage.Content, true
			}
		}
	}

	return "", false
}

// see if we can find an upload message for the link we want (even if its for another user)
func SearchForDownloadLinkInUploadMessages(link string, messages []RevoltMessage) (string, bool) {
	for _, revoltMessage := range messages {
		matches := DownloadLinkInUploadMessageRegex.FindAllStringSubmatch(revoltMessage.Embeds[0].Description, -1)[0]

		if matches[1] == link {
			return matches[2], true
		}
	}

	return "", false
}

func GetDownloadLinkFromSlavart(sessionToken, link string, quality int, timeoutTime time.Time) (string, error) {
	messages, err := GetUploadMessages(sessionToken)
	if err != nil {
		return "", err
	}

	if downloadLink, ok := SearchForDownloadLinkInUploadMessages(link, messages); ok {
		return downloadLink, nil
	}

	online, err := GetBotOnlineStatus(sessionToken)
	if err != nil {
		return "", err
	}

	if !online {
		return "", errors.New("bot isn't online")
	}

	requestMessageId, err := SendDownloadMessage(sessionToken, link, quality)
	if err != nil {
		return "", err
	}

	time.Sleep(time.Second * 5) // give time for the message to send

	requestMessages, err := GetRequestMessages(sessionToken)
	if err != nil {
		return "", err
	}

	// error found
	if errMessage, errorFound := CheckForErrorMessageInRequestMessages(requestMessageId, requestMessages); errorFound {
		return "", errors.New(errMessage)
	}

	for {
		if timeoutTime.Before(time.Now()) {
			return "", errors.New("timed-out before download link could be found")
		}

		time.Sleep(time.Second * 5)

		messages, err := GetUploadMessages(sessionToken)
		if err != nil {
			return "", err
		}

		// check for later errors
		if errMessage, errorFound := CheckForErrorMessageInRequestMessages(requestMessageId, messages); errorFound {
			return "", errors.New(errMessage)
		}

		if downloadLink, ok := SearchForDownloadLinkInUploadMessages(link, messages); ok {
			return downloadLink, nil
		}
	}
}

// The whole download routine
func DownloadUrl(url, sessionToken, logLevel string, quality int, timeoutTime time.Time, cooldownDuration time.Duration, outputDir string, skipUnzip, ignoreCover, ignoreSubdirs bool) {
	if logLevel == "all" {
		fmt.Println("Getting download link...")
	}
	downloadLink, err := GetDownloadLinkFromSlavart(sessionToken, url, quality, timeoutTime)
	if err != nil {
		if logLevel == "all" || logLevel == "errors" {
			log.Fatal(err)
		}
	}

	if logLevel == "all" {
		fmt.Println("\nDownloading zip...")
	}
	// this will create a temp file in the default location
	tempFile, err := os.CreateTemp("", "slavartdl.*.zip")
	if err != nil {
		if logLevel == "all" || logLevel == "errors" {
			log.Fatal(err)
		}
	}
	defer os.Remove(tempFile.Name())

	tempFilePath := tempFile.Name()
	err = helpers.DownloadFile(downloadLink, tempFilePath, logLevel != "all")
	if err != nil {
		if logLevel == "all" || logLevel == "errors" {
			log.Fatal(err)
		}
	}

	if !skipUnzip {
		if logLevel == "all" {
			fmt.Println("\nUnzipping...")
		}
		if err := helpers.Unzip(tempFilePath, outputDir, ignoreSubdirs, ignoreCover, logLevel != "all"); err != nil {
			if logLevel == "all" || logLevel == "errors" {
				log.Fatal(err)
			}
		}
	} else {
		zipName, err := helpers.GetZipName(tempFilePath)
		if err != nil {
			if logLevel == "all" || logLevel == "errors" {
				log.Fatal(err)
			}
		}

		if logLevel == "all" {
			fmt.Println(filepath.Clean(zipName))
		}
		outputFileDir := outputDir + string(os.PathSeparator) + filepath.Clean(zipName) + ".zip"
		// temp file gets deleted later
		if err := helpers.CopyFile(tempFilePath, outputFileDir); err != nil {
			if logLevel == "all" || logLevel == "errors" {
				log.Fatal(err)
			}
		}

	}

	if logLevel == "all" {
		fmt.Println("\nDone!")
	}
}

func GetSlavartInviteId() (string, error) {
	redirectUrl, err := common.CaptureRedirectRequest("GET", InviteRedirect, nil)
	if err != nil {
		return "", err
	}

	matches := DivoltInviteIdFromInviteUrlRegex.FindAllStringSubmatch(redirectUrl, -1)
	return matches[0][1], nil
}
