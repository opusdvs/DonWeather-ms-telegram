package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	_ "github.com/jackc/pgx/v5/stdlib"

	"github.com/opusdvs/DonWeather-ms-telegram/internal/delivery"
	"github.com/opusdvs/DonWeather-ms-telegram/internal/repository"
	"github.com/opusdvs/DonWeather-ms-telegram/internal/usecase"
)

func main() {
	appCtx, appCancel := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer appCancel()

	botToken := os.Getenv("TELEGRAM_BOT_TOKEN")
	if botToken == "" {
		log.Fatal("TELEGRAM_BOT_TOKEN environment variable is required")
	}
	webhookUrl := os.Getenv("WEBHOOK_URL")
	if webhookUrl == "" {
		log.Fatal("WEBHOOK_URL environment variable is required")
	}
	apiUrl := os.Getenv("API_URL")
	if apiUrl == "" {
		log.Fatal("API_URL environment variable is required")
	}
	messageRepository := repository.NewHTTPMessageRepository(apiUrl)
	messageService := usecase.NewMessageService(messageRepository)
	messageHandler := delivery.NewTelegramMessageHandler(*messageService, appCtx)

	setWebhookRequest := delivery.NewSetWebhookRequest(botToken)
	err := setWebhookRequest.SetWebhook(webhookUrl)
	if err != nil {
		log.Fatal(err)
	}
	log.Println("Webhook set successfully")
	mainMux := http.NewServeMux()
	mainMux.HandleFunc("/api/v1/telegram/webhook", messageHandler.Webhook)

	server := &http.Server{
		Addr:         ":8080",
		Handler:      mainMux,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	go func() {
		log.Println("Server started on port 8080")
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatal(err)
		}
		log.Println("Server stopped")
	}()

	<-appCtx.Done()

	log.Println("Shutdown signal received")

	shutdownCtx, cancel := context.WithTimeout(
		context.Background(),
		10*time.Second,
	)
	defer cancel()

	if err := server.Shutdown(shutdownCtx); err != nil {
		log.Fatal(err)
	}

	log.Println("Server stopped cleanly")
}
