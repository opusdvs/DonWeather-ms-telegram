package domain

import "context"

type MessageRepository interface {
	LinkTelegramIdToToken(ctx context.Context, telegramID int64, token string) (string, error)
}
