package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
)

// Config for user to connect to save game repository
type Config struct {
	APIKey       string `json:"apiKey"`
	RepoURL      string `json:"repoUrl"`
	SaveGamePath string `json:"saveGamePath"`
	SaveGameName string `json:"saveGameName"`
}

func (c Config) toString() string {
	return "  apiKey: " + c.APIKey + "\n  repoUrl: " + c.RepoURL
}

func writeConfig(config Config) error {
	os.Remove("./config.json")

	configFile, err := os.Create("./config.json")

	if err != nil {
		return err
	}

	defer configFile.Close()

	data := "{\"apiKey\":\"" + config.APIKey + "\", \"repoUrl\": \"" + config.RepoURL + "\"}"
	_, err = configFile.WriteString(data)

	if err != nil {
		return err
	}

	return nil
}

func readConfig() (Config, error) {
	fmt.Println("loading config...")

	_, err := os.Stat("./config.json")

	if os.IsNotExist(err) {
		fmt.Println("config file not found...")
		err := editConfig()

		if err != nil {
			return Config{}, err
		}
	}

	configFile, err := os.Open("./config.json")

	if err != nil {
		return Config{}, err
	}

	defer configFile.Close()

	byteValue, err := ioutil.ReadAll(configFile)

	var config Config
	json.Unmarshal(byteValue, &config)

	fmt.Println("config loaded...")

	return config, nil
}

func editConfig() error {
	var (
		repoURL string
		apiKey  string
	)

	fmt.Println("")
	fmt.Println("edit config")

	fmt.Print("repository url: ")
	fmt.Scanln(&repoURL)

	fmt.Print("api key: ")
	fmt.Scanln(&apiKey)

	config := Config{APIKey: apiKey, RepoURL: repoURL}

	err := writeConfig(config)

	if err != nil {
		return err
	}

	return nil
}
