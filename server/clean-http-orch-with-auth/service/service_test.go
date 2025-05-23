package service

import (
	"context"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"github.com/win-ts/go-service-boilerplate/server/clean-http-orch-with-auth/dto"
	"github.com/win-ts/go-service-boilerplate/server/clean-http-orch-with-auth/pkg/httpclient"
)

var (
	mockContext = context.Background()
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
