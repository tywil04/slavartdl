package cmd

import (
	"log"
	"math/rand"
	"net/url"
	"os"
	"path/filepath"
	"time"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"github.com/tywil04/slavartdl/cli/internal/config"
	"github.com/tywil04/slavartdl/cli/internal/helpers"
	"github.com/tywil04/slavartdl/discord"
	"github.com/tywil04/slavartdl/divolt"
	"github.com/tywil04/slavartdl/downloader"
)

const pathSeperator = string(os.PathSeparator)

var downloadCmd = &cobra.Command{
	Use:          "download [flags] url(s)",
	Short:        "Download music from url using SlavArt Divolt server",
	Long:         "Download music from url using SlavArt Divolt server (Supports: Tidal, Qobuz, SoundCloud, Deezer, Spotify, YouTube and Jiosaavn)",
	Args:         cobra.ArbitraryArgs,
	SilenceUsage: true,
	PreRun: func(cmd *cobra.Command, args []string) {
		for _, arg := range args {
			parsedUrl, err := url.ParseRequestURI(arg)
			if err != nil {
				log.Fatal(err.Error())
			}

			allowed := false
			for _, host := range divolt.SlavartAllowedHosts {
				if host == parsedUrl.Host {
					allowed = true
					break
				}
			}

			if !allowed {
				log.Fatal("host not allowed")
			}
		}
	},
	Run: func(cmd *cobra.Command, args []string) {
		flags := cmd.Flags()

		// optional
		configPathRel, err := flags.GetString("configPath")
		if err != nil {
			log.Fatal(err)
		}

		configPath, err := filepath.Abs(configPathRel)
		if err != nil {
			log.Fatal(err)
		}

		// load config
		if err := config.Load(configPathRel == "", configPath); err != nil {
			log.Fatal(err)
		}

		// optional
		logLevel, err := flags.GetString("logLevel")
		helpers.LogError(err, logLevel)

		if logLevel == "" {
			logLevel = viper.GetString("downloadcmd.loglevel")
		}

		if logLevel == "" {
			// set default
			logLevel = "all"
		}

		// required
		outputDirRel, err := flags.GetString("outputDir")
		helpers.LogError(err, logLevel)

		if outputDirRel == "" {
			outputDirRel = viper.GetString("downloadcmd.outputdir")
			if outputDirRel == "" {
				helpers.ManualLogError("no outputDir provided in config or '--outputDir'", logLevel)
			}
		}

		outputDir, err := filepath.Abs(outputDirRel)
		helpers.LogError(err, logLevel)

		// optional
		fromFile, err := flags.GetString("fromFile")
		helpers.LogError(err, logLevel)

		// optional
		fromStdin, err := flags.GetBool("fromStdin")
		helpers.LogError(err, logLevel)

		// optional
		quality, err := flags.GetInt("quality")
		helpers.LogError(err, logLevel)

		if quality == 0 {
			quality = viper.GetInt("downloadcmd.quality")
		}

		// normalise quality to the same scale as the slavart bot. if quality is -1 it gets ignored later on
		quality -= 1

		// optional
		timeout, err := flags.GetInt("timeout")
		helpers.LogError(err, logLevel)

		if timeout == 0 {
			timeout = viper.GetInt("downloadcmd.timeout")
		}

		if timeout == 0 {
			helpers.ManualLogError("total timeout is 0, unable to continue", logLevel)
		}

		// optional
		cooldown, err := flags.GetInt("cooldown")
		if err != nil {
			helpers.ManualLogError("unknown error when getting '--cooldown'", logLevel)
		}

		if cooldown == 0 {
			cooldown = viper.GetInt("downloadcmd.cooldown")
		}

		// optional
		ignoreCover, err := flags.GetBool("ignoreCover")
		helpers.LogError(err, logLevel)

		if !ignoreCover {
			ignoreCover = viper.GetBool("downloadcmd.ignore.cover")
		}

		// optional
		ignoreSubdirs, err := flags.GetBool("ignoreSubdirs")
		helpers.LogError(err, logLevel)

		if !ignoreSubdirs {
			ignoreSubdirs = viper.GetBool("downloadcmd.ignore.subdirs")
		}

		// optional
		skipUnzip, err := flags.GetBool("skipUnzip")
		helpers.LogError(err, logLevel)

		if !skipUnzip {
			skipUnzip = viper.GetBool("downloadcmd.skip.unzip")
		}

		timeoutTime := time.Second * time.Duration(timeout)
		cooldownDuration := time.Second * time.Duration(cooldown)

		useDiscord, err := flags.GetBool("useDiscord")
		helpers.LogError(err, logLevel)

		if fromFile != "" {
			// if a file is provided, add the urls to the list to be processed
			urls, err := helpers.GetUrlsFromFile(fromFile)
			helpers.LogError(err, logLevel)
			args = append(args, urls...)
		}

		if fromStdin {
			// if told to read from stdin
			urls, err := helpers.GetUrlsFromStdin()
			helpers.LogError(err, logLevel)
			args = append(args, urls...)
		}

		// this is now required because args dont have to be passed just via the cli anymore
		if len(args) == 0 {
			helpers.ManualLogError("no urls to download", logLevel)
		}

		if !useDiscord {
			// use divolt (default)

			sessionTokens := viper.GetStringSlice("divoltsessiontokens")
			loginCredentialsInterface := viper.Get("divoltlogincredentials")
			loginCredentials := loginCredentialsInterface.([]any)

			session := divolt.Session{}

			numberOfSessionTokens := len(sessionTokens)
			numberOfLoginCredentials := len(loginCredentials)

			randomlySelectedSource := -1
			if numberOfSessionTokens > 0 && numberOfLoginCredentials > 0 {
				randomlySelectedSource = rand.Intn(2)
			} else if numberOfSessionTokens > 0 && numberOfLoginCredentials == 0 {
				randomlySelectedSource = 1
			} else if numberOfSessionTokens == 0 && numberOfLoginCredentials > 0 {
				randomlySelectedSource = 0
			}

			switch randomlySelectedSource {
			case 0:
				var selectedCredential int
				if numberOfLoginCredentials == 1 {
					selectedCredential = 0
				} else {
					selectedCredential = rand.Intn(numberOfLoginCredentials)
				}

				credential := loginCredentials[selectedCredential].(map[string]string)
				err := session.AuthenticateWithCredentials(credential["email"], credential["password"])
				helpers.LogError(err, logLevel)
			case 1:
				var selectedToken int
				if numberOfSessionTokens == 1 {
					selectedToken = 0
				} else {
					selectedToken = rand.Intn(numberOfSessionTokens)
				}

				token := sessionTokens[selectedToken]
				err := session.AuthenticateWithSessionToken(token)
				helpers.LogError(err, logLevel)
			default:
				helpers.ManualLogError("no source to authenticated with divolt", logLevel)
			}

			for _, url := range args {
				status, err := session.SlavartGetBotStatus()
				helpers.LogError(err, logLevel)

				if status == divolt.SlavartBotStatusOffline {
					helpers.ManualLogError("slavart bot is offline", logLevel)
				}

				message, err := session.SlavartSendDownloadCommand(url, quality)
				helpers.LogError(err, logLevel)

				musicUrl, err := session.SlavartGetUploadUrl(message.Id, url, timeoutTime)
				helpers.LogError(err, logLevel)

				buffer, bytesWritten, err := downloader.DownloadFile(musicUrl)
				helpers.LogError(err, logLevel)

				if !skipUnzip {
					err := downloader.Unzip(buffer, bytesWritten, outputDir, ignoreSubdirs, ignoreCover)
					helpers.LogError(err, logLevel)
				} else {
					outputPath := outputDir + pathSeperator + filepath.Clean("slavart-"+time.Now().String()) + ".zip"
					err := downloader.CopyFile(buffer, outputPath)
					helpers.LogError(err, logLevel)
				}

				if url != args[len(args)-1] {
					time.Sleep(cooldownDuration)
				}
			}
		} else {
			// use discord

			sessionTokens := viper.GetStringSlice("discordsessiontokens")
			loginCredentialsInterface := viper.Get("discordlogincredentials")
			loginCredentials := loginCredentialsInterface.([]any)

			session := discord.Session{}

			numberOfSessionTokens := len(sessionTokens)
			numberOfLoginCredentials := len(loginCredentials)

			randomlySelectedSource := -1
			if numberOfSessionTokens > 0 && numberOfLoginCredentials > 0 {
				randomlySelectedSource = rand.Intn(2)
			} else if numberOfSessionTokens > 0 && numberOfLoginCredentials == 0 {
				randomlySelectedSource = 1
			} else if numberOfSessionTokens == 0 && numberOfLoginCredentials > 0 {
				randomlySelectedSource = 0
			}

			switch randomlySelectedSource {
			case 0:
				var selectedCredential int
				if numberOfLoginCredentials == 1 {
					selectedCredential = 0
				} else {
					selectedCredential = rand.Intn(numberOfLoginCredentials)
				}

				credential := loginCredentials[selectedCredential].(map[string]string)
				err := session.AuthenticateWithCredentials(credential["email"], credential["password"])
				helpers.LogError(err, logLevel)
			case 1:
				var selectedToken int
				if numberOfSessionTokens == 1 {
					selectedToken = 0
				} else {
					selectedToken = rand.Intn(numberOfSessionTokens)
				}

				token := sessionTokens[selectedToken]
				err := session.AuthenticateWithAuthorizationToken(token)
				helpers.LogError(err, logLevel)
			default:
				helpers.ManualLogError("no source to authenticated with discord", logLevel)
			}

			for _, url := range args {
				message, err := session.PixeldrainSendDownloadCommand(url, quality)
				helpers.LogError(err, logLevel)

				musicUrl, err := session.PixeldrainGetUploadUrl(message.Id, url, timeoutTime)
				helpers.LogError(err, logLevel)

				buffer, bytesWritten, err := downloader.DownloadFile(musicUrl)
				helpers.LogError(err, logLevel)

				if !skipUnzip {
					err := downloader.Unzip(buffer, bytesWritten, outputDir, ignoreSubdirs, ignoreCover)
					helpers.LogError(err, logLevel)
				} else {
					outputPath := outputDir + pathSeperator + filepath.Clean("slavart-"+time.Now().String()) + ".zip"
					err := downloader.CopyFile(buffer, outputPath)
					helpers.LogError(err, logLevel)
				}

				if url != args[len(args)-1] {
					time.Sleep(cooldownDuration)
				}
			}
		}
	},
}

