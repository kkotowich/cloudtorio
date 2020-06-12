package main

import (
	"io"
	"net/http"
)

// RepoFile is a file from a github repo
type RepoFile struct {
	Type     string `json:"type"`
	Encoding string `json:"encoding"`
	Size     int    `json:"size"`
	Name     string `json:"name"`
	Path     string `json:"path"`
	Content  string `json:"content"`
	SHA      string `json:"sha"`
}

// RepoDirectory is a slice of RepoFiles
type RepoDirectory struct {
	RepoFiles []RepoFile `json:""`
}

// FileCommit to upload a file to the repo
type FileCommit struct {
	Message string `json:"message"`
	// Committer Committer `json:"committer,omitempty"`
	Content string `json:"content"`
	SHA     string `json:"sha,omitempty"`
}

// Committer is the user that uploads a file to the repo
type Committer struct {
	Name  string `json:"name"`
	Email string `json:"email"`
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

	//TODO: http error code handling

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

// Put makes a PUT request
func (rh RequestHelper) Put(path string, body io.Reader) (*http.Response, error) {
	req, err := http.NewRequest("PUT", rh.BaseURL+path, body)

	if err != nil {
		return nil, err
	}

	return rh.makeRequest(req)
}
