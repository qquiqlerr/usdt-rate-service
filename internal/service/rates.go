package service

import (
	"context"
	"usdt-rate-service/internal/models"

	"go.uber.org/zap"
)

// DepthProvider is an interface that defines methods to get depth data for a specific market.
type DepthProvider interface {
	GetDepth(ctx context.Context, market string) (*models.Depth, error)
}

// RatesRepository is an interface that defines methods to save rates to a repository.
type RatesRepository interface {
	SaveRate(ctx context.Context, rate *models.Rate) error
}

// RatesService provides methods to get and save rates for a specific market.
type RatesService struct {
	logger          *zap.Logger
	depthProvider   DepthProvider
	ratesRepository RatesRepository
}

// NewRatesService creates a new RatesService with the provided logger, depth provider, and rates repository.
func NewRatesService(logger *zap.Logger, depthProvider DepthProvider, ratesRepository RatesRepository) *RatesService {
	return &RatesService{
		logger:          logger,
		depthProvider:   depthProvider,
		ratesRepository: ratesRepository,
	}
}

// GetRates retrieves the rates for a specific market by getting depth data from the provider,
// validating it, creating a Rate model, and saving it to the repository.
// It returns the Rate model or an error if any step fails.
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
	if err = depth.Validate(); err != nil {
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
	// TODO: Maybe we dont need to return error here, just log it
	if err = s.ratesRepository.SaveRate(ctx, rate); err != nil {
		logger.Error("Failed to save rate", zap.Error(err))
		return nil, err
	}
	logger.Info("Rate saved successfully", zap.Any("rate", rate))
	return rate, nil
}
