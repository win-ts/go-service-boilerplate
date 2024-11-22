package service

import (
	"context"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"github.com/win-ts/go-service-boilerplate/server/clean-http-polyrepo/dto"
	"github.com/win-ts/go-service-boilerplate/server/clean-http-polyrepo/pkg/httpclient"
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

type mockKafkaProducerRepository struct {
	mock.Mock
}

func (m *mockKafkaProducerRepository) Produce(message dto.Event) error {
	args := m.Called(message)
	return args.Error(0)
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

func TestDoKafkaProduce(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		mockKafkaProducerRepository := new(mockKafkaProducerRepository)
		mockKafkaProducerRepository.On("Produce", mock.Anything).Return(nil)

		s := New(Dependencies{
			KafkaProducerRepository: mockKafkaProducerRepository,
		})

		err := s.DoKafkaProduce(mockContext)

		assert.NoError(t, err)
	})
}
