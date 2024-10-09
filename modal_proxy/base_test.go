package modal_proxy

import "testing"

func TestGetMaximAPIKey(t *testing.T) {
	// Test case where x-maxim-api-key exists
	t.Run("API key exists", func(t *testing.T) {
		input := map[string][]string{
			"x-maxim-api-key": {"example-api-key"},
			"another-key":     {"value1"},
		}

		expected := "example-api-key"
		result, err := GetMaximApiKey(input)
		if err != nil {
			t.Fatalf("expected no error, but got %v", err)
		}
		if result != expected {
			t.Errorf("expected %v, got %v", expected, result)
		}
	})

	// Test case where x-maxim-api-key does not exist
	t.Run("API key does not exist", func(t *testing.T) {
		input := map[string][]string{
			"some-other-key": {"value1"},
		}

		_, err := GetMaximApiKey(input)
		if err == nil {
			t.Fatal("expected an error but got none")
		}
		expectedErrMsg := "x-maxim-api-key not found"
		if err.Error() != expectedErrMsg {
			t.Errorf("expected error message %v, but got %v", expectedErrMsg, err.Error())
		}
	})

	// Test case where x-maxim-api-key exists but has no associated value
	t.Run("API key exists but no value", func(t *testing.T) {
		input := map[string][]string{
			"x-maxim-api-key": {},
			"another-key":     {"value1"},
		}

		_, err := GetMaximApiKey(input)
		if err == nil {
			t.Fatal("expected an error but got none")
		}
		expectedErrMsg := "x-maxim-api-key exists but no value associated with it"
		if err.Error() != expectedErrMsg {
			t.Errorf("expected error message %v, but got %v", expectedErrMsg, err.Error())
		}
	})
}
