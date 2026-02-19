package delivery

import (
	"context"

	"github.com/opusdvs/DonWeather-ms-telegram/internal/usecase"
)

type EventDelivery struct {
	eventService usecase.MessageService
}

func NewEventDelivery(eventService usecase.MessageService) *EventDelivery {
	return &EventDelivery{eventService: eventService}
}

func (e *EventDelivery) StartEventDelivery(ctx context.Context) error {
	return e.eventService.GetEventWeatherMessage(ctx)
}
