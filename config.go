package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"strings"
)

// Config for user to connect to save game repository
type Config struct {
	Username     string `json:"username"`
	Repo         string `json:"repo"`
	RepoURL      string `json:"repoUrl"`
	APIKey       string `json:"apiKey"`
	SaveGamePath string `json:"saveGamePath"`
	SaveGameName string `json:"saveGameName"`
}

func (c Config) toString() string {
	return "  username: " + c.Username +
		"\n  repo: " + c.Repo +
		"\n  repoUrl: " + c.RepoURL +
		"\n  apiKey: " + c.APIKey +
		"\n  saveGamePath: " + c.SaveGamePath +
		"\n  saveGameName: " + c.SaveGameName
}

func writeConfig(config Config) error {
	os.Remove("./config.json")

	configFile, err := os.Create("./config.json")

	if err != nil {
		return err
	}

	defer configFile.Close()

	data := "{\n" +
		"  \"username\":\"" + config.Username + "\",\n" +
		"  \"repo\": \"" + config.Repo + "\",\n" +
		"  \"repoUrl\": \"" + config.RepoURL + "\",\n" +
		"  \"apiKey\": \"" + config.APIKey + "\",\n" +
		"  \"saveGamePath\": \"" + config.SaveGamePath + "\",\n" +
		"  \"saveGameName\": \"" + config.SaveGameName + "\"\n" +
		"}"
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
		repo         string
		username     string
		repoURL      string
		apiKey       string
		saveGamePath string
		saveGameName string
	)

	fmt.Println("")
	fmt.Println("edit config")

	fmt.Print("api key: ")
	fmt.Scanln(&apiKey)

	fmt.Print("repository url (https://github.com/kkotowich/cloudtorio-save): ")
	fmt.Scanln(&repoURL)
	if len(repoURL) == 0 {
		repoURL = "https://github.com/kkotowich/cloudtorio-save"
	}

	//TODO: escape string into config
	// fmt.Print("save game directory (%appdata\\Factorio\\saves): ")
	// fmt.Scanln(&saveGamePath)
	// if len(saveGamePath) == 0 {
	// 	saveGamePath = "%appdata\\Factorio\\saves"
	// }

	fmt.Print("save game name (cloudtorio-save): ")
	fmt.Scanln(&saveGameName)
	if len(saveGameName) == 0 {
		saveGameName = "cloudtorio-save"
	}

	username, repo = parseRepoURL(repoURL)

	config := Config{
		APIKey:       apiKey,
		Repo:         repo,
		Username:     username,
		RepoURL:      repoURL,
		SaveGamePath: saveGamePath,
		SaveGameName: saveGameName,
	}

	err := writeConfig(config)

	if err != nil {
		return err
	}

	return nil
}

func parseRepoURL(repoURL string) (username string, repo string) {
	tokens := strings.Split(repoURL, "/")
	return tokens[3], tokens[4]
}
