-- +goose Up
CREATE TABLE holdings (
    id BIGSERIAL PRIMARY KEY,
    user_asset_id BIGINT NOT NULL REFERENCES user_assets(id) ON DELETE CASCADE,
    total_quantity NUMERIC NOT NULL DEFAULT 0,
    average_price NUMERIC NOT NULL DEFAULT 0,
    total_invested NUMERIC NOT NULL DEFAULT 0,
    updated_at TIMESTAMPTZ DEFAULT now(),
    UNIQUE(user_asset_id)
);
-- +goose StatementBegin
SELECT 'up SQL query';
-- +goose StatementEnd
-- +goose Down
DROP TABLE IF EXISTS holdings;
-- +goose StatementBegin
SELECT 'down SQL query';
-- +goose StatementEnd