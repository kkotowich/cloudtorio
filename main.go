package main

import (
	"bufio"
	"bytes"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
)

func main() {
	var result string

	config, err := readConfig()

	if err != nil {
		fmt.Printf("error opening config... %v\n", err)
		panic("panic!")
	}

	fmt.Println("")
	fmt.Println(config.toString())

	for result != "4" {
		result = mainMenu(config)
	}
}

func mainMenu(config Config) string {
	var (
		result string
		err    error
	)
	fmt.Println("")
	fmt.Println("Main Menu")
	fmt.Println("1. download")
	fmt.Println("2. upload")
	fmt.Println("3. edit config")
	fmt.Println("4. exit")
	fmt.Println("")
	fmt.Print("Select 1, 2, 3, 4: ")
	fmt.Scanln(&result)

	switch result {
	case "1":
		err = downloadToDisk(config)
	case "2":
		err = upload(config)
	case "3":
		err = editConfig()
	}

	if err != nil {
		fmt.Printf("error: %v\n", err)
	}

	return result
}

func downloadFile(config Config) (RepoFile, error) {
	rh := NewRequestHelper(config.APIKey)

	// https://api.github.com/repos/YOUR_ACCOUNT/YOUR_REPO/git/blobs/THE_SHA
	resp, err := rh.Get("/repos/" + config.RepoOwner + "/" + config.RepoName + "/contents")

	fmt.Printf("resp: %v", resp)

	if err != nil {
		fmt.Println("error downloading file from repo")
		return RepoFile{}, err
	}

	var repoFile RepoFile
	json.NewDecoder(resp.Body).Decode(&repoFile)

	resp.Body.Close()

	return repoFile, nil
}

// save game data is larger than the 1mb download restriction
// must download file as blob, but first need to scan repo to get the file sha
func downloadRepoContents(config Config) (RepoDirectory, error) {
	rh := NewRequestHelper(config.APIKey)

	path := "/repos/" + config.RepoOwner + "/" + config.RepoName + "/contents"
	resp, err := rh.Get(path)
	defer resp.Body.Close()

	if err != nil {
		fmt.Println("error downloading file from repo")
		return RepoDirectory{}, err
	}

	var repoDir RepoDirectory
	json.NewDecoder(resp.Body).Decode(&repoDir.RepoFiles)

	return repoDir, nil
}

func downloadToDisk(config Config) error {
	fmt.Println("downloading file...")

	repoDir, err := downloadRepoContents(config)

	if err != nil {
		fmt.Println("error downloading file")
		return err
	}

	repoFileMeta := findSaveGameRepoFile(repoDir, config.SaveGameName)

	fmt.Printf("repoFile: %v", repoFileMeta)

	// decoded, err := base64.StdEncoding.DecodeString(repoFile.Content)

	// if err != nil {
	// 	fmt.Println("error decoding contents")
	// 	return err
	// }

	// fmt.Printf("repoFile: %v\n", repoFile.Content)
	// fmt.Println(string(decoded))

	// fmt.Println("file successfully downloaded...")
	// fmt.Println("writing to disk...")

	// //TODO: write to disk

	// fmt.Println("file successfully saved to disk...")

	return nil
}

func upload(config Config) error {
	fmt.Println("encoding file...")

	encoded, err := encodeFileBase64(config.SaveGamePath + "\\" + config.SaveGameName + ".zip")

	if err != nil {
		fmt.Println("error encoding file")
		return err
	}

	fmt.Println("file encoded...")

	fmt.Println("fetching sha...")
	sha, err := getLatestSHA(config)

	if err != nil {
		fmt.Println("error getting latest sha")
		return err
	}

	if len(sha) == 0 {
		fmt.Println("sha: not found")
	} else {
		fmt.Printf("sha: %s\n", sha)
	}

	fc := composeFileCommit(config, encoded, sha)

	fmt.Println("uploading file...")

	body, err := json.Marshal(fc)

	rh := NewRequestHelper(config.APIKey)
	resp, err := rh.Put("/repos/"+config.RepoOwner+"/"+config.RepoName+"/contents/"+config.SaveGameName+".zip", bytes.NewBuffer(body))

	if err != nil {
		fmt.Println("error uploading file...")
		return err
	}

	fmt.Println("/repos/" + config.RepoOwner + "/" + config.RepoName + "/contents/" + config.SaveGameName + ".zip")
	fmt.Printf("response: %v\n", resp)

	fmt.Println("file successfully uploaded...")

	return nil
}

func getLatestSHA(config Config) (string, error) {
	repoDir, err := downloadRepoContents(config)

	if err != nil {
		fmt.Println("error downloading file")
		return "", err
	}

	repoFileMeta := findSaveGameRepoFile(repoDir, config.SaveGameName)

	return repoFileMeta.SHA, nil
}

func encodeFileBase64(path string) (string, error) {
	file, err := os.Open(path)

	if err != nil {
		fmt.Println("error opening file")
		return "", err
	}

	reader := bufio.NewReader(file)
	content, err := ioutil.ReadAll(reader)

	if err != nil {
		fmt.Println("error reading file")
		return "", err
	}

	encoded := base64.StdEncoding.EncodeToString(content)

	return encoded, nil
}

func composeFileCommit(config Config, encoded string, sha string) FileCommit {
	var message string
	fmt.Print("enter commit message: ")
	fmt.Scanln(&message)

	return FileCommit{
		Message: message,
		Content: encoded,
		SHA:     sha,
	}
}

func findSaveGameRepoFile(rd RepoDirectory, saveGameName string) RepoFile {
	for _, rf := range rd.RepoFiles {
		if rf.Name == saveGameName+".zip" {
			return rf
		}
	}

	return RepoFile{}
}
