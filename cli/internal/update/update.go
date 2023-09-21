package update

import (
	"archive/tar"
	"archive/zip"
	"bytes"
	"compress/gzip"
	"crypto/md5"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"runtime"
	"strconv"
	"strings"

	"github.com/minio/selfupdate"
)

const Version = "v1.1.11"

const assetNameTemplate = "slavartdl-%s-%s-%s.%s"
const signatureNameTemplate = "slavartdl-%s-%s-%s.%s.md5"

// accepts version like 'v0.0.0'.
// returns major, minor, patch ints
func parseVersionTag(version string) (int, int, int, error) {
	versionTrim := strings.TrimPrefix(version, "v")
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
func Update(force bool) (string, error) {
	releasesResponse := struct {
		TagName string `json:"tag_name"`
		Assets  []struct {
			BrowserDownloadUrl string `json:"browser_download_url"`
			Name               string `json:"name"`
		} `json:"assets"`
	}{}

	response, err := http.Get("https://api.github.com/repos/tywil04/slavartdl/releases/latest")
	if err != nil {
		return "", err
	}
	json.NewDecoder(response.Body).Decode(&releasesResponse)

	apiMajor, apiMinor, apiPatch, err := parseVersionTag(releasesResponse.TagName)
	if err != nil {
		return "", err
	}

	pkgMajor, pkgMinor, pkgPatch, err := parseVersionTag(Version)
	if err != nil {
		return "", err
	}

	// check if current version is equal to or greater than fetched version
	if !force && (apiMajor <= pkgMajor && apiMinor <= pkgMinor && apiPatch <= pkgPatch) {
		// no update, no error
		return "", nil
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
		return releasesResponse.TagName, fmt.Errorf("failed to get asset and asset signature download url for your system, supported platforms are linux [arm64 amd64], darwin [arm64 amd64], windows [amd64]")
	}

	assetResponse, err := http.Get(assetDownloadUrl)
	if err != nil || (assetResponse.StatusCode < 200 && assetResponse.StatusCode > 299) {
		return releasesResponse.TagName, fmt.Errorf("failed to download asset")
	}
	defer assetResponse.Body.Close()

	signatureResponse, err := http.Get(signatureDownloadUrl)
	if err != nil || (assetResponse.StatusCode < 200 && assetResponse.StatusCode > 299) {
		return releasesResponse.TagName, fmt.Errorf("failed to download asset signature")
	}
	defer signatureResponse.Body.Close()

	assetBuffer := bytes.NewBuffer([]byte{})
	if _, err := io.Copy(assetBuffer, assetResponse.Body); err != nil {
		return releasesResponse.TagName, fmt.Errorf("failed to copy assetResponse.body to assetBuffer")
	}
	assetRaw := assetBuffer.Bytes()

	signatureRaw, err := io.ReadAll(signatureResponse.Body)
	if err != nil {
		return releasesResponse.TagName, fmt.Errorf("failed to parse signature from signature response")
	}
	signature := strings.ReplaceAll(string(signatureRaw), "\n", "")

	downloadedFileSignatureRaw := md5.Sum(assetRaw)
	downloadedFileSignature := hex.EncodeToString(downloadedFileSignatureRaw[:])

	if signature != downloadedFileSignature {
		return releasesResponse.TagName, fmt.Errorf("downloaded file doesnt match signature")
	}

	if extension == "tar.gz" {
		archive, err := gzip.NewReader(assetBuffer)
		if err != nil {
			return releasesResponse.TagName, fmt.Errorf("failed to load asset into gzip reader")
		}
		defer archive.Close()

		tarball := tar.NewReader(archive)
		for {
			header, err := tarball.Next()
			if err == io.EOF {
				break
			} else if err != nil {
				continue
			}

			if strings.Contains(header.Name, "slavartdl") {
				if err := selfupdate.Apply(tarball, selfupdate.Options{}); err != nil {
					if err := selfupdate.RollbackError(err); err != nil {
						return "", fmt.Errorf("an unknown error has occured while updating")
					}
				}

				break
			}
		}
	} else {
		archive, err := zip.NewReader(bytes.NewReader(assetRaw), int64(len(assetRaw)))
		if err != nil {
			return releasesResponse.TagName, fmt.Errorf("failed to load asset into zip")
		}

		for _, zipped := range archive.File {
			if strings.Contains(zipped.Name, "slavartdl") {
				zipFile, err := zipped.Open()
				if err != nil {
					return releasesResponse.TagName, fmt.Errorf("failed to open file in zip archive")
				}
				defer zipFile.Close()

				if err := selfupdate.Apply(zipFile, selfupdate.Options{}); err != nil {
					if err := selfupdate.RollbackError(err); err != nil {
						return "", fmt.Errorf("an unknown error has occured while updating")
					}
				}

				break
			}
		}
	}

	// successfully updated with no errors
	return releasesResponse.TagName, nil
}
