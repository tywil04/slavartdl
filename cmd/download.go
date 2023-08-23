package cmd

import (
	"errors"
	"fmt"
	"math/rand"
	"net/url"
	"os"
	"path/filepath"
	"time"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"github.com/tywil04/slavartdl/internal/config"
	"github.com/tywil04/slavartdl/internal/helpers"
	"github.com/tywil04/slavartdl/internal/slavart"
)

var downloadCmd = &cobra.Command{
	Use:          "download [flags] url(s)",
	Short:        "Download music from url using SlavArt Divolt server",
	Long:         "Download music from url using SlavArt Divolt server (Supports: Tidal, Qobuz, SoundCloud, Deezer, Spotify, YouTube and Jiosaavn)",
	Args:         cobra.ArbitraryArgs,
	SilenceUsage: true,
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

		// optional
		configPathRel, err := flags.GetString("configPath")
		if err != nil {
			return fmt.Errorf("unknown error when getting '--configPath'")
		}

		configPath, err := filepath.Abs(configPathRel)
		if err != nil {
			return fmt.Errorf("failed to resolve relative 'configPath' into absolute path")
		}

		// load config
		if err := config.Load(configPathRel == "", configPath); err != nil {
			return err
		}

		// required
		outputDirRel, err := flags.GetString("outputDir")
		if err != nil {
			return fmt.Errorf("unknown error when getting '--outputDir'")
		}

		if outputDirRel == "" {
			outputDirRel = viper.GetString("downloadcmd.outputdir")
			if outputDirRel == "" {
				return fmt.Errorf("no outputDir provided in config or '--outputDir'")
			}
		}

		outputDir, err := filepath.Abs(outputDirRel)
		if err != nil {
			return fmt.Errorf("failed to resolve relative 'outputDir' into absolute path")
		}

		// optional
		fromFileRel, err := flags.GetString("fromFile")
		if err != nil {
			return fmt.Errorf("unknown error when getting '--fromFile'")
		}

		fromFile, err := filepath.Abs(fromFileRel)
		if err != nil {
			return fmt.Errorf("failed to resolve relative 'fromFile' into absolute path")
		}

		// optional
		fromStdin, err := flags.GetBool("fromStdin")
		if err != nil {
			return fmt.Errorf("unknown error when getting '--fromStdin'")
		}

		// optional
		quality, err := flags.GetInt("quality")
		if err != nil {
			return fmt.Errorf("unknown error when getting '--quality'")
		}

		if quality == 0 {
			quality = viper.GetInt("downloadcmd.quality")
		}

		// normalise quality to the same scale as the slavart bot. if quality is -1 it gets ignored later on
		quality -= 1

		// optional
		timeout, err := flags.GetInt("timeout")
		if err != nil {
			return fmt.Errorf("unknown error when getting '--timeout'")
		}

		if timeout == 0 {
			timeout = viper.GetInt("downloadcmd.timeout")
		}

		if timeout == 0 {
			return fmt.Errorf("total timeout is 0, unable to continue")
		}

		// optional
		cooldown, err := flags.GetInt("cooldown")
		if err != nil {
			return fmt.Errorf("unknown error when getting '--cooldown'")
		}

		if cooldown == 0 {
			cooldown = viper.GetInt("downloadcmd.cooldown")
		}

		// optional
		ignoreCover, err := flags.GetBool("ignoreCover")
		if err != nil {
			return fmt.Errorf("unknown error when getting '--ignoreCover'")
		}

		if !ignoreCover {
			ignoreCover = viper.GetBool("downloadcmd.ignore.cover")
		}

		// optional
		ignoreSubdirs, err := flags.GetBool("ignoreSubdirs")
		if err != nil {
			return fmt.Errorf("unknown error when getting '--ignoreSubdirs'")
		}

		if !ignoreSubdirs {
			ignoreSubdirs = viper.GetBool("downloadcmd.ignore.subdirs")
		}

		// optional
		skipUnzip, err := flags.GetBool("skipUnzip")
		if err != nil {
			return fmt.Errorf("unknown error when getting '--skipUnzip'")
		}

		if !skipUnzip {
			skipUnzip = viper.GetBool("downloadcmd.skip.unzip")
		}

		timeoutTime := time.Now().Add(time.Second * time.Duration(timeout))
		cooldownDuration := time.Second * time.Duration(cooldown)

		if fromFile != "" {
			// if a file is provided, add the urls to the list to be processed
			urls, err := helpers.GetUrlsFromFile(fromFile)
			if err != nil {
				return fmt.Errorf("failed to read urls from file")
			}
			args = append(args, urls...)
		}

		if fromStdin {
			// if told to read from stdin
			urls, err := helpers.GetUrlsFromStdin()
			if err != nil {
				return fmt.Errorf("failed to read urls from stdin")
			}
			args = append(args, urls...)
		}

		// this is now required because args dont have to be passed just via the cli anymore
		if len(args) == 0 {
			return fmt.Errorf("no urls to download")
		}

		for _, link := range args {
			// randomly select a session token to avoid using the same account all the time
			var sessionToken string
			sessionTokens := viper.GetStringSlice("divoltsessiontokens")
			loginCredentialsInterface := viper.Get("divoltlogincredentials")

			if loginCredentials, ok := loginCredentialsInterface.([]any); ok {
				for _, credentialAny := range loginCredentials {
					// if any issue is encountered skip these credentials

					credential, ok := credentialAny.(map[string]any)
					if !ok {
						continue
					}

					emailInterface, ok := credential["email"]
					if !ok {
						continue
					}

					email, ok := emailInterface.(string)
					if !ok {
						continue
					}

					passwordInterface, ok := credential["password"]
					if !ok {
						continue
					}

					password, ok := passwordInterface.(string)
					if !ok {
						continue
					}

					token, err := slavart.GetSessionTokenFromCredentials(email, password)
					if err != nil {
						continue
					}

					sessionTokens = append(sessionTokens, token)
				}
			}

			length := len(sessionTokens)
			if length == 0 {
				return fmt.Errorf("no session tokens found in config")
			} else if length == 1 {
				sessionToken = sessionTokens[0]
			} else {
				sessionToken = sessionTokens[rand.Intn(length)]
			}

			fmt.Println("Getting download link...")
			downloadLink, err := slavart.GetDownloadLinkFromSlavart(sessionToken, link, quality, timeoutTime)
			if err != nil {
				return err
			}

			fmt.Println("\nDownloading zip...")
			// this will create a temp file in the default location
			tempFile, err := os.CreateTemp("", "slavartdl.*.zip")
			if err != nil {
				return err
			}
			defer os.Remove(tempFile.Name())

			tempFilePath := tempFile.Name()
			err = helpers.DownloadFile(downloadLink, tempFilePath)
			if err != nil {
				return err
			}

			if !skipUnzip {
				fmt.Println("\nUnzipping...")
				if err := helpers.Unzip(tempFilePath, outputDir, ignoreSubdirs, ignoreCover); err != nil {
					return err
				}
			} else {
				zipName, err := helpers.GetZipName(tempFilePath)
				if err != nil {
					return err
				}

				fmt.Println(filepath.Clean(zipName))

				outputFileDir := outputDir + string(os.PathSeparator) + filepath.Clean(zipName) + ".zip"
				// temp file gets deleted later
				if err := helpers.CopyFile(tempFilePath, outputFileDir); err != nil {
					return nil
				}

			}

			fmt.Println("\nDone!")

			if link != args[len(args)-1] {
				time.Sleep(cooldownDuration)
			}
		}

		return nil
	},
}

