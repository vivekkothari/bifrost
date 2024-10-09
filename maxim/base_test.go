package maxim

import (
	"encoding/json"
	"reflect"
	"testing"
)

func TestAccountsResponse(t *testing.T) {
	// Sample JSON data similar to the one provided
	jsonData := `{
		"data": {
			"openai": [
				{
					"apiKey" : "openai-api-key",
					"name": "OpenAI",
					"modelAvailable": [
						{
							"name" : "gpt3.5",
							"id": "gpt-3.5-turbo-16k"
						}
					]
				}
			],
			"azure": [
				{
					"baseUrl": "https://azure.example.com",
					"apiKey1": "azure-api-key-1",
					"apiKey2": "azure-api-key-2",
					"deploymentIds" : [
						{
							"id": "gpt3.5",
							"model": "gpt-3.5-turbo-16k"
						}
					]
				}
			],
			"anthropic": [
				{
					"name": "Anthropic",
					"apiKey": "anthropic-api-key"
				}
			]
		}
	}`

	// Expected output based on the given JSON data
	expected := AccountsResponse{
		Data: Accounts{
			OpenAI: []OpenAI{
				{
					APIKey: "openai-api-key",
					Name:   "OpenAI",
					ModelAvailable: []ModelAvailable{
						{
							Name: "gpt3.5",
							ID:   "gpt-3.5-turbo-16k",
						},
					},
				},
			},
			Azure: []Azure{
				{
					BaseURL: "https://azure.example.com",
					APIKey1: "azure-api-key-1",
					APIKey2: "azure-api-key-2",
					DeploymentIds: []DeploymentID{
						{
							ID:    "gpt3.5",
							Model: "gpt-3.5-turbo-16k",
						},
					},
				},
			},
			Anthropic: []Anthropic{
				{
					Name:   "Anthropic",
					APIKey: "anthropic-api-key",
				},
			},
		},
	}

	// Unmarshal JSON into the Root struct
	var result AccountsResponse
	err := json.Unmarshal([]byte(jsonData), &result)
	if err != nil {
		t.Fatalf("Failed to unmarshal JSON: %v", err)
	}

	// Compare the unmarshaled result with the expected output
	if !reflect.DeepEqual(result, expected) {
		t.Errorf("Unmarshaled result does not match expected output.\nGot: %+v\nExpected: %+v", result, expected)
	}

}
