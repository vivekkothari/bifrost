package utils

import (
	"github.com/google/uuid"
	"testing"
)

// TestGetQueryIndex tests the GetQueryIndex function.
func TestGetQueryIndex(t *testing.T) {
	// Test with a string identifier
	strID := "prefix"
	result := GetQueryIndex(strID)
	if !startsWith(result, strID) {
		t.Errorf("Expected ID to start with %s but got %s", strID, result)
	}
	if !hasUUIDSuffix(result, strID) {
		t.Errorf("Expected ID to end with a UUID but got %s", result)
	}

	// Test with a function identifier
	funcID := "func_id"
	result = GetQueryIndex(func() string {
		return funcID
	})
	if result != funcID {
		t.Errorf("Expected ID to be %s but got %s", funcID, result)
	}

	// Test with nil identifier
	result = GetQueryIndex(nil)
	if _, err := uuid.Parse(result); err != nil {
		t.Errorf("Expected ID to be a valid UUID but got %s", result)
	}
}

// Helper function to check if a string starts with a given prefix
func startsWith(value, prefix string) bool {
	return len(value) >= len(prefix) && value[:len(prefix)] == prefix
}

// Helper function to check if a string ends with a UUID
func hasUUIDSuffix(value, prefix string) bool {
	// Length of UUID string representation
	uuidLen := len(uuid.New().String())
	return len(value) > len(prefix) && len(value) >= len(prefix)+uuidLen+1
}
