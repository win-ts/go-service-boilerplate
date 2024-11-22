package service

import (
	"context"

	"github.com/win-ts/go-service-boilerplate/server/clean-http-polyrepo/dto"
)

func (s *service) DoKafkaProduce(ctx context.Context) error {
	_ = ctx
	return s.kafkaProducerRepository.Produce(dto.Event{
		Event: "example-event",
		Payload: dto.EventPayload{
			Data: "example data",
		},
	})
}