func init() {
	flags := downloadCmd.Flags()

	flags.StringP("outputDir", "o", "", "the output directory to store the downloaded music")
	flags.StringP("fromFile", "f", "", "the path to a text file to read urls from, urls must be separated by a newline")
	downloadCmd.MarkFlagDirname("outputDir")
	downloadCmd.MarkFlagDirname("fromFile")

	flags.BoolP("fromStdin", "s", false, "should urls be read from standard input, urls must be separated by a newline")

	flags.StringP("configPath", "C", "", "a directory that contains an override config.json file\nor a file which contains an override config\n[a custom config file must end in .json]")

	flags.StringP("logLevel", "l", "", "what level of logs should be outputted to standard output\n(errors in regard with config file will always be reported)\n- all: everything gets logged\n- errors: only errors get logged\n- silent: nothing gets logged")
	flags.IntP("quality", "q", 0, "the quality of music to download\n- 0: best quality available\n- 1: 128kbps MP3/AAC\n- 2: 320kbps MP3/AAC\n- 3: 16bit 44.1kHz\n- 4: 24bit ≤96kHz\n- 5: 24bit ≤192kHz")
	flags.IntP("timeout", "t", 0, "how long before link search is timed out in seconds")
	flags.Int("cooldown", 0, "how long to wait after downloading first url in seconds\n(only matters if you are downloading multiple urls at once)")

	flags.BoolP("ignoreCover", "c", false, "ignore cover.jpg when unzipping downloaded music")
	flags.BoolP("ignoreSubdirs", "d", false, "ignore subdirectories when unzipping downloaded music")
	flags.BoolP("skipUnzip", "z", false, "skip unzipping downloaded music")
	flags.BoolP("useDiscord", "D", false, "use of discord instead of divolt, requires discord session token being installed in config")

	rootCmd.AddCommand(downloadCmd)
}
