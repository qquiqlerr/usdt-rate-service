package service

import (
	"context"
	"usdt-rate-service/internal/models"

	"go.uber.org/zap"
)

type DepthProvider interface {
	GetDepth(ctx context.Context, market string) (*models.Depth, error)
}

type RatesRepository interface {
	SaveRate(ctx context.Context, rate *models.Rate) error
}

type RatesService struct {
	logger          *zap.Logger
	depthProvider   DepthProvider
	ratesRepository RatesRepository
}

func NewRatesService(logger *zap.Logger, depthProvider DepthProvider, ratesRepository RatesRepository) *RatesService {
	return &RatesService{
		logger:          logger,
		depthProvider:   depthProvider,
		ratesRepository: ratesRepository,
	}
}

func (s *RatesService) GetRates(ctx context.Context, market string) (*models.Rate, error) {
	logger := s.logger.With(
		zap.String("service", "RatesService"),
		zap.String("method", "GetRates"),
		zap.String("market", market),
	)
	logger.Debug("Getting rates")
	// 1. Get depth data from the provider
	depth, err := s.depthProvider.GetDepth(ctx, market)
	if err != nil {
		logger.Error("Failed to get depth", zap.Error(err))
		return nil, err
	}

	// Validate the depth data
	if err := depth.Validate(); err != nil {
		logger.Error("Invalid depth data", zap.Error(err))
		return nil, err
	}

	// 2. Create a Rate model from the depth data
	rate := &models.Rate{
		Market:    market,
		AskPrice:  depth.Asks[0].Price,
		BidPrice:  depth.Bids[0].Price,
		Timestamp: depth.Timestamp,
	}

	logger.Debug("save rate", zap.Any("rate", rate))
	// 3. Save the rate to the repository
	if err := s.ratesRepository.SaveRate(ctx, rate); err != nil {
		logger.Error("Failed to save rate", zap.Error(err))
	}

	return rate, nil
}
