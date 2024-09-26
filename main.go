package main

import (
	"bifrost/modal_proxy"
	"fmt"
	"github.com/gofiber/fiber/v2"
	"os"
	"os/signal"
	"time"
)

const PORT = 3000

func main() {
	// Initialize a new Fiber app
	app := fiber.New(
		fiber.Config{
			Prefork:           false,           // Disable prefork mode (uses multiple Go processes)
			IdleTimeout:       5 * time.Second, // Timeout for idle connections
			ReduceMemoryUsage: true,            // Reduces memory usage by freeing up resources more aggressively
		})

	openAiModalProvider := modal_proxy.NewOpenAIProvider("https://api.openai.com")
	anthropicAiModalProvider := modal_proxy.NewAnthropicModalProvider("https://api.anthropic.com")

	//OpenAI proxy
	app.Post("/v1/chat/completions", func(ctx *fiber.Ctx) error {
		return openAiModalProvider.GetCompletion(ctx, "/v1/chat/completions")
	})
	//Python client adds the v1 prefix to the endpoint, thus need to not add it here.
	app.Post("/chat/completions", func(ctx *fiber.Ctx) error {
		return openAiModalProvider.GetCompletion(ctx, "/v1/chat/completions")
	})
	//llamaindex uses completions API
	app.Post("/completions", func(ctx *fiber.Ctx) error {
		return openAiModalProvider.GetCompletion(ctx, "/v1/completions")
	})
	app.Post("/v1/messages", func(ctx *fiber.Ctx) error {
		return anthropicAiModalProvider.GetCompletion(ctx, "/v1/messages")
	})

	// Setup graceful shutdown
	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, os.Interrupt, os.Kill)

	go func() {
		<-sigs
		fmt.Println("\nShutting down gracefully...")
		// Gracefully shutdown the server
		if err := app.Shutdown(); err != nil {
			fmt.Println("Error shutting down the server:", err)
		}
	}()

	fmt.Printf("Starting proxy server on :%d\n", PORT)
	if err := app.Listen(fmt.Sprintf(":%d", PORT)); err != nil {
		fmt.Println("Error starting server:", err)
	}
}
