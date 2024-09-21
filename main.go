package main

import (
	"bifrost/modal_proxy"
	"fmt"
	"github.com/gofiber/fiber/v2"
	"os"
	"os/signal"
)

const PORT = 3000

func main() {
	// Initialize a new Fiber app
	app := fiber.New()

	openAiModalProvider := modal_proxy.NewOpenAIProvider("https://api.openai.com")

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

	// Graceful shutdown
	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, os.Interrupt)
	go func() {
		<-sigs
		fmt.Println("\nShutting down gracefully...")
		os.Exit(0)
	}()

	fmt.Printf("Starting proxy server on :%d\n", PORT)
	if err := app.Listen(fmt.Sprintf(":%d", PORT)); err != nil {
		fmt.Println("Error starting server:", err)
	}
}
