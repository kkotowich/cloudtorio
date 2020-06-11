package main

import (
	"encoding/base64"
	"encoding/json"
	"fmt"

	"github.com/manifoldco/promptui"
)

func main() {
	var result string

	config, err := readConfig()

	if err != nil {
		fmt.Printf("error opening config... %v\n", err)
		panic("panic!")
	}

	fmt.Println(config.toString())

	for result != "exit" {
		result = mainMenu(config)
	}
}

func mainMenu(config Config) string {
	prompt := promptui.Select{
		Label: "Select",
		Items: []string{"download", "upload", "edit config", "exit"},
	}

	_, result, err := prompt.Run()

	if err != nil {
		fmt.Printf("Prompt failed %v\n", err)
		return ""
	}

	switch result {
	case "download":
		download(config)
	case "upload":
		upload()
	case "edit config":
		editConfig()
	}

	return result
}

func download(config Config) error {
	fmt.Println("downloading file...")

	resp, err := getRequest("/repos/"+config.Username+"/"+config.Repo+"/contents/"+config.SaveGameName+".zip", config.APIKey)

	if err != nil {
		return err
	}

	var repoFile RepoFile
	json.NewDecoder(resp.Body).Decode(&repoFile)

	decoded, err := base64.StdEncoding.DecodeString(repoFile.Content)

	if err != nil {
		return err
	}

	fmt.Printf("repoFile: %v\n", repoFile.Content)
	fmt.Println(string(decoded))

	return nil
}

func upload() {
	fmt.Println("welcome to upload!")
}
