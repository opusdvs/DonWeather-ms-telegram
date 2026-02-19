package domain

import "context"

type MessageRepository interface {
	LinkTelegramIdToToken(ctx context.Context, telegramID int64, token string) (string, error)
}

type EventWeatherRepository interface {
	GetEventWeatherMessage(ctx context.Context, f func(EventWeatherMessage) error) error
}

type SendMessageRepository interface {
	SendMessage(ctx context.Context, telegramID int64, text string) error
	SendInlineKeyboardMessage(ctx context.Context, telegramID int64, text string, keyboard map[string]any) error
	AnswerCallbackQuery(ctx context.Context, callbackQueryID string, text string) error
}
