package repository

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"time"
)

type LinkTelegramIdToTokenRequestBody struct {
	TelegramId int64  `json:"telegram_id"`
	Token      string `json:"token"`
}
type LinkTelegramIdToTokenResponseBody struct {
	Token string `json:"token"`
}
type HTTPMessageRepository struct {
	client *http.Client
	apiUrl string
}

func NewHTTPMessageRepository(apiUrl string) *HTTPMessageRepository {
	return &HTTPMessageRepository{client: &http.Client{
		Timeout: 10 * time.Second,
	}, apiUrl: apiUrl}
}

func (r *HTTPMessageRepository) LinkTelegramIdToToken(ctx context.Context, telegramId int64, token string) (string, error) {
	log.Println("Setting telegram ID:", telegramId, "for token:", token)
	var responseBody LinkTelegramIdToTokenResponseBody
	requestBody := LinkTelegramIdToTokenRequestBody{
		TelegramId: telegramId,
		Token:      token,
	}
	requestBodyBytes, err := json.Marshal(requestBody)
	if err != nil {
		log.Println("Error marshalling request body:", err)
		return "", err
	}
	request, err := http.NewRequestWithContext(ctx, http.MethodPost, r.apiUrl, bytes.NewBuffer(requestBodyBytes))
	if err != nil {
		log.Println("Error creating request:", err)
		return "", err
	}
	request.Header.Set("Content-Type", "application/json")
	response, err := r.client.Do(request)
	if err != nil {
		log.Println("Error doing request:", err)
		return "", err
	}
	defer response.Body.Close()
	if response.StatusCode != http.StatusOK && response.StatusCode != http.StatusCreated {
		log.Println("Error status code:", response.StatusCode)
		return "", fmt.Errorf("failed to link telegram ID to token: %s", responseBody.Token)
	}
	bodyBytes, err := io.ReadAll(response.Body)
	if err != nil {
		log.Println("Error reading response body:", err)
		return "", err
	}
	err = json.Unmarshal(bodyBytes, &responseBody)
	if err != nil {
		// API может вернуть при 200/201 не JSON (например "Success") — считаем успехом, возвращаем переданный токен
		log.Println("Response is not JSON (status OK), treating as success:", string(bodyBytes))
		return token, nil
	}
	if responseBody.Token == "" {
		log.Println("Error empty token in response, using request token")
		return token, nil
	}
	return responseBody.Token, nil
}
