package modal_proxy

import "github.com/gofiber/fiber/v2"

// ModalProviderInterface defines the interface for calling different modals.
type ModalProviderInterface interface {
	GetCompletion(c *fiber.Ctx) error
}
