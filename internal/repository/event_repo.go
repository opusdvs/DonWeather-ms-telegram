package repository

import (
	"context"
	"encoding/json"
	"log"
	"time"

	"github.com/nats-io/nats.go"
	"github.com/opusdvs/DonWeather-ms-telegram/internal/domain"
)

type EventConsumerRepository struct {
	natsConn nats.JetStreamContext
}

func NewEventConsumerRepository(natsConn *nats.Conn) *EventConsumerRepository {
	jetStreamContext, err := natsConn.JetStream()
	if err != nil {
		log.Fatalf("failed to get jetstream context: %s", err)
	}
	return &EventConsumerRepository{natsConn: jetStreamContext}
}

func (e *EventConsumerRepository) GetEventWeatherMessage(ctx context.Context, f func(domain.EventWeatherMessage) error) error {
	// Подключаемся к JetStream и создаём durable pull subscription
	sub, err := e.natsConn.PullSubscribe(
		"event.telegram.message.send",
		"weather-event-consumer",
		nats.ManualAck(),
		nats.DeliverLast(),
	)
	if err != nil {
		log.Printf("failed to pull subscribe: %s", err)
		return err
	}
	defer sub.Unsubscribe()

	for {
		// Fetch с MaxWait вместо использования контекста с таймаутом
		msgs, err := sub.Fetch(1, nats.MaxWait(5*time.Second))
		if err != nil {
			if err == nats.ErrTimeout {
				continue // просто ждем следующего сообщения
			}
			if ctx.Err() != nil {
				log.Printf("consumer context cancelled: %s", ctx.Err())
				return nil
			}
			log.Printf("failed to fetch message: %s", err)
			continue
		}

		for _, msg := range msgs {
			var event domain.EventWeatherMessage
			if err := json.Unmarshal(msg.Data, &event); err != nil {
				log.Printf("failed to unmarshal message: %s", err)
				msg.Term() // Terminate message для удаления из очереди
				continue
			}
			log.Printf("event: %+v", event)
			if err := f(event); err != nil {
				log.Printf("failed to process message: %s", err)
				msg.Nak() // повторно в очередь
				continue
			}

			// После успешной обработки — Ack, чтобы не приходило снова
			if err := msg.Ack(); err != nil {
				log.Printf("failed to ack message: %s", err)
			}
		}
	}
}
