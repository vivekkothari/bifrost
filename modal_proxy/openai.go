package modal_proxy

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/gofiber/fiber/v2"
	"io"
	"net/http"
	"strings"
)

type Message struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

type RequestBody struct {
	Model    string    `json:"model"`
	Messages []Message `json:"messages"`
	Stream   bool      `json:"stream"`
}

type OpenAIModalProvider struct {
	apiUrl string
}

func NewOpenAIProvider(apiUrl string) *OpenAIModalProvider {
	return &OpenAIModalProvider{
		apiUrl: apiUrl,
	}
}

var client = &http.Client{}

// GetCompletion Implement method.
func (mp *OpenAIModalProvider) GetCompletion(c *fiber.Ctx) error {
	if c.Method() != http.MethodPost {
		return c.Status(fiber.StatusMethodNotAllowed).SendString("Only POST method is allowed")
	}
	var reqBody RequestBody
	if err := c.BodyParser(&reqBody); err != nil {
		return c.Status(fiber.StatusBadRequest).SendString("Invalid request body")
	}

	body, err := json.Marshal(reqBody)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).SendString("Error marshalling request body")
	}

	req, err := http.NewRequest(http.MethodPost, mp.apiUrl, bytes.NewBuffer(body))
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).SendString("Error creating request")
	}
	reqHeaders := c.GetReqHeaders()
	for key, values := range reqHeaders {
		for _, value := range values {
			req.Header.Add(key, value)
		}
	}
	resp, err := client.Do(req)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).SendString("Error making request to OpenAI API")
	}

	for key, value := range resp.Header {
		c.Set(key, value[0])
	}

	if resp.StatusCode != http.StatusOK {
		return c.Status(resp.StatusCode).SendString("Error response from OpenAI API: " + resp.Status)
	}

	reader := bufio.NewReader(resp.Body)

	c.Status(fiber.StatusOK).Context().SetBodyStreamWriter(func(w *bufio.Writer) {
		for {
			lineBytes, err := reader.ReadBytes('\n')
			if err != nil {
				if err == io.EOF {
					break
				}
				c.Status(fiber.StatusInternalServerError)
				_, err := fmt.Fprint(w, "Error reading response from OpenAI API")
				if err != nil {
					break
				}
			}
			line := string(lineBytes)
			_, err = fmt.Fprint(w, line)
			err = w.Flush()
			if strings.Contains(line, "[DONE]") {
				break
			}
			if len(lineBytes) == 0 {
				continue
			}
		}
	})
	return nil
}
