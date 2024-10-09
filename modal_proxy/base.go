package modal_proxy

import (
	"bufio"
	"errors"
	"fmt"
	"github.com/gofiber/fiber/v2"
	"io"
	"net/http"
	"strings"
	"time"
)

var client = &http.Client{
	Timeout: 2 * time.Minute, // Add timeout for the HTTP client
	Transport: &http.Transport{
		MaxIdleConns:        100,
		IdleConnTimeout:     90 * time.Second,
		MaxConnsPerHost:     100,
		MaxIdleConnsPerHost: 10,
	},
}

// ModalProviderInterface defines the interface for calling different modals.
type ModalProviderInterface interface {
	GetCompletion(c *fiber.Ctx) error

	// GetApiKey returns the API key for the modal provider and the selected modal.
	GetApiKey(reqHeaders map[string][]string, modal string) (string, error)
}

func closeResponse(resp *http.Response) {
	func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			fmt.Printf("Error closing response body: %v\n", err)
		}
	}(resp.Body)
}

func GetMaximApiKey(reqHeaders map[string][]string) (string, error) {
	const apiKey = "x-maxim-api-key"
	// Check if the key exists in the map
	if values, ok := reqHeaders[apiKey]; ok {
		// Check if there is at least one value associated with the key
		if len(values) > 0 {
			return values[0], nil
		}
		return "", errors.New("x-maxim-api-key exists but no value associated with it")
	}
	return "", errors.New("x-maxim-api-key not found")
}

func copyReadersToOutgoingResponse(c *fiber.Ctx, resp *http.Response) {
	for key, value := range resp.Header {
		for val := range value {
			c.Response().Header.Add(key, value[val])
		}
	}
}

func copyHeadersFromIncomingRequest(c *fiber.Ctx, req *http.Request) {
	reqHeaders := c.GetReqHeaders()
	for key, values := range reqHeaders {
		for _, value := range values {
			req.Header.Add(key, value)
		}
	}
}

func streamResponse(c *fiber.Ctx, resp *http.Response, bufReader *bufio.Reader) {
	c.Status(fiber.StatusOK).Context().SetBodyStreamWriter(func(w *bufio.Writer) {
		bufWriter := bufio.NewWriter(w)
		defer closeResponse(resp)
		for {
			lineBytes, err := bufReader.ReadBytes('\n')
			if err != nil {
				if err == io.EOF {
					break
				}
				_, err := fmt.Fprintf(bufWriter, "Error reading response from ModalProvider API\n")
				if err != nil {
					return
				}
				err = bufWriter.Flush()
				if err != nil {
					return
				}
				break
			}
			if len(lineBytes) == 0 {
				continue
			}

			line := string(lineBytes)
			_, err = bufWriter.WriteString(line)
			if err != nil {
				fmt.Printf("Error writing response: %v\n", err)
				break
			}

			err = bufWriter.Flush()
			if err != nil {
				fmt.Printf("Error flushing buffer: %v\n", err)
				break
			}
			if strings.Contains(line, "[DONE]") {
				break
			}
		}
	})
}

func blockingResponse(c *fiber.Ctx, reader io.Reader) error {
	bodyBytes, err := io.ReadAll(reader)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).SendString("Error reading response body")
	}
	// Remove the Content-Encoding header because the content has been decompressed
	c.Response().Header.Del("Content-Encoding")
	return c.Status(fiber.StatusOK).SendString(string(bodyBytes))
}
