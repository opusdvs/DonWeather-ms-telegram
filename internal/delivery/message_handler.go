package delivery

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"regexp"
	"strings"

	"github.com/opusdvs/DonWeather-ms-telegram/internal/domain"
	"github.com/opusdvs/DonWeather-ms-telegram/internal/usecase"
)

type TelegramMessageHandler struct {
	messageService usecase.MessageService
	appCtx         context.Context
}

func NewTelegramMessageHandler(messageService usecase.MessageService, appCtx context.Context) *TelegramMessageHandler {
	return &TelegramMessageHandler{messageService: messageService, appCtx: appCtx}
}

func (h *TelegramMessageHandler) Webhook(w http.ResponseWriter, r *http.Request) {
	var update domain.Update
	fmt.Println("Webhook started")
	if err := json.NewDecoder(r.Body).Decode(&update); err != nil {
		log.Println("Error decoding request body:", err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if update.Message == (domain.Message{}) {
		log.Println("Message is required")
		http.Error(w, "Message is required", http.StatusBadRequest)
		return
	}
	w.WriteHeader(http.StatusOK)

	if update.Message.Chat.ID == 0 {
		log.Println("Chat ID is required")
		http.Error(w, "Chat ID is required", http.StatusBadRequest)
		return
	}
	telegramId := update.Message.Chat.ID

	if update.Message.Text == "" {
		log.Println("Message text is required")
		http.Error(w, "Message text is required", http.StatusBadRequest)
		return
	}
	token, err := parseToken(update.Message.Text)
	if err != nil {
		log.Println("Error parsing token:", err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	log.Println("Token:", token)

	_, err = h.messageService.LinkTelegramIdToToken(h.appCtx, telegramId, token)
	if err != nil {
		log.Println("Error linking telegram ID to token:", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func (h *TelegramMessageHandler) SetWebhook(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Webhook set successfully"))
}

func parseToken(message string) (string, error) {
	parts := strings.Split(message, " ")
	if len(parts) != 2 {
		return "", fmt.Errorf("invalid token format")
	}

	token := parts[1]
	if err := validateToken(token); err != nil {
		return "", err
	}
	return token, nil
}

func validateToken(token string) error {
	match, err := regexp.MatchString("^[0-9a-f]{64}$", token)
	if err != nil {
		return err
	}
	if !match {
		return fmt.Errorf("invalid token format")
	}
	return nil
}
