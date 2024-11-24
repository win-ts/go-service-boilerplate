package service

import (
	"context"
	"net/http"
	"testing"
	"time"

	"github.com/redis/go-redis/v9"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"github.com/win-ts/go-service-boilerplate/job/clean-job-cronjob/dto"
	"github.com/win-ts/go-service-boilerplate/job/clean-job-cronjob/pkg/httpclient"
)

var (
	mockContext       = context.Background()
	mockTestDBResults = []dto.TestEntity{
		{
			ID:      "1",
			Message: "test1",
		},
		{
			ID:      "2",
			Message: "test2",
		},
	}
)

type mockExampleRepository struct {
	mock.Mock
}

func (m *mockExampleRepository) DoExample(ctx context.Context) (string, error) {
	args := m.Called(ctx)
	return args.String(0), args.Error(1)
}

type mockWiremockAPIRepository struct {
	mock.Mock
}

func (m *mockWiremockAPIRepository) GetTest(ctx context.Context, h dto.WiremockGetTestHeader) (*httpclient.Response[dto.WiremockGetTestResponse], error) {
	args := m.Called(ctx, h)
	return args.Get(0).(*httpclient.Response[dto.WiremockGetTestResponse]), args.Error(1)
}

type mockDatabaseRepository struct {
	mock.Mock
}

func (m *mockDatabaseRepository) QueryTest() (*[]dto.TestEntity, error) {
	args := m.Called()
	return args.Get(0).(*[]dto.TestEntity), args.Error(1)
}

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

func TestDoExample(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		mockExampleRepository := new(mockExampleRepository)
		mockExampleRepository.On("DoExample", mock.Anything).Return("example", nil)

		s := New(Dependencies{
			ExampleRepository: mockExampleRepository,
		})

		result, err := s.DoExample(mockContext)

		assert.NoError(t, err)
		assert.Equal(t, "example", result)
	})
}

func TestDoWiremock(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		mockWiremockAPIRepository := new(mockWiremockAPIRepository)
		mockWiremockAPIRepository.On("GetTest", mock.Anything, mock.Anything).Return(&httpclient.Response[dto.WiremockGetTestResponse]{
			HTTPStatusCode: http.StatusOK,
			Response: dto.WiremockGetTestResponse{
				Message: "test",
			},
		}, nil)

		s := New(Dependencies{
			WiremockAPIRepository: mockWiremockAPIRepository,
		})

		result, err := s.DoWiremock(mockContext)

		assert.NoError(t, err)
		assert.Equal(t, "test", result.Message)
	})
}

func TestDoDatabase(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		mockDatabaseRepository := new(mockDatabaseRepository)
		mockDatabaseRepository.On("QueryTest").Return(&mockTestDBResults, nil)

		s := New(Dependencies{
			DatabaseRepository: mockDatabaseRepository,
		})

		result, err := s.DoDBTest()

		assert.NoError(t, err)
		assert.Equal(t, mockTestDBResults[0].ID, (*result)[0].ID)
		assert.Equal(t, mockTestDBResults[0].Message, (*result)[0].Message)
		assert.Equal(t, mockTestDBResults[1].ID, (*result)[1].ID)
		assert.Equal(t, mockTestDBResults[1].Message, (*result)[1].Message)
	})
}

func TestDoSetGetCache(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		mockCacheRepository := new(mockCacheRepository)
		mockCacheRepository.On("Set", mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(&redis.StatusCmd{}, nil)
		mockCacheRepository.On("Get", mock.Anything, mock.Anything).Return(redis.NewStringResult(`{"id":"1","message":"example"}`, nil))

		s := New(Dependencies{
			CacheRepository: mockCacheRepository,
		})

		result, err := s.DoSetGetCache(mockContext)

		assert.NoError(t, err)
		assert.Equal(t, "1", result.ID)
		assert.Equal(t, "example", result.Message)
	})
}
