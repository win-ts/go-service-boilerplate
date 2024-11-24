package service

import (
	"context"
	"time"

	"github.com/win-ts/go-service-boilerplate/consumer/amqt-consumer/dto"
	"github.com/win-ts/go-service-boilerplate/consumer/amqt-consumer/pkg/utils"
)

func (s *service) DoSetGetCache(ctx context.Context) (*dto.TestModel, error) {
	key := "exampleKey"
	value := dto.TestModel{
		ID:      "1",
		Message: "example",
	}
	expiration := 5 * time.Minute

	if _, err := s.cacheRepository.Set(ctx, key, value, expiration).Result(); err != nil {
		return nil, err
	}

	val, err := s.cacheRepository.Get(ctx, key).Result()
	if err != nil {
		return nil, err
	}
	res := utils.DecodeJSONfromString[dto.TestModel](val)

	return res, nil
}
