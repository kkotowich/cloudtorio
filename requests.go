package main

import (
	"io"
	"net/http"
)

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

// RequestHelper makes http requests to github api
type RequestHelper struct {
	BaseURL string
	APIKey  string
}

// NewRequestHelper constructs a new RequestHelper
func NewRequestHelper(apiKey string) RequestHelper {
	return RequestHelper{
		BaseURL: "https://api.github.com",
		APIKey:  apiKey,
	}
}

func (rh RequestHelper) makeRequest(req *http.Request) (*http.Response, error) {

	req.Header.Add("Accept", "application/vnd.github.v3+json")
	req.Header.Add("Authorization", "token "+rh.APIKey)

	client := &http.Client{}
	return client.Do(req)
}

// Get makes a GET request
func (rh RequestHelper) Get(path string) (*http.Response, error) {

	req, err := http.NewRequest("GET", rh.BaseURL+path, nil)

	if err != nil {
		return nil, err
	}

	return rh.makeRequest(req)
}

// Post makes a POST request
func (rh RequestHelper) Post(path string, body io.Reader) (*http.Response, error) {
	req, err := http.NewRequest("POST", rh.BaseURL+path, body)

	if err != nil {
		return nil, err
	}

	return rh.makeRequest(req)
}
