-- +goose Up
CREATE TABLE assets (
    id BIGSERIAL PRIMARY KEY,
    symbol TEXT NOT NULL,
    name TEXT NOT NULL,
    instrument_type TEXT NOT NULL CHECK (
        instrument_type IN ('stock', 'mutual_fund')
    ),
    isin TEXT UNIQUE NOT NULL,
    exchange TEXT NOT NULL,
    currency TEXT DEFAULT 'INR',
    external_platform_id TEXT UNIQUE,
    created_at TIMESTAMPTZ DEFAULT now()
);
CREATE INDEX idx_assets_isin ON assets(isin);
-- +goose StatementBegin
SELECT 'up SQL query';
-- +goose StatementEnd
-- +goose Down
DROP INDEX IF EXISTS idx_assets_isin;
DROP TABLE IF EXISTS assets;
-- +goose StatementBegin
SELECT 'down SQL query';
-- +goose StatementEnd