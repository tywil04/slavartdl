package config

import (
	"bytes"
	"fmt"
	"os"
	"path/filepath"

	"github.com/spf13/viper"
)

var ConfigDir string
var ConfigName string
var ConfigFile string

func Load(defaultHandling bool, customPath string) error {
	// default handling means it uses the default config
	if defaultHandling {
		var err error
		ConfigDir, err = os.UserConfigDir()
		if err != nil {
			return fmt.Errorf("failed to find user config directory")
		}

		ConfigName = "config.json"
		ConfigDir += string(os.PathSeparator) + "SlavartDL"
		ConfigFile = ConfigDir + string(os.PathSeparator) + ConfigName
	} else {
		fileInfo, err := os.Stat(customPath)
		if err == nil && fileInfo.IsDir() {
			// directory exists
			ConfigName = "config.json"
			ConfigDir = customPath
			ConfigFile = ConfigDir + string(os.PathSeparator) + ConfigName
		} else if err == nil || (err != nil && os.IsNotExist(err)) {
			// either file does exist, or it doesnt
			// if it doesnt exist viper will create it
			dirName, fileName := filepath.Split(customPath)
			ext := filepath.Ext(fileName)
			if ext != ".json" {
				return fmt.Errorf("custom config file must end in .json")
			}

			ConfigName = fileName
			ConfigDir = dirName
			ConfigFile = customPath
		} else {
			return fmt.Errorf("unknown error when handling custom config")
		}
	}

	// configure viper
	viper.SetConfigName(ConfigName)
	viper.SetConfigType("json")
	viper.AddConfigPath(ConfigDir)

	if err := viper.ReadInConfig(); err != nil {
		// file not found
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			if err := os.Mkdir(ConfigDir, os.ModePerm); err != nil && !os.IsExist(err) {
				return fmt.Errorf("failed to create user config directory")
			}

			var defaultConfig = bytes.NewBuffer([]byte(`{
	"divoltsessiontokens": [],
	"downloadcmd": {
		"outputdir": "",
		"quality": 0,
		"timeout": {
			"seconds": 0,
			"minutes": 2
		},
		"ignore": {
			"cover": false,
			"subdirs": false
		},
		"skip": {
			"unzip": false
		}
	}
}`))

			// load default config
			if err := viper.ReadConfig(defaultConfig); err != nil {
				return fmt.Errorf("failed to load default config")
			}

			// write default config
			if err := viper.WriteConfigAs(ConfigFile); err != nil {
				return fmt.Errorf("failed to write default config")
			}
		} else {
			// config was found but there was some other error
			return fmt.Errorf("unknown error with config")
		}
	}

	return nil
}

func Offload() error {
	return viper.WriteConfigAs(ConfigFile)
}
