package repository

import (
	"context"
	"usdt-rate-service/internal/models"

	"github.com/jackc/pgx/v5/pgxpool"
)

type Rates struct {
	pool *pgxpool.Pool
}

func NewRates(pool *pgxpool.Pool) *Rates {
	return &Rates{
		pool: pool,
	}
}

func (r *Rates) SaveRate(ctx context.Context, rate *models.Rate) error {
	query := `
	  INSERT INTO rates (market, ask, bid, timestamp) 
	  VALUES ($1, $2, $3, $4)
	`
	_, err := r.pool.Exec(ctx, query, rate.Market, rate.AskPrice, rate.BidPrice, rate.Timestamp)
	return err
}
