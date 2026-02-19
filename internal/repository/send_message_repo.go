package repository

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"
	"time"
)

type SendMessageRequestBody struct {
	ChatID    int64  `json:"chat_id"`
	Text      string `json:"text"`
	ParseMode string `json:"parse_mode"`
}

type SendMessageRepository struct {
	client *http.Client
	apiUrl string
}

func NewSendMessageRepository(apiUrl string) *SendMessageRepository {
	finalApiUrl := fmt.Sprintf("%s/sendMessage", apiUrl)
	return &SendMessageRepository{
		client: &http.Client{Timeout: 10 * time.Second},
		apiUrl: finalApiUrl,
	}
}

// AnswerCallbackQuery отвечает на callback_query от Telegram (убирает "loading" состояние кнопки)
func (r *SendMessageRepository) AnswerCallbackQuery(ctx context.Context, callbackQueryID string, text string) error {
	if callbackQueryID == "" {
		return fmt.Errorf("callback query ID is required")
	}
	answerUrl := strings.Replace(r.apiUrl, "/sendMessage", "/answerCallbackQuery", 1)
	requestBody := map[string]any{
		"callback_query_id": callbackQueryID,
	}
	if text != "" {
		requestBody["text"] = text
	}
	requestBodyBytes, err := json.Marshal(requestBody)
	if err != nil {
		return err
	}
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, answerUrl, bytes.NewBuffer(requestBodyBytes))
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/json")
	response, err := r.client.Do(req)
	if err != nil {
		return err
	}
	defer response.Body.Close()
	if response.StatusCode != http.StatusOK {
		return fmt.Errorf("failed to answer callback query: %s", response.Status)
	}
	return nil
}

func (r *SendMessageRepository) SendMessage(ctx context.Context, telegramID int64, text string) error {
	if telegramID == 0 {
		return fmt.Errorf("telegram ID is required")
	}
	if text == "" {
		return fmt.Errorf("text is required")
	}
	requestBody := SendMessageRequestBody{
		ChatID:    telegramID,
		Text:      text,
		ParseMode: "Markdown",
	}
	log.Println("Sending message:", requestBody)
	requestBodyBytes, err := json.Marshal(requestBody)
	if err != nil {
		return err
	}
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, r.apiUrl, bytes.NewBuffer(requestBodyBytes))
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/json")
	response, err := r.client.Do(req)
	if err != nil {
		return err
	}
	defer response.Body.Close()
	if response.StatusCode != http.StatusOK {
		return fmt.Errorf("failed to send message: %s", response.Status)
	}
	return nil
}

func (r *SendMessageRepository) SendInlineKeyboardMessage(ctx context.Context, telegramID int64, text string, keyboard map[string]any) error {
	if telegramID == 0 {
		return fmt.Errorf("telegram ID is required")
	}
	if text == "" {
		return fmt.Errorf("text is required")
	}
	if keyboard == nil {
		return fmt.Errorf("keyboard is required")
	}
	requestBody := map[string]any{
		"chat_id":      telegramID,
		"text":         text,
		"reply_markup": keyboard,
	}
	log.Println("Sending inline keyboard message:", requestBody)
	requestBodyBytes, err := json.Marshal(requestBody)
	if err != nil {
		return err
	}
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, r.apiUrl, bytes.NewBuffer(requestBodyBytes))
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/json")
	response, err := r.client.Do(req)
	if err != nil {
		return err
	}
	defer response.Body.Close()
	if response.StatusCode != http.StatusOK {
		return fmt.Errorf("failed to send inline keyboard message: %s", response.Status)
	}
	return nil
}
