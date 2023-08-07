package update

import (
	"archive/tar"
	"archive/zip"
	"bytes"
	"compress/gzip"
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"io"
	"net/http"
	"os"
	"runtime"
	"strconv"
	"strings"

	"github.com/tywil04/slavartdl/internal/helpers"
)

var Version = "v1.1.8"

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

	err := helpers.JsonApiRequest(
		http.MethodGet,
		"https://api.github.com/repos/tywil04/slavartdl/releases/latest",
		&releasesResponse,
		nil,
		nil,
	)
	if err != nil {
		return "", err
	}

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

	executablePath, err := os.Executable()
	if err != nil {
		return releasesResponse.TagName, fmt.Errorf("failed to find path to currently running executable")
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

			if header.Name == "slavartdl" {
				if err := os.Remove(executablePath); err != nil {
					return releasesResponse.TagName, fmt.Errorf("failed to remove currently running executable")
				}

				file, err := os.OpenFile(executablePath, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0755)
				if err != nil {
					fmt.Println(err)
					return releasesResponse.TagName, fmt.Errorf("failed to open currently running executable")
				}
				defer file.Close()

				if _, err := io.Copy(file, tarball); err != nil {
					fmt.Println(err)
					return releasesResponse.TagName, fmt.Errorf("failed to copy slavartdl from tarball into currently running executable")
				}
			}
		}
	} else {
		archive, err := zip.NewReader(bytes.NewReader(assetRaw), int64(len(assetRaw)))
		if err != nil {
			return releasesResponse.TagName, fmt.Errorf("failed to load asset into zip")
		}

		for _, zipped := range archive.File {
			if zipped.Name == "slavartdl" {
				if err := os.Remove(executablePath); err != nil {
					return releasesResponse.TagName, fmt.Errorf("failed to remove currently running executable")
				}

				file, err := os.OpenFile(executablePath, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0755)
				if err != nil {
					return releasesResponse.TagName, fmt.Errorf("failed to open currently running executable")
				}
				defer file.Close()

				zipFile, err := zipped.Open()
				if err != nil {
					return releasesResponse.TagName, fmt.Errorf("failed to open file in zip archive")
				}
				defer zipFile.Close()

				if _, err := io.Copy(file, zipFile); err != nil {
					return releasesResponse.TagName, fmt.Errorf("failed to copy slavartdl from zip into currently running archive")
				}

				break
			}
		}
	}

	// successfully updated with no errors
	return releasesResponse.TagName, nil
}
