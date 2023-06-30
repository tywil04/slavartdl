package cmd

import (
	"errors"
	"fmt"
	"net/url"
	"os"
	"path/filepath"
	"strconv"

	"github.com/spf13/cobra"

	"slavartdl/lib/helpers"
	"slavartdl/lib/slavart"
)

var downloadCmd = &cobra.Command{
	Use:       "download url [flags]",
	Short:     "download music from url using slavart (supports: tidal, qobuz, soundcloud, deezer, spotify, youtube and jiosaavn)",
	Long:      "download music from url using slavart (supports: x, y, z)",
	Args:      cobra.ExactArgs(1),
	ValidArgs: []string{"url"},
	PreRunE: func(cmd *cobra.Command, args []string) error {
		parsedUrl, err := url.ParseRequestURI(args[0])
		if err != nil {
			return err
		}

		allowedHosts := []string{
			"tidal.com",
			"www.qobuz.com",
			"soundcloud.com",
			"www.deezer.com",
			"open.spotify.com",
			"music.youtube.com",
			"www.jiosaavn.com",
		}

		allowed := false
		for _, host := range allowedHosts {
			if host == parsedUrl.Host {
				allowed = true
				break
			}
		}

		if !allowed {
			return errors.New("host not allowed")
		}

		return nil
	},
	RunE: func(cmd *cobra.Command, args []string) error {
		flags := cmd.Flags()

		// required
		outputDirectory, err := filepath.Abs(flags.Lookup("output-directory").Value.String())
		if err != nil {
			return err
		}

		// optional
		quality, err := strconv.Atoi(flags.Lookup("quality").Value.String())
		if err != nil {
			return err
		}

		fmt.Println("Getting download link...")
		downloadLink, err := slavart.GetDownloadLinkFromSlavart(args[0], quality)
		if err != nil {
			return err
		}

		fmt.Println("\nDownloading zip...")
		tempFile, err := os.CreateTemp("/tmp", "slavartdownloader.*.zip")
		if err != nil {
			return err
		}
		defer os.Remove(tempFile.Name())

		tempFilePath := tempFile.Name()
		err = helpers.DownloadFile(downloadLink, tempFilePath)
		if err != nil {
			return err
		}

		fmt.Println("\nUnzipping...")
		err = helpers.Unzip(tempFilePath, outputDirectory)
		if err != nil {
			return err
		}

		fmt.Println("\nDone!")

		return nil
	},
}

func init() {
	flags := downloadCmd.Flags()

	flags.StringP("output-directory", "o", "", "the output directory to store the downloaded music")
	downloadCmd.MarkFlagRequired("output-directory")
	downloadCmd.MarkFlagDirname("output-directory")

	flags.IntP("quality", "q", -1, "the quality of music to download, omit (or -1) for best quality available (1: 128kbps MP3/AAC, 2: 320kbps MP3/AAC, 3: 16bit 44.1kHz, 4: 24bit ≤96kHz, 5: 24bit ≤192kHz)")

	rootCmd.AddCommand(downloadCmd)
}
