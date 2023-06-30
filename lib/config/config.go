package config

import (
	"encoding/json"
	"fmt"
	"os"
)

type Config struct {
	DivoltSessionTokens []string `json:"DivoltSessionTokens"`
}

const (
	ConfigDirectory = "/slavart"
	ConfigFilePath  = "/config.json"
)

var Public Config

func CreateConfigIfNotExist() error {
	userConfigDirectory, err := os.UserConfigDir()
	if err != nil {
		return err
	}

	configDirectory := userConfigDirectory + ConfigDirectory
	configFilePath := configDirectory + ConfigFilePath

	_, err = os.Stat(configDirectory)
	if os.IsNotExist(err) {
		os.Mkdir(configDirectory, 0777)
	} else if err != nil {
		return err
	}

	_, err = os.Stat(configFilePath)
	if os.IsNotExist(err) {
		fmt.Println(err)

		configFile, err := os.Create(configFilePath)
		if err != nil {
			return err
		}

		if err := json.NewEncoder(configFile).Encode(Config{}); err != nil {
			return err
		}

		configFile.Close()
	} else if err != nil {
		return err
	}

	return nil
}

func LoadConfig() error {
	userConfigDirectory, err := os.UserConfigDir()
	if err != nil {
		return err
	}

	if err := CreateConfigIfNotExist(); err != nil {
		return err
	}

	configFile, err := os.Open(userConfigDirectory + ConfigDirectory + ConfigFilePath)
	if err != nil {
		return err
	}
	defer configFile.Close()

	if err := json.NewDecoder(configFile).Decode(&Public); err != nil {
		return err
	}

	return nil
}

func WriteConfig() error {
	userConfigDirectory, err := os.UserConfigDir()
	if err != nil {
		return err
	}

	configFile, err := os.OpenFile(userConfigDirectory+ConfigDirectory+ConfigFilePath, os.O_RDWR, os.ModePerm)
	if err != nil {
		return err
	}
	defer configFile.Close()

	return json.NewEncoder(configFile).Encode(Public)
}
