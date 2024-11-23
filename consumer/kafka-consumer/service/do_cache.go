package service

import (
	"context"
	"time"

	"github.com/win-ts/go-service-boilerplate/consumer/kafka-consumer/dto"
	"github.com/win-ts/go-service-boilerplate/consumer/kafka-consumer/pkg/utils"
)

func (s *service) DoSetGetCache(ctx context.Context, value dto.Event) (*dto.Event, error) {
	key := "exampleKey"
	expiration := 5 * time.Minute

	if _, err := s.cacheRepository.Set(ctx, key, value, expiration).Result(); err != nil {
		return nil, err
	}

	val, err := s.cacheRepository.Get(ctx, key).Result()
	if err != nil {
		return nil, err
	}
	res := utils.DecodeJSONfromString[dto.Event](val)

	return res, nil
}
