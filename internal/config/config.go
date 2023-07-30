package config

import (
	"bytes"
	"fmt"
	"os"
	"path/filepath"

	"github.com/spf13/viper"
)

func Load(defaultHandling bool, customPath string) error {
	var configDir string
	var configName string
	var configFile string

	if defaultHandling {
		var err error
		configDir, err = os.UserConfigDir()
		if err != nil {
			return fmt.Errorf("failed to find user config directory")
		}

		configName = "config.json"
		configDir += string(os.PathSeparator) + "SlavartDL"
		configFile = configDir + string(os.PathSeparator) + configName
	} else {
		file, err := os.Stat(customPath)
		if err != nil {
			return fmt.Errorf("failed to find custom config from 'configPath'")
		}

		if file.IsDir() {
			configName = "config.json"
			configDir = customPath
			configFile = configDir + string(os.PathSeparator) + configName
		} else {
			configName = file.Name()
			configDir = filepath.Dir(customPath)
			configFile = customPath
		}
	}

	viper.SetConfigName(configName)
	viper.SetConfigType("json")
	viper.AddConfigPath(configDir)

	if err := viper.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			if err := os.Mkdir(configDir, os.ModePerm); err != nil && !os.IsExist(err) {
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

			if err := viper.ReadConfig(defaultConfig); err != nil {
				return fmt.Errorf("failed to load default config")
			}

			if err := viper.WriteConfigAs(configFile); err != nil {
				return fmt.Errorf("failed to write default config")
			}
		} else {
			// config was found but there was some other error
			return fmt.Errorf("unknown error with config")
		}
	}

	return nil
}
