-- +goose Up
-- +goose StatementBegin
CREATE TABLE rates (
    id SERIAL PRIMARY KEY,
    market VARCHAR(20) NOT NULL,
    ask DECIMAL(20,8) NOT NULL,
    bid DECIMAL(20,8) NOT NULL,
    timestamp BIGINT NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX idx_rates_market_timestamp ON rates(market, timestamp);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP INDEX IF EXISTS idx_rates_market_timestamp;
DROP TABLE IF EXISTS rates;
-- +goose StatementEnd
