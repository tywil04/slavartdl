package cmd

import (
	"errors"
	"fmt"
	"net/url"
	"os"
	"path/filepath"
	"strconv"
	"time"

	"github.com/spf13/cobra"

	"github.com/tywil04/slavartdl/internal/helpers"
	"github.com/tywil04/slavartdl/internal/slavart"
)

var downloadCmd = &cobra.Command{
	Use:   "download [flags] url(s)",
	Short: "Download music from url using SlavArt Divolt server",
	Long:  "Download music from url using SlavArt Divolt server (Supports: Tidal, Qobuz, SoundCloud, Deezer, Spotify, YouTube and Jiosaavn)",
	Args:  cobra.MinimumNArgs(1),
	PreRunE: func(cmd *cobra.Command, args []string) error {
		for _, arg := range args {
			parsedUrl, err := url.ParseRequestURI(arg)
			if err != nil {
				return err
			}

			allowed := false
			for _, host := range slavart.AllowedHosts {
				if host == parsedUrl.Host {
					allowed = true
					break
				}
			}

			if !allowed {
				return errors.New("host not allowed")
			}
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

		// optional
		timeoutDurationSeconds, err := strconv.Atoi(flags.Lookup("timeout-duration-seconds").Value.String())
		if err != nil {
			return err
		}

		// optional
		timeoutDurationMinutes, err := strconv.Atoi(flags.Lookup("timeout-duration-minutes").Value.String())
		if err != nil {
			return err
		}

		// optional
		ignoreCover, err := strconv.ParseBool(flags.Lookup("ignore-cover").Value.String())
		if err != nil {
			return err
		}

		// optional
		ignoreSubdirectories, err := strconv.ParseBool(flags.Lookup("ignore-subdirectories").Value.String())
		if err != nil {
			return err
		}

		timeoutTime := time.Now().
			Add(time.Minute * time.Duration(timeoutDurationMinutes)).
			Add(time.Second * time.Duration(timeoutDurationSeconds))

		for _, arg := range args {
			fmt.Println("Getting download link...")
			downloadLink, err := slavart.GetDownloadLinkFromSlavart(arg, quality, timeoutTime)
			if err != nil {
				return err
			}

			fmt.Println("\nDownloading zip...")
			// this will create a temp file in the default location
			tempFile, err := os.CreateTemp("", "slavartdownloader.*.zip")
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
			err = helpers.Unzip(tempFilePath, outputDirectory, ignoreSubdirectories, ignoreCover)
			if err != nil {
				return err
			}

			fmt.Println("\nDone!")
		}

		return nil
	},
}

func init() {
	flags := downloadCmd.Flags()

	flags.StringP("outputDirs", "o", "", "the output directory to store the downloaded music")
	downloadCmd.MarkFlagRequired("output-directory")
	downloadCmd.MarkFlagDirname("output-directory")

	flags.IntP("quality", "q", 0, "the quality of music to download, if omited best quality\n- 1: 128kbps MP3/AAC\n- 2: 320kbps MP3/AAC\n- 3: 16bit 44.1kHz\n- 4: 24bit ≤96kHz\n- 5: 24bit ≤192kHz")
	flags.IntP("timeoutSeconds", "s", 0, "how long before link search is timed out in seconds\n[combines with --timeoutMinutes]")
	flags.IntP("timeoutMinutes", "m", 0, "how long before link search is timed out in minutes\n[combines with --timeoutSeconds]")

	flags.BoolP("ignoreCover", "c", false, "ignore cover.jpg when unzipping downloaded music")
	flags.BoolP("ignoreSubdirs", "d", false, "ignore subdirectories when unzipping downloaded music")

	rootCmd.AddCommand(downloadCmd)
}
