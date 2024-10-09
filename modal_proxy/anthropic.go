package modal_proxy

import (
	"bifrost/maxim"
	"bufio"
	"bytes"
	"compress/gzip"
	"fmt"
	"github.com/andybalholm/brotli"
	"github.com/gofiber/fiber/v2"
	"io"
	"math/rand"
	"net/http"
	"strings"
)

type AnthropicModalProvider struct {
	apiUrl string
}

func NewAnthropicModalProvider(apiUrl string) *AnthropicModalProvider {
	return &AnthropicModalProvider{
		apiUrl: apiUrl,
	}
}

func (mp *AnthropicModalProvider) GetApiKey(reqHeaders map[string][]string, modal string) (string, error) {
	maximApiKey, err := GetMaximApiKey(reqHeaders)
	if err != nil {
		return "", err
	}
	account, err := maxim.GetMaximAccount(maximApiKey)
	if err != nil {
		return "", err
	}
	//FIXME: This is a temporary solution to get a random API key from the list of API keys,
	// eventually we will need to implement a more sophisticated way to select the API key which
	// looks at the response as well
	index := rand.Int() % len(account.Data.Anthropic)
	return account.Data.Anthropic[index].APIKey, nil
}

// GetCompletion Implement method.
func (mp *AnthropicModalProvider) GetCompletion(c *fiber.Ctx, apiPath string) error {
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
		return c.Status(resp.StatusCode).SendString("Error response from Anthropic API: " + resp.Status)
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
