package service

import (
	"context"
	"testing"
	"time"

	"github.com/redis/go-redis/v9"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"github.com/win-ts/go-service-boilerplate/consumer/kafka-consumer/dto"
)

var (
	mockContext = context.Background()
	mockMessage = dto.Event{
		Event: "example-event",
		Payload: dto.EventPayload{
			Data: "example",
		},
	}
	mockMessageJSON = `{"event":"example-event","payload":{"data":"example"}}`
)

type mockCacheRepository struct {
	mock.Mock
}

func (m *mockCacheRepository) Get(ctx context.Context, key string) *redis.StringCmd {
	args := m.Called(ctx, key)
	return args.Get(0).(*redis.StringCmd)
}

func (m *mockCacheRepository) Set(ctx context.Context, key string, value interface{}, expiration time.Duration) *redis.StatusCmd {
	args := m.Called(ctx, key, value, expiration)
	return args.Get(0).(*redis.StatusCmd)
}

func TestDoSetGetCache(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		mockCacheRepository := new(mockCacheRepository)
		mockCacheRepository.On("Set", mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(&redis.StatusCmd{}, nil)
		mockCacheRepository.On("Get", mock.Anything, mock.Anything).Return(redis.NewStringResult(mockMessageJSON, nil))

		s := New(Dependencies{
			CacheRepository: mockCacheRepository,
		})

		result, err := s.DoSetGetCache(mockContext, mockMessage)

		assert.NoError(t, err)
		assert.Equal(t, mockMessage, *result)
	})
}
