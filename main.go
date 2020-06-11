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

	rh := NewRequestHelper(config.APIKey)
	resp, err := rh.Get("/repos/" + config.Username + "/" + config.Repo + "/contents/" + config.SaveGameName + ".zip")

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

	fmt.Println("file successfully downloaded...")
	fmt.Println("writing to disk...")

	//TODO: write to disk

	fmt.Println("file successfully saved to disk...")

	return nil
}

func upload() error {
	fmt.Println("uploading file...")

	// TODO
	// "/projects/:id/repository/commits"

	fmt.Println("file successfully uploaded...")

	return nil
}
