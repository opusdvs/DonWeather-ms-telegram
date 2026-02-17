package delivery

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

type SetWebhookRequestBody struct {
	URL string `json:"url"`
}
type SetWebhookRequest struct {
	BotToken string `json:"bot_token"`
	Client   *http.Client
}

func NewSetWebhookRequest(botToken string) *SetWebhookRequest {
	return &SetWebhookRequest{BotToken: botToken, Client: &http.Client{}}
}

func (h *SetWebhookRequest) SetWebhook(webhookUrl string) error {
	apiUrl := fmt.Sprintf("https://api.telegram.org/bot%s/setWebhook", h.BotToken)
	body := SetWebhookRequestBody{
		URL: webhookUrl,
	}
	jsonBody, err := json.Marshal(body)
	if err != nil {
		return err
	}
	request, err := http.NewRequest(http.MethodPost, apiUrl, bytes.NewBuffer(jsonBody))
	if err != nil {
		return err
	}
	request.Header.Set("Content-Type", "application/json")
	response, err := h.Client.Do(request)
	if err != nil {
		return err
	}
	defer response.Body.Close()
	respBody, err := io.ReadAll(response.Body)
	if err != nil {
		return err
	}
	if response.StatusCode != http.StatusOK {
		return fmt.Errorf("failed to set webhook: %s", string(respBody))
	}
	return nil
}
