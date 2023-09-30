package cmd

import (
	"log"
	"math/rand"
	"net/url"
	"os"
	"path/filepath"
	"time"

	"github.com/spf13/cobra"

	"github.com/tywil04/slavartdl/cli/internal/config"
	"github.com/tywil04/slavartdl/cli/internal/helpers"
	"github.com/tywil04/slavartdl/discord"
	"github.com/tywil04/slavartdl/divolt"
	"github.com/tywil04/slavartdl/downloader"
)

const pathSeparator = string(os.PathSeparator)

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
		configPath, err := flags.GetString("configPath")
		if err != nil {
			log.Fatal(err)
		}

		// load config
		if err := config.OpenConfig(configPath); err != nil {
			log.Fatal(err)
		}

		// optional
		logLevel, err := flags.GetString("logLevel")
		helpers.LogError(err, logLevel)

		if logLevel == "" {
			logLevel = config.Open.DownloadCmd.LogLevel
		}

		// required
		outputDirRel, err := flags.GetString("outputDir")
		helpers.LogError(err, logLevel)

		if outputDirRel == "" {
			outputDirRel = config.Open.DownloadCmd.OutputDir
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
			quality = config.Open.DownloadCmd.Quality
		}

		// normalise quality to the same scale as the slavart bot. if quality is -1 it gets ignored later on
		quality -= 1

		// optional
		timeout, err := flags.GetInt("timeout")
		helpers.LogError(err, logLevel)

		if timeout == 0 {
			timeout = config.Open.DownloadCmd.Timeout
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
			cooldown = config.Open.DownloadCmd.Cooldown
		}

		// optional
		ignoreCover, err := flags.GetBool("ignoreCover")
		helpers.LogError(err, logLevel)

		if !ignoreCover {
			ignoreCover = config.Open.DownloadCmd.Ignore.Cover
		}

		// optional
		ignoreSubdirs, err := flags.GetBool("ignoreSubdirs")
		helpers.LogError(err, logLevel)

		if !ignoreSubdirs {
			ignoreSubdirs = config.Open.DownloadCmd.Ignore.SubDirs
		}

		// optional
		skipUnzip, err := flags.GetBool("skipUnzip")
		helpers.LogError(err, logLevel)

		if !skipUnzip {
			skipUnzip = config.Open.DownloadCmd.Skip.Unzip
		}

		timeoutTime := time.Second * time.Duration(timeout)
		cooldownDuration := time.Second * time.Duration(cooldown)

		// optional
		useDiscord, err := flags.GetBool("useDiscord")
		helpers.LogError(err, logLevel)

		if !useDiscord {
			useDiscord = config.Open.DownloadCmd.UseDiscord
		}

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

			session := divolt.Session{}

			numberOfSessionTokens := len(config.Open.DivoltSessionTokens)
			numberOfLoginCredentials := len(config.Open.DivoltLoginCredentials)

			randomlySelectedSource := -1
			if numberOfSessionTokens > 0 && numberOfLoginCredentials > 0 {
				randomlySelectedSource = rand.Intn(2)
			} else if numberOfSessionTokens > 0 && numberOfLoginCredentials == 0 {
				randomlySelectedSource = 1
			} else if numberOfSessionTokens == 0 && numberOfLoginCredentials > 0 {
				randomlySelectedSource = 0
			}

			helpers.Println("[DIVOLT]: Starting authentication...", logLevel)

			switch randomlySelectedSource {
			case 0:
				var selectedCredential int
				if numberOfLoginCredentials == 1 {
					selectedCredential = 0
				} else {
					selectedCredential = rand.Intn(numberOfLoginCredentials)
				}

				credential := config.Open.DivoltLoginCredentials[selectedCredential]
				err := session.AuthenticateWithCredentials(credential.Email, credential.Password)
				helpers.LogError(err, logLevel)
			case 1:
				var selectedToken int
				if numberOfSessionTokens == 1 {
					selectedToken = 0
				} else {
					selectedToken = rand.Intn(numberOfSessionTokens)
				}

				token := config.Open.DivoltSessionTokens[selectedToken]
				err := session.AuthenticateWithSessionToken(token)
				helpers.LogError(err, logLevel)
			default:
				helpers.ManualLogError("no source to authenticated with divolt", logLevel)
			}

			helpers.Println("[DIVOLT]: Successfully authenticated!", logLevel)

			session.SlavartTryInviteUser()

			for _, url := range args {
				status, err := session.SlavartGetBotStatus()
				helpers.LogError(err, logLevel)

				if status == divolt.SlavartBotStatusOffline {
					helpers.ManualLogError("slavart bot is offline", logLevel)
				}

				message, err := session.SlavartSendDownloadCommand(url, quality)
				helpers.LogError(err, logLevel)

				helpers.Println("[DIVOLT]: Sent download command for "+url+".", logLevel)
				helpers.Println("[DIVOLT]: Waiting for download url...", logLevel)

				musicUrl, err := session.SlavartGetUploadUrl(message.Id, url, timeoutTime)
				helpers.LogError(err, logLevel)

				helpers.Println("[DIVOLT]: Successfully fetched download url!", logLevel)
				helpers.Println("[DIVOLT]: Starting download...", logLevel)

				buffer, bytesWritten, err := downloader.DownloadFile(musicUrl)
				helpers.LogError(err, logLevel)

				helpers.Println("[DIVOLT]: Successfully downloaded music archive!", logLevel)

				if !skipUnzip {
					helpers.Println("[DIVOLT]: Starting unzip...", logLevel)

					err := downloader.Unzip(buffer, bytesWritten, outputDir, ignoreSubdirs, ignoreCover)
					helpers.LogError(err, logLevel)

					helpers.Println("[DIVOLT]: Successfully unzipped music archive into download location!", logLevel)
				} else {
					helpers.Println("[DIVOLT]: Starting copy...", logLevel)

					outputPath := outputDir + pathSeparator + filepath.Clean("slavart-"+time.Now().String()) + ".zip"
					err := downloader.CopyFile(buffer, outputPath)
					helpers.LogError(err, logLevel)

					helpers.Println("[DIVOLT]: Successfully copied music archive into download location!", logLevel)
				}

				helpers.Println("[DIVOLT]: Successfully downloaded "+url+".", logLevel)

				if url != args[len(args)-1] {
					time.Sleep(cooldownDuration)
				}
			}
		} else {
			// use discord

			session := discord.Session{}

			numberOfSessionTokens := len(config.Open.DiscordSessionTokens)
			numberOfLoginCredentials := len(config.Open.DiscordLoginCredentials)

			randomlySelectedSource := -1
			if numberOfSessionTokens > 0 && numberOfLoginCredentials > 0 {
				randomlySelectedSource = rand.Intn(2)
			} else if numberOfSessionTokens > 0 && numberOfLoginCredentials == 0 {
				randomlySelectedSource = 1
			} else if numberOfSessionTokens == 0 && numberOfLoginCredentials > 0 {
				randomlySelectedSource = 0
			}

			helpers.Println("[DISCORD]: Starting authentication...", logLevel)

			switch randomlySelectedSource {
			case 0:
				var selectedCredential int
				if numberOfLoginCredentials == 1 {
					selectedCredential = 0
				} else {
					selectedCredential = rand.Intn(numberOfLoginCredentials)
				}

				credential := config.Open.DiscordLoginCredentials[selectedCredential]
				err := session.AuthenticateWithCredentials(credential.Email, credential.Password)
				helpers.LogError(err, logLevel)
			case 1:
				var selectedToken int
				if numberOfSessionTokens == 1 {
					selectedToken = 0
				} else {
					selectedToken = rand.Intn(numberOfSessionTokens)
				}

				token := config.Open.DiscordSessionTokens[selectedToken]
				err := session.AuthenticateWithAuthorizationToken(token)
				helpers.LogError(err, logLevel)
			default:
				helpers.ManualLogError("no source to authenticated with discord", logLevel)
			}

			helpers.Println("[DISCORD]: Successfully authenticated!", logLevel)

			session.PixeldrainTryInviteUser()

			for _, url := range args {
				message, err := session.PixeldrainSendDownloadCommand(url, quality)
				helpers.LogError(err, logLevel)

				helpers.Println("[DISCORD]: Sent download command for "+url+".", logLevel)
				helpers.Println("[DISCORD]: Waiting for download url...", logLevel)

				musicUrl, err := session.PixeldrainGetUploadUrl(message.Id, url, timeoutTime)
				helpers.LogError(err, logLevel)

				helpers.Println("[DISCORD]: Successfully fetched download url!", logLevel)
				helpers.Println("[DISCORD]: Starting download...", logLevel)

				buffer, bytesWritten, err := downloader.DownloadFile(musicUrl)
				helpers.LogError(err, logLevel)

				helpers.Println("[DISCORD]: Successfully downloaded music archive!", logLevel)

				if !skipUnzip {
					helpers.Println("[DISCORD]: Starting unzip...", logLevel)

					err := downloader.Unzip(buffer, bytesWritten, outputDir, ignoreSubdirs, ignoreCover)
					helpers.LogError(err, logLevel)

					helpers.Println("[DISCORD]: Successfully unzipped music archive into download location!", logLevel)
				} else {
					helpers.Println("[DISCORD]: Starting copy...", logLevel)

					outputPath := outputDir + pathSeparator + filepath.Clean("slavart-"+time.Now().String()) + ".zip"
					err := downloader.CopyFile(buffer, outputPath)
					helpers.LogError(err, logLevel)

					helpers.Println("[DISCORD]: Successfully copied music archive into download location!", logLevel)
				}

				helpers.Println("[DISCORD]: Successfully downloaded "+url+".", logLevel)

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
