package usecase

import (
	"context"
	"fmt"
	"strconv"

	"github.com/opusdvs/DonWeather-ms-telegram/internal/domain"
)

type MessageService struct {
	messageRepository     domain.MessageRepository
	eventRepository       domain.EventWeatherRepository
	sendMessageRepository domain.SendMessageRepository
}

func NewMessageService(
	messageRepository domain.MessageRepository,
	eventRepository domain.EventWeatherRepository,
	sendMessageRepository domain.SendMessageRepository,
) *MessageService {
	return &MessageService{
		messageRepository:     messageRepository,
		eventRepository:       eventRepository,
		sendMessageRepository: sendMessageRepository,
	}
}

func (m *MessageService) LinkTelegramIdToToken(ctx context.Context, telegramID int64, token string) (string, error) {
	return m.messageRepository.LinkTelegramIdToToken(ctx, telegramID, token)
}

// SubscriptionSuccessMessage возвращает текст сообщения в Telegram при успешной подписке.
func (m *MessageService) SubscriptionSuccessMessage() string {
	return "✅ Вы успешно подписались на уведомления о погоде.\n\n" +
		"Теперь вы будете получать сообщения об изменениях погоды."
}

func (m *MessageService) SubscriptionErrorMessage() string {
	return "❌ Ошибка при подписке на уведомления о погоде.\n\n" +
		"Напоминаем: нельзя оформить две подписки на один и тот же город.\n\n" +
		"Пожалуйста, попробуйте еще раз или обратитесь к администратору.\n\n" +
		"Для связи с администратором используйте следующие контакты:\n" +
		"Telegram: opus\\_dv\n" +
		"Email: support@buildbyte.ru"
}

func (m *MessageService) UnknownCommandMessage() string {
	return "❓ Неизвестная команда. Пожалуйста, используйте команду /help для получения помощи."
}

func (m *MessageService) ErrorParsingTokenMessage() string {
	return "❌ Ошибка при парсинге токена. Для оформления подписки на погоду посетите сайт https://donweather.com"
}

func (m *MessageService) HelpMessage() string {
	return "🆘 Помощь\n\n" +
		"/menu - меню\n" +
		"/help - получить помощь\n\n" +
		"Для связи с администратором используйте следующие контакты:\n" +
		"Telegram: @opus\\_dv\n" +
		"Email: support@buildbyte.ru"
}

func (m *MessageService) MenuMessage() string {
	return "🏠 Меню\n\n" +
		"Выберите действие:"
}

func (m *MessageService) GetEventWeatherMessage(ctx context.Context) error {
	return m.eventRepository.GetEventWeatherMessage(ctx, func(message domain.EventWeatherMessage) error {
		telegramID, err := strconv.ParseInt(message.TelegramID, 10, 64)
		if err != nil {
			return err
		}

		text := BuildMessageText(message)
		return m.SendMessage(ctx, telegramID, text)
	})
}

func (m *MessageService) SendMessage(ctx context.Context, telegramID int64, text string) error {
	return m.sendMessageRepository.SendMessage(ctx, telegramID, text)
}

// BuildInlineKeyboard возвращает раскладку inline-кнопок для сообщения о погоде.
func (m *MessageService) BuildInlineKeyboard() map[string]any {
	keyboard := map[string]any{
		"inline_keyboard": [][]map[string]string{
			{
				{"text": "📋 Мои подписки", "callback_data": "my_subscriptions"},
			},
			{
				{"text": "❌ Отменить все подписки", "callback_data": "cancel_all_subscriptions"},
			},
			{
				{"text": "📋 Подписаться на погоду", "callback_data": "subscribe_to_weather"},
			},
			{
				{"text": "🆘 Помощь", "callback_data": "help"},
			},
		},
	}
	return keyboard
}

func (m *MessageService) SendInlineKeyboardMessage(ctx context.Context, telegramID int64, text string) error {
	keyboard := m.BuildInlineKeyboard()
	return m.sendMessageRepository.SendInlineKeyboardMessage(ctx, telegramID, text, keyboard)
}

func (m *MessageService) HandleCallbackQuery(ctx context.Context, callbackQueryID string, callbackData string, chatID int64) error {
	if err := m.sendMessageRepository.AnswerCallbackQuery(ctx, callbackQueryID, ""); err != nil {
		return fmt.Errorf("failed to answer callback query: %w", err)
	}
	switch callbackData {
	case "my_subscriptions":
		text := "📋 Ваши подписки:\n\n" +
			"Список подписок будет здесь..."
		return m.SendMessage(ctx, chatID, text)
	case "cancel_all_subscriptions":
		text := "❌ Все подписки отменены."
		return m.SendMessage(ctx, chatID, text)
	case "subscribe_to_weather":
		text := "📋 Для подписки на погоду посетите сайт https://donweather.com"
		return m.SendMessage(ctx, chatID, text)
	case "help":
		return m.SendMessage(ctx, chatID, m.HelpMessage())
	default:
		text := "❓ Неизвестная команда."
		return m.SendMessage(ctx, chatID, text)
	}
}

func BuildMessageText(message domain.EventWeatherMessage) string {
	header := "🌤️ *Погода*\n\n📍 *%s*\n\n──────────────────\n\n"
	text := fmt.Sprintf(header, message.Data.City)
	switch message.Data.EventType {
	case domain.EventTempDrop:
		text += fmt.Sprintf("❄️ Температура упала на %.1f°C", message.Data.NewValue)
	case domain.EventTempRise:
		text += fmt.Sprintf("🌡️ Температура поднялась на %.1f°C", message.Data.NewValue)
	case domain.EventRainStarted:
		text += "🌧️ Начался дождь"
	case domain.EventRainStopped:
		text += "🌤️ Дождь закончился"
	case domain.EventStrongWind:
		text += "💨 Сильный ветер"
	case domain.EventHumidityDrop:
		text += fmt.Sprintf("💧 Влажность упала на %.1f%%", message.Data.NewValue)
	case domain.EventHumidityRise:
		text += fmt.Sprintf("💧 Влажность поднялась на %.1f%%", message.Data.NewValue)
	case domain.EventPressureDrop:
		text += fmt.Sprintf("📉 Давление упало на %.1f мм рт.ст.", message.Data.NewValue)
	case domain.EventPressureRise:
		text += fmt.Sprintf("📈 Давление поднялось на %.1f мм рт.ст.", message.Data.NewValue)
	default:
		text += "❓ Неизвестное"
	}
	return text
}
