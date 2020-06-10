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
		result = mainMenu()
	}
}

func mainMenu() string {
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
		download()
	case "upload":
		upload()
	case "edit config":
		editConfig()
	}

	return result
}

func download() error {
	fmt.Println("downloading file...")

	//TODO: put path and file into config
	resp, err := getRequest("/repos/kkotowich/cloudtorio-save/contents/README.md")

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
