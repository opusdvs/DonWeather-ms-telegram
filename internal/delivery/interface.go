package delivery

import (
	"net/http"
)

type MessageHandler interface {
	Webhook(w http.ResponseWriter, r *http.Request)
	SetWebhook(w http.ResponseWriter, r *http.Request)
}
