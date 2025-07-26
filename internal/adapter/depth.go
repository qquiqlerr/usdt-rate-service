package adapter

import (
	"context"
	"usdt-rate-service/internal/infra/grinex"
	"usdt-rate-service/internal/models"
)

type GrinexDepthProvider struct {
	client *grinex.Client
}

func NewGrinexDepthProvider(client *grinex.Client) *GrinexDepthProvider {
	return &GrinexDepthProvider{client: client}
}

func (p *GrinexDepthProvider) GetDepth(ctx context.Context, market string) (*models.Depth, error) {
	dto, err := p.client.GetDepth(ctx, market)
	if err != nil {
		return nil, err
	}
	return ToDepth(dto), nil
}

func ToDepth(dto *grinex.DepthResponse) *models.Depth {
	asks := make([]models.Order, len(dto.Asks))
	for i, ask := range dto.Asks {
		asks[i] = models.Order{
			Price:  ask.Price,
			Volume: ask.Volume,
		}
	}

	bids := make([]models.Order, len(dto.Bids))
	for i, bid := range dto.Bids {
		bids[i] = models.Order{
			Price:  bid.Price,
			Volume: bid.Volume,
		}
	}

	return &models.Depth{
		Timestamp: dto.Timestamp,
		Asks:      asks,
		Bids:      bids,
	}
}
