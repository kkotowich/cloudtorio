package main

import "net/http"

var baseURL = "https://api.github.com"

// RepoFile is a file from a github repo
type RepoFile struct {
	Type        string        `json:"type"`
	Encoding    string        `json:"encoding"`
	Size        int           `json:"size"`
	Name        string        `json:"name"`
	Path        string        `json:"path"`
	Content     string        `json:"content"`
	SHA         string        `json:"sha"`
	URL         string        `json:"url"`
	GitURL      string        `json:"git_url"`
	HTMLURL     string        `json:"html_url"`
	DownloadURL string        `json:"download_url"`
	Links       RepoFileLinks `json:"_links"`
}

// RepoFileLinks is a collection of links for a repo file
type RepoFileLinks struct {
	Git  string `json:"git"`
	Self string `json:"self"`
	HTML string `json:"html"`
}

func makeRequest(req *http.Request) (*http.Response, error) {

	req.Header.Add("Accept", "application/vnd.github.v3+json")
	req.Header.Add("Authorization", "token ab48a64d6bc93ebb5df6bd0edb86ce3fb9fe66a8")

	client := &http.Client{}
	return client.Do(req)
}

func getRequest(path string) (*http.Response, error) {

	req, err := http.NewRequest("GET", baseURL+path, nil)

	if err != nil {
		return nil, err
	}

	return makeRequest(req)
}
