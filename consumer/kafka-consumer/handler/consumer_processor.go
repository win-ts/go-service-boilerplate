package handler

import (
	"context"
	"encoding/json"
	"log/slog"

	"github.com/IBM/sarama"
	"github.com/win-ts/go-service-boilerplate/consumer/kafka-consumer/dto"
	"github.com/win-ts/go-service-boilerplate/consumer/kafka-consumer/service"
)

const (
	eventName = "example-event"
)

// ProcessorInterface represents the interface for the processor
type ProcessorInterface interface {
	Setup(sarama.ConsumerGroupSession) error
	Cleanup(sarama.ConsumerGroupSession) error
	ConsumeClaim(sarama.ConsumerGroupSession, sarama.ConsumerGroupClaim) error
}

type processor struct {
	service service.Port
}

// ProcessorDependencies represents the dependencies for the processor
type ProcessorDependencies struct {
	Service service.Port
}

// NewProcessor creates a new processor for the consumer
func NewProcessor(d ProcessorDependencies) ProcessorInterface {
	return &processor{
		service: d.Service,
	}
}

// Setup implements sarama.ConsumerGroupHandler to setup the consumer
func (p *processor) Setup(sarama.ConsumerGroupSession) error {
	slog.Info("[handler.Setup] Kafka consumer setup complete")
	return nil
}

// Cleanup implements sarama.ConsumerGroupHandler to cleanup the consumer
func (p *processor) Cleanup(sarama.ConsumerGroupSession) error {
	return nil
}

// ConsumeClaim implements sarama.ConsumerGroupHandler to consume events
func (p *processor) ConsumeClaim(session sarama.ConsumerGroupSession, claim sarama.ConsumerGroupClaim) error {
	ctx := context.Background()

	for {
		select {
		case message, ok := <-claim.Messages():
			if !ok {
				slog.Info("[handler.ConsumeClaim] no more messages to consume")
				return nil
			}

			if message == nil {
				slog.Warn("[handler.ConsumeClaim] event is nil, skipping...")
				continue
			}

			var e dto.Event
			if err := json.Unmarshal(message.Value, &e); err != nil {
				slog.Error("error - [handler.ConsumeClaim] unable to unmarshal event, skipping...")
				continue
			}

			slog.Info("[handler.ConsumeClaim] received event",
				slog.String("event", e.Event),
				slog.Any("payload", e.Payload),
			)

			switch e.Event {
			case eventName:
				res, err := p.service.DoSetGetCache(ctx, e)
				if err != nil {
					slog.Error("error - [handler.ConsumeClaim] unable to set and get message into cache",
						slog.Any("error", err),
					)
				}
				slog.Info("[handler.ConsumeClaim] cache result",
					slog.Any("result", res),
				)
			default:
				slog.Warn("[KafkaConsumer.ConsumeClaim] unknown event, skipping...",
					slog.String("event", e.Event),
				)
			}

			session.MarkMessage(message, "")
		case <-session.Context().Done():
			return nil
		}
	}
}
