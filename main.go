package main

import (
	"fmt"
	"github.com/gofiber/fiber/v2"
	"os"
	"os/signal"
	"vector-cache/modal_proxy"
)

const PORT = 3000

func main() {
	// Initialize a new Fiber app
	app := fiber.New()

	openAiModalProvider := modal_proxy.NewOpenAIProvider("https://api.openai.com/v1/chat/completions")

	//OpenAI proxy
	app.Post("/v1/chat/completions", func(ctx *fiber.Ctx) error {
		return openAiModalProvider.GetCompletion(ctx)
	})

	// Graceful shutdown
	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, os.Interrupt)
	go func() {
		<-sigs
		fmt.Println("\nShutting down gracefully...")
		os.Exit(0)
	}()

	fmt.Println("Starting proxy server on :%d", PORT)
	if err := app.Listen(fmt.Sprintf(":%d", PORT)); err != nil {
		fmt.Println("Error starting server:", err)
	}
}
