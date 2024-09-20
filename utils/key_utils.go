package utils

import (
	"fmt"
	"github.com/google/uuid"
)

// GetQueryIndex generates a unique vector ID based on the provided identifier.
func GetQueryIndex(identifier interface{}) string {
	var vectorID string

	switch v := identifier.(type) {
	case string:
		// If identifier is a string, append a new UUID to it
		vectorID = fmt.Sprintf("%s_%s", v, uuid.New().String())
	case func() string:
		// If identifier is a function, call it to get the vector ID
		vectorID = v()
	default:
		// If identifier is nil or not a string/function, generate a UUID
		vectorID = uuid.New().String()
	}

	return vectorID
}
