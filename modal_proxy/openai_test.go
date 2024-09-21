package modal_proxy

import (
	"bytes"
	"compress/gzip"
	"github.com/andybalholm/brotli"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/stretchr/testify/assert"
)

// Helper function to create Fiber App
func setupApp(provider *OpenAIModalProvider) *fiber.App {
	app := fiber.New()
	app.Post("/completion", func(ctx *fiber.Ctx) error {
		return provider.GetCompletion(ctx, "/v1/chat/completions")
	})
	return app
}

// MockRoundTripper simulates the RoundTrip for HTTP requests.
type MockRoundTripper struct {
	StatusCode int
	Body       string
	Headers    map[string]string
}

func (mrt *MockRoundTripper) RoundTrip(_ *http.Request) (*http.Response, error) {
	// Create a mock response
	response := &http.Response{
		StatusCode: mrt.StatusCode,
		Body:       io.NopCloser(strings.NewReader(mrt.Body)),
		Header:     make(http.Header),
	}

	// Add headers to the response if any
	for k, v := range mrt.Headers {
		response.Header.Set(k, v)
	}

	return response, nil
}

func mockClient(statusCode int, body string, headers map[string]string) {
	client = &http.Client{
		Transport: &MockRoundTripper{
			StatusCode: statusCode,
			Body:       body,
			Headers:    headers,
		},
		Timeout: time.Second,
	}
}

// Test for Invalid HTTP Method
func TestInvalidMethod(t *testing.T) {
	provider := NewOpenAIProvider("https://api.openai.com")
	app := setupApp(provider)

	req := httptest.NewRequest(http.MethodGet, "/completion", nil) // Invalid method (GET)
	resp, _ := app.Test(req)

	assert.Equal(t, fiber.StatusMethodNotAllowed, resp.StatusCode)
	assert.Contains(t, resp.Header.Get("Content-Type"), "text/plain")
}

// Test for Request Creation Failure
func TestRequestCreationFailure(t *testing.T) {
	provider := NewOpenAIProvider("https://invalid-url") // Simulate a bad URL
	app := setupApp(provider)

	req := httptest.NewRequest(http.MethodPost, "/completion", strings.NewReader("test body"))
	resp, _ := app.Test(req)

	assert.Equal(t, fiber.StatusInternalServerError, resp.StatusCode)
}

// Test for HTTP Client Failure
func TestHttpClientFailure(t *testing.T) {
	provider := NewOpenAIProvider("https://api.openai.com")
	app := setupApp(provider)

	// Mock the client to simulate a timeout
	client = &http.Client{
		Timeout: 1 * time.Millisecond, // Timeout to force client failure
	}

	req := httptest.NewRequest(http.MethodPost, "/completion", strings.NewReader("test body"))
	resp, _ := app.Test(req)

	assert.Equal(t, fiber.StatusInternalServerError, resp.StatusCode)
	assert.Contains(t, resp.Header.Get("Content-Type"), "text/plain")
}

// Test for Non-200 Status from OpenAI API
func TestNonOKResponse(t *testing.T) {
	provider := NewOpenAIProvider("https://api.openai.com")
	app := setupApp(provider)

	// Mock response with 404 status code
	mockClient(http.StatusNotFound, "Not Found", nil)

	req := httptest.NewRequest(http.MethodPost, "/completion", strings.NewReader("test body"))
	resp, _ := app.Test(req)

	assert.Equal(t, http.StatusNotFound, resp.StatusCode)
	assert.Contains(t, resp.Header.Get("Content-Type"), "text/plain")
}

// Test for Brotli Encoding Response
func TestBrotliEncodingResponse(t *testing.T) {
	provider := NewOpenAIProvider("https://api.openai.com")
	app := setupApp(provider)

	// Compress the body using Brotli
	var compressedBody bytes.Buffer
	brotliWriter := brotli.NewWriter(&compressedBody)
	_, err := brotliWriter.Write([]byte("Brotli encoded data"))
	if err != nil {
		t.Fatalf("Failed to compress data with Brotli: %v", err)
	}
	brotliWriter.Close()

	headers := map[string]string{"Content-Encoding": "br"}

	// Mock response with Brotli encoding
	mockClient(http.StatusOK, compressedBody.String(), headers)

	req := httptest.NewRequest(http.MethodPost, "/completion", strings.NewReader("test body"))
	resp, _ := app.Test(req)

	assert.Equal(t, http.StatusOK, resp.StatusCode)

	// Read the response body
	body, _ := io.ReadAll(resp.Body)
	assert.Contains(t, string(body), "Brotli encoded data") // Validate the decompressed content
}

// Test for Gzip Encoding Response
func TestGzipEncodingResponse(t *testing.T) {
	provider := NewOpenAIProvider("https://api.openai.com")
	app := setupApp(provider)

	// Simulate a Gzip-encoded response body
	var buf bytes.Buffer
	gz := gzip.NewWriter(&buf)
	_, err := gz.Write([]byte("compressed gzip data"))
	assert.NoError(t, err)
	gz.Close()

	headers := map[string]string{"Content-Encoding": "gzip"}

	// Mock response with Gzip encoding
	mockClient(http.StatusOK, buf.String(), headers)

	req := httptest.NewRequest(http.MethodPost, "/completion", strings.NewReader("test body"))
	resp, _ := app.Test(req)

	assert.Equal(t, http.StatusOK, resp.StatusCode)
}

// Test for Non-Stream (Regular) Response
func TestNonStreamResponse(t *testing.T) {
	provider := NewOpenAIProvider("https://api.openai.com")
	app := setupApp(provider)

	// Mock regular response with plain text
	mockClient(http.StatusOK, "plain text response", nil)

	req := httptest.NewRequest(http.MethodPost, "/completion", strings.NewReader("test body"))
	resp, _ := app.Test(req)

	assert.Equal(t, http.StatusOK, resp.StatusCode)
	body, _ := io.ReadAll(resp.Body)
	assert.Equal(t, "plain text response", string(body))
}

// Test for Error During Event-Stream Handling
func TestErrorDuringEventStream(t *testing.T) {
	provider := NewOpenAIProvider("https://api.openai.com")
	app := setupApp(provider)

	// Simulate a streaming response that causes an error mid-stream
	stream := "data: First event\n\n" + "data: Second event\n\n" + "error: Unexpected Error\n\n"
	headers := map[string]string{"Content-Type": "text/event-stream"}

	// Mock event-stream response
	mockClient(http.StatusOK, stream, headers)

	req := httptest.NewRequest(http.MethodPost, "/completion", strings.NewReader("test body"))
	resp, _ := app.Test(req)

	assert.Equal(t, http.StatusOK, resp.StatusCode)
	// You would test for proper handling of event-stream response here
}
