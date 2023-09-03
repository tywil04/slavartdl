package cmd

import (
	"errors"
	"log"
	"math/rand"
	"net/url"
	"path/filepath"
	"time"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"github.com/tywil04/slavartdl/common"
	"github.com/tywil04/slavartdl/slavart"

	"github.com/tywil04/slavartdl/cli/internal/config"
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
	Run: func(cmd *cobra.Command, args []string) {
		flags := cmd.Flags()

		// optional
		configPathRel, err := flags.GetString("configPath")
		if err != nil {
			log.Fatal("unknown error when getting '--configPath'")
		}

		configPath, err := filepath.Abs(configPathRel)
		if err != nil {
			log.Fatal("failed to resolve relative 'configPath' into absolute path")
		}

		// load config
		if err := config.Load(configPathRel == "", configPath); err != nil {
			log.Fatal(err)
		}

		// optional
		logLevel, err := flags.GetString("logLevel")
		if err != nil {
			if logLevel == "all" || logLevel == "errors" {
				log.Fatal("unknown error when getting '--logLevel'")
			}
		}

		if logLevel == "" {
			logLevel = viper.GetString("downloadcmd.loglevel")
		}

		if logLevel == "" {
			// set default
			logLevel = "all"
		}

		if logLevel != "all" && logLevel != "errors" && logLevel != "silent" {
			if logLevel == "all" || logLevel == "errors" {
				log.Fatal("'--logLevel' should be one of 'all', 'errors' or 'silent'")
			}
		}

		// required
		outputDirRel, err := flags.GetString("outputDir")
		if err != nil {
			if logLevel == "all" || logLevel == "errors" {
				log.Fatal("unknown error when getting '--outputDir'")
			}
		}

		if outputDirRel == "" {
			outputDirRel = viper.GetString("downloadcmd.outputdir")
			if outputDirRel == "" {
				if logLevel == "all" || logLevel == "errors" {
					log.Fatal("no outputDir provided in config or '--outputDir'")
				}
			}
		}

		outputDir, err := filepath.Abs(outputDirRel)
		if err != nil {
			if logLevel == "all" || logLevel == "errors" {
				log.Fatal("failed to resolve relative 'outputDir' into absolute path")
			}
		}

		// optional
		fromFile, err := flags.GetString("fromFile")
		if err != nil {
			if logLevel == "all" || logLevel == "errors" {
				log.Fatal("unknown error when getting '--fromFile'")
			}
		}

		// optional
		fromStdin, err := flags.GetBool("fromStdin")
		if err != nil {
			if logLevel == "all" || logLevel == "errors" {
				log.Fatal("unknown error when getting '--fromStdin'")
			}
		}

		// optional
		quality, err := flags.GetInt("quality")
		if err != nil {
			if logLevel == "all" || logLevel == "errors" {
				log.Fatal("unknown error when getting '--quality'")
			}
		}

		if quality == 0 {
			quality = viper.GetInt("downloadcmd.quality")
		}

		// normalise quality to the same scale as the slavart bot. if quality is -1 it gets ignored later on
		quality -= 1

		// optional
		timeout, err := flags.GetInt("timeout")
		if err != nil {
			if logLevel == "all" || logLevel == "errors" {
				log.Fatal("unknown error when getting '--timeout'")
			}
		}

		if timeout == 0 {
			timeout = viper.GetInt("downloadcmd.timeout")
		}

		if timeout == 0 {
			if logLevel == "all" || logLevel == "errors" {
				log.Fatal("total timeout is 0, unable to continue")
			}
		}

		// optional
		cooldown, err := flags.GetInt("cooldown")
		if err != nil {
			if logLevel == "all" || logLevel == "errors" {
				log.Fatal("unknown error when getting '--cooldown'")
			}
		}

		if cooldown == 0 {
			cooldown = viper.GetInt("downloadcmd.cooldown")
		}

		// optional
		ignoreCover, err := flags.GetBool("ignoreCover")
		if err != nil {
			if logLevel == "all" || logLevel == "errors" {
				log.Fatal("unknown error when getting '--ignoreCover'")
			}
		}

		if !ignoreCover {
			ignoreCover = viper.GetBool("downloadcmd.ignore.cover")
		}

		// optional
		ignoreSubdirs, err := flags.GetBool("ignoreSubdirs")
		if err != nil {
			if logLevel == "all" || logLevel == "errors" {
				log.Fatal("unknown error when getting '--ignoreSubdirs'")
			}
		}

		if !ignoreSubdirs {
			ignoreSubdirs = viper.GetBool("downloadcmd.ignore.subdirs")
		}

		// optional
		skipUnzip, err := flags.GetBool("skipUnzip")
		if err != nil {
			if logLevel == "all" || logLevel == "errors" {
				log.Fatal("unknown error when getting '--skipUnzip'")
			}
		}

		if !skipUnzip {
			skipUnzip = viper.GetBool("downloadcmd.skip.unzip")
		}

		timeoutTime := time.Now().Add(time.Second * time.Duration(timeout))
		cooldownDuration := time.Second * time.Duration(cooldown)

		if fromFile != "" {
			// if a file is provided, add the urls to the list to be processed
			urls, err := common.GetUrlsFromFile(fromFile)
			if err != nil {
				if logLevel == "all" || logLevel == "errors" {
					log.Fatal("failed to read urls from file")
				}
			}
			args = append(args, urls...)
		}

		if fromStdin {
			// if told to read from stdin
			urls, err := common.GetUrlsFromStdin()
			if err != nil {
				if logLevel == "all" || logLevel == "errors" {
					log.Fatal("failed to read urls from stdin")
				}
			}
			args = append(args, urls...)
		}

		// this is now required because args dont have to be passed just via the cli anymore
		if len(args) == 0 {
			if logLevel == "all" || logLevel == "errors" {
				log.Fatal("no urls to download")
			}
		}

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
			if logLevel == "all" || logLevel == "errors" {
				log.Fatal("no session tokens found in config")
			}
		} else if length == 1 {
			sessionToken = sessionTokens[0]
		} else {
			sessionToken = sessionTokens[rand.Intn(length)]
		}

		slavart.Download(args, sessionToken, logLevel, quality, timeoutTime, cooldownDuration, outputDir, skipUnzip, ignoreCover, ignoreSubdirs)
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

	flags.StringP("logLevel", "l", "", "what level of logs should be outputted to standard output\n(errors in regard with config file will always be reported)\n- all: everything gets logged\n- errors: only errors get logged\n- silent: nothing gets logged")
	flags.IntP("quality", "q", 0, "the quality of music to download\n- 0: best quality available\n- 1: 128kbps MP3/AAC\n- 2: 320kbps MP3/AAC\n- 3: 16bit 44.1kHz\n- 4: 24bit ≤96kHz\n- 5: 24bit ≤192kHz")
	flags.IntP("timeout", "t", 0, "how long before link search is timed out in seconds")
	flags.Int("cooldown", 0, "how long to wait after downloading first url in seconds\n(only matters if you are downloading multiple urls at once)")

	flags.BoolP("ignoreCover", "c", false, "ignore cover.jpg when unzipping downloaded music")
	flags.BoolP("ignoreSubdirs", "d", false, "ignore subdirectories when unzipping downloaded music")
	flags.BoolP("skipUnzip", "z", false, "skip unzipping downloaded music")

	rootCmd.AddCommand(downloadCmd)
}
