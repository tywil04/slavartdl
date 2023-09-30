package config

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
)

var Open *Config
var openFilePath string

const pathSeparator = string(os.PathSeparator)

func defaultConfig() *Config {
	return &Config{
		DivoltSessionTokens:     []string{},
		DivoltLoginCredentials:  []*ConfigCredential{},
		DiscordSessionTokens:    []string{},
		DiscordLoginCredentials: []*ConfigCredential{},
		DownloadCmd: &ConfigDownloadCmd{
			UseDiscord: false,
			OutputDir:  "",
			LogLevel:   "all",
			Quality:    0,
			Timeout:    120,
			Cooldown:   0,
			Ignore: &ConfigDownloadCmdIgnore{
				Cover:   false,
				SubDirs: false,
			},
			Skip: &ConfigDownloadCmdSkip{
				Unzip: false,
			},
		},
	}
}

func DefaultConfigLocation() (string, error) {
	userConfigDir, err := os.UserConfigDir()
	if err != nil {
		return "", fmt.Errorf("failed to find user config directory")
	}
	return userConfigDir + pathSeparator + "SlavartDL" + pathSeparator + "config.json", nil
}

func OpenConfig(customPath string) error {
	filePath := customPath
	if filePath == "" {
		var err error
		filePath, err = DefaultConfigLocation()
		if err != nil {
			return err
		}
	}

	fileStat, err := os.Stat(filePath)
	if err != nil && !os.IsNotExist(err) {
		return err
	}

	if err == nil && fileStat.IsDir() {
		filePath += pathSeparator + "/config.json"
	}

	dirPath := filepath.Dir(filePath)
	if err := os.MkdirAll(dirPath, os.ModeDir); err != nil && !os.IsExist(err) {
		return err
	}

	file, err := os.OpenFile(filePath, os.O_RDWR|os.O_CREATE, os.ModePerm)
	if err != nil {
		return err
	}

	config := &Config{}
	if err := json.NewDecoder(file).Decode(config); err != nil {
		if err.Error() == "EOF" {
			config = defaultConfig()
		} else {
			return err
		}
	}

	Open = config
	openFilePath = filePath
	return nil
}

func SaveConfig() error {
	file, err := os.OpenFile(openFilePath, os.O_RDWR|os.O_CREATE, os.ModePerm)
	if err != nil {
		return err
	}
	defer file.Close()

	jsonEncoder := json.NewEncoder(file)
	jsonEncoder.SetIndent("", "    ")
	if err := jsonEncoder.Encode(Open); err != nil {
		return err
	}

	return nil
}