func init() {
	flags := downloadCmd.Flags()

	flags.StringP("outputDir", "o", "", "the output directory to store the downloaded music")
	flags.StringP("fromFile", "f", "", "the path to a text file to read urls from, urls must be seperated by a newline")
	downloadCmd.MarkFlagDirname("outputDir")
	downloadCmd.MarkFlagDirname("fromFile")

	flags.BoolP("fromStdin", "s", false, "should urls be read from standard input, urls must be seperated by a newline")

	flags.StringP("configPath", "C", "", "a directory that contains an override config.json file\nor a file which contains an override config\n[a custom config file must end in .json]")

	flags.IntP("quality", "q", 0, "the quality of music to download\n- 0: best quality available\n- 1: 128kbps MP3/AAC\n- 2: 320kbps MP3/AAC\n- 3: 16bit 44.1kHz\n- 4: 24bit ≤96kHz\n- 5: 24bit ≤192kHz")
	flags.IntP("timeout", "t", 0, "how long before link search is timed out in seconds")
	flags.Int("cooldown", 0, "how long to wait after downloading first url in seconds\n(only matters if you are downloading multiple urls at once)")

	flags.BoolP("ignoreCover", "c", false, "ignore cover.jpg when unzipping downloaded music")
	flags.BoolP("ignoreSubdirs", "d", false, "ignore subdirectories when unzipping downloaded music")
	flags.BoolP("skipUnzip", "z", false, "skip unzipping downloaded music")

	rootCmd.AddCommand(downloadCmd)
}
