package service

import (
	"context"
	"errors"
	"testing"
	"time"
	"usdt-rate-service/internal/models"
	"usdt-rate-service/mocks"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"go.uber.org/zap"
)
func TestRatesService_GetRates(t *testing.T) {
	logger := zap.NewNop()
	ctx := context.Background()
	market := "usdtrub"

	validDepth := &models.Depth{
		Asks: []models.Order{
			{Price: 50000.0, Amount: 1.0},
		},
		Bids: []models.Order{
			{Price: 49900.0, Amount: 1.5},
		},
		Timestamp: time.Now().Unix(),
	}

	invalidDepth := &models.Depth{
		Asks: []models.Order{},
		Bids: []models.Order{},
	}

	tests := []struct {
		name            string
		mockDepth       *models.Depth
		mockDepthErr    error
		mockSaveErr     error
		expectedErr     bool
		expectedAskPrice float64
		expectedBidPrice float64
		shouldCallSave  bool
	}{
		{
			name:            "successful rate retrieval",
			mockDepth:       validDepth,
			expectedAskPrice: 50000.0,
			expectedBidPrice: 49900.0,
			shouldCallSave:  true,
		},
		{
			name:            "depth provider error",
			mockDepth:       nil,
			mockDepthErr:    errors.New("provider error"),
			mockSaveErr:     nil,
			expectedErr:     true,
			shouldCallSave:  false,
		},
		{
			name:            "invalid depth data",
			mockDepth:       invalidDepth,
			mockDepthErr:    nil,
			mockSaveErr:     nil,
			expectedErr:     true,
			shouldCallSave:  false,
		},
		{
			name:            "repository save error",
			mockDepth:       validDepth,
			mockDepthErr:    nil,
			mockSaveErr:     errors.New("save error"),
			expectedErr:     false,
			expectedAskPrice: 50000.0,
			expectedBidPrice: 49900.0,
			shouldCallSave:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockProvider := &mocks.MockDepthProvider{}
			mockRepo := &mocks.MockRatesRepository{}

			mockProvider.On("GetDepth", ctx, market).Return(tt.mockDepth, tt.mockDepthErr)
			
			if tt.shouldCallSave {
				mockRepo.On("SaveRate", ctx, mock.MatchedBy(func(rate *models.Rate) bool {
					return rate != nil && rate.AskPrice > 0 && rate.BidPrice > 0
				})).Return(tt.mockSaveErr)
			}

			service := NewRatesService(logger, mockProvider, mockRepo)
			rate, err := service.GetRates(ctx, market)

			if tt.expectedErr {
				assert.Error(t, err)
				assert.Nil(t, rate)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, rate)
				assert.Equal(t, tt.expectedAskPrice, rate.AskPrice)
				assert.Equal(t, tt.expectedBidPrice, rate.BidPrice)
			}

			mockProvider.AssertExpectations(t)
			if tt.shouldCallSave {
				mockRepo.AssertExpectations(t)
			} else {
				mockRepo.AssertNotCalled(t, "SaveRate")
			}
		})
	}
}
