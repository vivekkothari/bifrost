package maxim

import (
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
