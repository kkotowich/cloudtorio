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

func downloadFile(config Config, sha string) (RepoFile, error) {
	rh := NewRequestHelper(config.APIKey)

	resp, err := rh.Get("/repos/" + config.RepoOwner + "/" + config.RepoName + "/git/blobs/" + sha)

	defer resp.Body.Close()

	if err != nil {
		fmt.Println("error downloading file from repo")
		return RepoFile{}, err
	}

	var repoFile RepoFile
	json.NewDecoder(resp.Body).Decode(&repoFile)

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

	repoFile, err := downloadFile(config, repoFileMeta.SHA)

	if err != nil {
		fmt.Println("error downloading blob")
		return err
	}

	decoded, err := base64.StdEncoding.DecodeString(repoFile.Content)

	if err != nil {
		fmt.Println("error decoding contents")
		return err
	}

	fmt.Println("file successfully downloaded...")
	fmt.Println("writing to disk...")

	file, err := os.Create(config.SaveGamePath + "\\" + config.SaveGameName + ".zip")

	if err != nil {
		fmt.Println("error creating file")
		return err
	}

	defer file.Close()

	if _, err := file.Write(decoded); err != nil {
		fmt.Println("error writing contents to file")
		return err
	}
	if err := file.Sync(); err != nil {
		fmt.Println("error saving file")
		return err
	}

	fmt.Println("file successfully saved to disk...")

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
		fmt.Println("uploading new file...")
	} else {
		fmt.Println("uploading existing file...")
	}

	fc := composeFileCommit(config, encoded, sha)

	fmt.Println("uploading file...")

	body, err := json.Marshal(fc)

	rh := NewRequestHelper(config.APIKey)
	_, err = rh.Put("/repos/"+config.RepoOwner+"/"+config.RepoName+"/contents/"+config.SaveGameName+".zip", bytes.NewBuffer(body))

	if err != nil {
		fmt.Println("error uploading file...")
		return err
	}

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
	reader := bufio.NewReader(os.Stdin)
	fmt.Print("enter commit message: ")
	message, _ := reader.ReadString('\n')

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
