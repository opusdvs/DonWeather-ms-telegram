package domain

import "time"

type EventType string

const (
	EventTempDrop     EventType = "TEMP_DROP"
	EventTempRise     EventType = "TEMP_RISE"
	EventRainStarted  EventType = "RAIN_STARTED"
	EventRainStopped  EventType = "RAIN_STOPPED"
	EventStrongWind   EventType = "STRONG_WIND"
	EventHumidityDrop EventType = "HUMIDITY_DROP"
	EventHumidityRise EventType = "HUMIDITY_RISE"
	EventPressureDrop EventType = "PRESSURE_DROP"
	EventPressureRise EventType = "PRESSURE_RISE"
)

type Update struct {
	Message       Message       `json:"message"`
	CallbackQuery CallbackQuery `json:"callback_query"`
}

type Message struct {
	Chat Chat   `json:"chat"`
	Text string `json:"text"`
	From From   `json:"from"`
}

type Chat struct {
	ID int64 `json:"id"`
}

type From struct {
	ID       int64  `json:"id"`
	Username string `json:"username"`
}

type CallbackQuery struct {
	ID      string  `json:"id"`
	Data    string  `json:"data"`
	Message Message `json:"message"`
}

type EventWeatherMessage struct {
	ID         string       `json:"id"`
	TelegramID string       `json:"telegram_id"`
	Data       EventWeather `json:"data"`
	Timestamp  time.Time    `json:"timestamp"`
	Type       EventType    `json:"type"`
	Source     string       `json:"source"`
	Version    string       `json:"version"`
}

type EventWeather struct {
	City      string    `json:"city"`
	EventType EventType `json:"event_type"`
	OldValue  float64   `json:"old_value"`
	NewValue  float64   `json:"new_value"`
	CreatedAt time.Time `json:"created_at"`
}
