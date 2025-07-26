package service_test

import (
	"context"
	"errors"
	"strconv"
	"testing"
	"time"
	"usdt-rate-service/internal/models"
	"usdt-rate-service/internal/service"
	"usdt-rate-service/mocks"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap"
)

func setupTestService(
	t *testing.T,
	depth *models.Depth,
	depthErr error,
	saveErr error,
	expectSave bool,
) (*service.RatesService, *mocks.MockDepthProvider, *mocks.MockRatesRepository) {
	t.Helper()
	ctx := context.Background()
	logger := zap.NewNop()

	mockProvider := &mocks.MockDepthProvider{}
	mockRepo := &mocks.MockRatesRepository{}

	mockProvider.On("GetDepth", ctx, "usdtrub").Return(depth, depthErr)

	if expectSave {
		mockRepo.On("SaveRate", ctx, mock.MatchedBy(matchSavedRate(depth))).Return(saveErr)
	}

	svc := service.NewRatesService(logger, mockProvider, mockRepo)
	return svc, mockProvider, mockRepo
}

func matchSavedRate(depth *models.Depth) func(*models.Rate) bool {
	return func(rate *models.Rate) bool {
		if rate == nil {
			return false
		}
		ask, _ := strconv.ParseFloat(depth.Asks[0].Price, 64)
		bid, _ := strconv.ParseFloat(depth.Bids[0].Price, 64)
		return ask > 0 && bid > 0
	}
}

func assertValidRate(t *testing.T, rate *models.Rate, ask, bid string) {
	t.Helper()
	assert.NotNil(t, rate)
	assert.Equal(t, ask, rate.AskPrice)
	assert.Equal(t, bid, rate.BidPrice)
}

func TestRatesService_GetRates(t *testing.T) {
	ctx := context.Background()

	var (
		validDepth = &models.Depth{
			Asks:      []models.Order{{Price: "50000.0", Amount: "1.0"}},
			Bids:      []models.Order{{Price: "49900.0", Amount: "1.5"}},
			Timestamp: time.Now().Unix(),
		}

		invalidDepth = &models.Depth{
			Asks: []models.Order{},
			Bids: []models.Order{},
		}
	)

	t.Run("successful rate retrieval", func(t *testing.T) {
		svc, provider, repo := setupTestService(t, validDepth, nil, nil, true)
		rate, err := svc.GetRates(ctx, "usdtrub")
		require.NoError(t, err)
		assertValidRate(t, rate, "50000.0", "49900.0")
		provider.AssertExpectations(t)
		repo.AssertExpectations(t)
	})

	t.Run("depth provider error", func(t *testing.T) {
		svc, provider, repo := setupTestService(t, nil, errors.New("provider error"), nil, false)
		rate, err := svc.GetRates(ctx, "usdtrub")
		require.Error(t, err)
		assert.Nil(t, rate)
		provider.AssertExpectations(t)
		repo.AssertNotCalled(t, "SaveRate")
	})

	t.Run("invalid depth data", func(t *testing.T) {
		svc, provider, repo := setupTestService(t, invalidDepth, nil, nil, false)
		rate, err := svc.GetRates(ctx, "usdtrub")
		require.Error(t, err)
		assert.Nil(t, rate)
		provider.AssertExpectations(t)
		repo.AssertNotCalled(t, "SaveRate")
	})

	t.Run("repository save error", func(t *testing.T) {
		svc, provider, repo := setupTestService(t, validDepth, nil, errors.New("save error"), true)
		_, err := svc.GetRates(ctx, "usdtrub")
		require.Error(t, err)
		provider.AssertExpectations(t)
		repo.AssertExpectations(t)
	})
}
