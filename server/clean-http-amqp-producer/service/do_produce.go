package service

import (
	"context"

	"github.com/win-ts/go-service-boilerplate/server/clean-http-amqp-producer/dto"
)

func (s *service) DoProduce(ctx context.Context, event dto.Event) error {
	return s.amqpProducerRepository.ProduceJSONMessage(ctx, event)
}
