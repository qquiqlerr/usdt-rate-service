package adapter

import (
	"context"
	"usdt-rate-service/internal/infra/grinex"
	"usdt-rate-service/internal/models"
)

// GrinexDepthProvider is an implementation of the DepthProvider interface that uses the Grinex client to get depth data.
type GrinexDepthProvider struct {
	client *grinex.Client
}

// NewGrinexDepthProvider creates a new GrinexDepthProvider with the provided Grinex client.
func NewGrinexDepthProvider(client *grinex.Client) *GrinexDepthProvider {
	return &GrinexDepthProvider{client: client}
}

// GetDepth retrieves the depth data for a specific market using the Grinex client.
func (p *GrinexDepthProvider) GetDepth(ctx context.Context, market string) (*models.Depth, error) {
	dto, err := p.client.GetDepth(ctx, market)
	if err != nil {
		return nil, err
	}
	return ToDepth(dto), nil
}

// ToDepth converts a Grinex DepthResponse DTO to a models.Depth.
func ToDepth(dto *grinex.DepthResponse) *models.Depth {
	asks := make([]models.Order, len(dto.Asks))
	for i, ask := range dto.Asks {
		asks[i] = models.Order{
			Price:  ask.Price,
			Volume: ask.Volume,
			Amount: ask.Amount,
			Factor: ask.Factor,
			Type:   ask.Type,
		}
	}

	bids := make([]models.Order, len(dto.Bids))
	for i, bid := range dto.Bids {
		bids[i] = models.Order{
			Price:  bid.Price,
			Volume: bid.Volume,
			Amount: bid.Amount,
			Factor: bid.Factor,
			Type:   bid.Type,
		}
	}

	return &models.Depth{
		Timestamp: dto.Timestamp,
		Asks:      asks,
		Bids:      bids,
	}
}
