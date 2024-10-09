package maxim

import (
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"time"
)

/**
This file has methods to interact with the Maxim API.
*/

// Struct to represent the 'modelAvailable' field in 'openai'
type ModelAvailable struct {
	Name string `json:"name"`
	ID   string `json:"id"`
}

// Struct to represent an item in the 'openai' field
type OpenAI struct {
	APIKey         string           `json:"apiKey"`
	Name           string           `json:"name"`
	ModelAvailable []ModelAvailable `json:"modelAvailable"`
}

// Struct to represent the 'deploymentIds' field in 'azure'
type DeploymentID struct {
	ID    string `json:"id"`
	Model string `json:"model"`
}

// Struct to represent an item in the 'azure' field
type Azure struct {
	BaseURL       string         `json:"baseUrl"`
	APIKey1       string         `json:"apiKey1"`
	APIKey2       string         `json:"apiKey2"`
	DeploymentIds []DeploymentID `json:"deploymentIds"`
}

// Struct to represent an item in the 'anthropic' field
type Anthropic struct {
	Name   string `json:"name"`
	APIKey string `json:"apiKey"`
}

// AccountsResponse struct to represent the entire JSON
type Accounts struct {
	OpenAI    []OpenAI    `json:"openai"`
	Azure     []Azure     `json:"azure"`
	Anthropic []Anthropic `json:"anthropic"`
}

// Struct to represent the accounts response
type AccountsResponse struct {
	Data Accounts `json:"data"`
}

var client = &http.Client{
	Timeout: 2 * time.Minute, // Add timeout for the HTTP client
	Transport: &http.Transport{
		MaxIdleConns:        100,
		IdleConnTimeout:     90 * time.Second,
		MaxConnsPerHost:     100,
		MaxIdleConnsPerHost: 10,
	},
}

// GetMaximAccount gets the accounts from the Maxim API
func GetMaximAccount(key string) (AccountsResponse, error) {
	// Call the Maxim API to get the accounts
	var result AccountsResponse
	//FIXME: Fix the API path
	req, err := http.NewRequest(http.MethodGet, "https://maxim.example.com/api/bifrost/v1/accounts", nil)
	if err != nil {
		return AccountsResponse{}, err
	}
	req.Header.Set("x-maxim-api-key", key)
	resp, err := client.Do(req)
	if err != nil {
		return AccountsResponse{}, err
	}
	if resp.Body == nil {
		return AccountsResponse{}, errors.New("response body is nil")
	}
	bodyBytes, err := io.ReadAll(resp.Body)
	err = json.Unmarshal(bodyBytes, &result)
	if err != nil {
		return AccountsResponse{}, err
	}
	return result, nil
}
