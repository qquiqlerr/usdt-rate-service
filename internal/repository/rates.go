package repository

import (
	"context"
	"usdt-rate-service/internal/models"

	"github.com/jackc/pgx/v5/pgxpool"
)

// Rates is a repository for managing rates in the database.
type Rates struct {
	pool *pgxpool.Pool
}

// NewRates creates a new Rates repository with the provided database connection pool.
func NewRates(pool *pgxpool.Pool) *Rates {
	return &Rates{
		pool: pool,
	}
}

// SaveRate saves a Rate model to the database.
// It inserts the market, ask price, bid price, and timestamp into the rates table.
func (r *Rates) SaveRate(ctx context.Context, rate *models.Rate) error {
	query := `
	  INSERT INTO rates (market, ask, bid, timestamp) 
	  VALUES ($1, $2, $3, $4)
	`
	_, err := r.pool.Exec(ctx, query, rate.Market, rate.AskPrice, rate.BidPrice, rate.Timestamp)
	return err
}
