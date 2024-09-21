package modal_proxy

import (
	"fmt"
	"github.com/gofiber/fiber/v2"
	"io"
	"net/http"
)

// ModalProviderInterface defines the interface for calling different modals.
type ModalProviderInterface interface {
	GetCompletion(c *fiber.Ctx) error
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
