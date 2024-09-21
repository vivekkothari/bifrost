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
	"time"
)

type OpenAIModalProvider struct {
	apiUrl string
}

func NewOpenAIProvider(apiUrl string) *OpenAIModalProvider {
	return &OpenAIModalProvider{
		apiUrl: apiUrl,
	}
}

var client = &http.Client{
	Timeout: 2 * time.Minute, // Add timeout for the HTTP client
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
		c.Status(fiber.StatusOK).Context().SetBodyStreamWriter(func(w *bufio.Writer) {
			bufWriter := bufio.NewWriter(w)
			defer closeResponse(resp)
			for {
				lineBytes, err := bufReader.ReadBytes('\n')
				if err != nil {
					if err == io.EOF {
						break
					}
					_, err := fmt.Fprintf(bufWriter, "Error reading response from OpenAI API\n")
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
	} else {
		// Handle non-streaming content (read all at once)
		bodyBytes, err := io.ReadAll(reader)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).SendString("Error reading response body")
		}
		// Remove the Content-Encoding header because the content has been decompressed
		c.Response().Header.Del("Content-Encoding")
		return c.Status(fiber.StatusOK).SendString(string(bodyBytes))
	}
	return nil
}

func closeResponse(resp *http.Response) {
	func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			fmt.Printf("Error closing response body: %v\n", err)
		}
	}(resp.Body)
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
