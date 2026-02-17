package usecase

import (
	"context"

	"github.com/opusdvs/DonWeather-ms-telegram/internal/domain"
)

type MessageService struct {
	messageRepository domain.MessageRepository
}

func NewMessageService(messageRepository domain.MessageRepository) *MessageService {
	return &MessageService{messageRepository: messageRepository}
}

func (m *MessageService) LinkTelegramIdToToken(ctx context.Context, telegramID int64, token string) (string, error) {
	return m.messageRepository.LinkTelegramIdToToken(ctx, telegramID, token)
}
