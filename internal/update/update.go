package update

import (
	"bytes"
	"crypto"
	"encoding/hex"
	"fmt"
	"io"
	"net/http"
	"runtime"
	"strconv"
	"strings"

	"github.com/minio/selfupdate"

	"github.com/tywil04/slavartdl/internal/helpers"
)

var Version = "v1.1.5"

const assetNameTemplate = "slavartdl-%s-%s-%s.%s"
const signatureNameTemplate = "slavartdl-%s-%s-%s.%s.md5"

// accepts version like 'v0.0.0'.
// returns major, minor, patch ints
func parseVersionTag(version string) (int, int, int, error) {
	versionTrim := strings.TrimPrefix(Version, "v")
	majorMinorPatch := strings.Split(versionTrim, ".")

	if len(majorMinorPatch) != 3 {
		return 0, 0, 0, fmt.Errorf("unknown error when parsing version tag")
	}

	major, err := strconv.Atoi(majorMinorPatch[0])
	if err != nil {
		return 0, 0, 0, fmt.Errorf("failed to parse major part of version tag")
	}

	minor, err := strconv.Atoi(majorMinorPatch[1])
	if err != nil {
		return 0, 0, 0, fmt.Errorf("failed to parse minor part of version tag")
	}

	patch, err := strconv.Atoi(majorMinorPatch[2])
	if err != nil {
		return 0, 0, 0, fmt.Errorf("failed to parse patch part of version tag")
	}

	return major, minor, patch, nil
}

// checks if there is a new update avaiable, if there is it updates
func Update() (bool, error) {
	releasesResponse := struct {
		TagName string `json:"tag_name"`
		Assets  []struct {
			BrowserDownloadUrl string `json:"browser_download_url"`
			Name               string `json:"name"`
		} `json:"assets"`
	}{}

	err := helpers.JsonApiRequest(
		http.MethodGet,
		"https://api.github.com/repos/tywil04/slavartdl/releases/latest",
		&releasesResponse,
		nil,
		nil,
	)
	if err != nil {
		return false, err
	}

	apiMajor, apiMinor, apiPatch, err := parseVersionTag(releasesResponse.TagName)
	if err != nil {
		return false, err
	}

	pkgMajor, pkgMinor, pkgPatch, err := parseVersionTag(Version)
	if err != nil {
		return false, err
	}

	// check if current version is equal to or greater than fetched version
	if apiMajor <= pkgMajor && apiMinor <= pkgMinor && apiPatch <= pkgPatch {
		// no update, no error
		return false, nil
	}

	// update is available
	extension := "tar.gz"
	if runtime.GOOS == "windows" {
		extension = "zip"
	}

	assetName := fmt.Sprintf(assetNameTemplate, releasesResponse.TagName, runtime.GOOS, runtime.GOARCH, extension)
	signatureName := fmt.Sprintf(signatureNameTemplate, releasesResponse.TagName, runtime.GOOS, runtime.GOARCH, extension)

	var assetDownloadUrl string
	var signatureDownloadUrl string

	for _, asset := range releasesResponse.Assets {
		if assetDownloadUrl != "" && signatureDownloadUrl != "" {
			break // done
		}

		if asset.Name == assetName {
			assetDownloadUrl = asset.BrowserDownloadUrl
			continue
		}

		if asset.Name == signatureName {
			signatureDownloadUrl = asset.BrowserDownloadUrl
			continue
		}
	}

	if assetDownloadUrl == "" || signatureDownloadUrl == "" {
		return false, fmt.Errorf("failed to get asset and asset signature download url for your system, supported platforms are linux [arm64 amd64], darwin [arm64 amd64], windows [amd64]")
	}

	assetResponse, err := http.Get(assetDownloadUrl)
	if err != nil || (assetResponse.StatusCode < 200 && assetResponse.StatusCode > 299) {
		return false, fmt.Errorf("failed to download asset")
	}
	defer assetResponse.Body.Close()

	signatureResponse, err := http.Get(signatureDownloadUrl)
	if err != nil || (assetResponse.StatusCode < 200 && assetResponse.StatusCode > 299) {
		return false, fmt.Errorf("failed to download asset signature")
	}
	defer signatureResponse.Body.Close()

	assetBody := bytes.NewBuffer([]byte{})
	io.Copy(assetBody, assetResponse.Body)

	signatureBody := bytes.NewBuffer([]byte{})
	io.Copy(signatureBody, signatureResponse.Body)

	checksum, err := hex.DecodeString(assetBody.String())
	if err != nil {
		return false, fmt.Errorf("failed to parse asset signature")
	}

	updateOptions := selfupdate.Options{
		Hash:     crypto.MD5,
		Checksum: checksum,
	}

	if err := selfupdate.Apply(assetBody, updateOptions); err != nil {
		if err = selfupdate.RollbackError(err); err != nil {
			return false, fmt.Errorf("failed to rollback after unsuccessful update. %s", err.Error())
		}
	}

	Version = releasesResponse.TagName

	// successfully updated with no errors
	return true, nil
}
