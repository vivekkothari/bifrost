package modal_proxy

import (
	"bufio"
	"bytes"
	"compress/gzip"
	"fmt"
	"github.com/andybalholm/brotli"
	"github.com/gofiber/fiber/v2"
	"io"
	"net/http"
	"strings"
)

type OpenAIModalProvider struct {
	apiUrl string
}

func NewOpenAIProvider(apiUrl string) *OpenAIModalProvider {
	return &OpenAIModalProvider{
		apiUrl: apiUrl,
	}
}

// GetCompletion Implement method.
func (mp *OpenAIModalProvider) GetCompletion(c *fiber.Ctx, apiPath string) error {
	if c.Method() != http.MethodPost {
		return c.Status(fiber.StatusMethodNotAllowed).SendString("Only POST method is allowed")
	}
	fmt.Printf("Received request to OpenAI API %s\n", string(c.Body()))
	req, err := http.NewRequest(http.MethodPost, mp.apiUrl+apiPath, bytes.NewBuffer(c.Body()))
	if err != nil || req == nil {
		return c.Status(fiber.StatusInternalServerError).SendString("Error creating request")
	}
	copyHeadersFromIncomingRequest(c, req)
	resp, err := client.Do(req)
	if err != nil || resp == nil {
		return c.Status(fiber.StatusInternalServerError).SendString("Error making request to OpenAI API")
	}
	copyReadersToOutgoingResponse(c, resp)
	if resp.Body == nil {
		return c.Status(fiber.StatusInternalServerError).SendString("Error: response body is nil")
	}
	if resp.StatusCode != http.StatusOK {
		return c.Status(resp.StatusCode).SendString("Error response from OpenAI API: " + resp.Status)
	}

	// Detect Content-Encoding and handle Brotli, Gzip, or plain text
	var reader io.Reader = resp.Body

	switch resp.Header.Get("Content-Encoding") {
	case "br":
		reader = brotli.NewReader(resp.Body)
	case "gzip":
		gzipReader, err := gzip.NewReader(resp.Body)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).SendString("Error reading gzip response")
		}
		reader = gzipReader
	}

	bufReader := bufio.NewReader(reader)

	contentType := resp.Header.Get("Content-Type")
	if strings.Contains(contentType, "text/event-stream") {
		streamResponse(c, resp, bufReader)
	} else {
		// Handle non-streaming content (read all at once)
		return blockingResponse(c, reader)
	}
	return nil
}
