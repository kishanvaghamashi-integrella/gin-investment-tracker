-- +goose Up
CREATE TABLE transactions (
    id BIGSERIAL PRIMARY KEY,
    user_asset_id BIGINT NOT NULL REFERENCES user_assets(id) ON DELETE CASCADE,
    txn_type TEXT NOT NULL CHECK (txn_type IN ('BUY', 'SELL')),
    quantity NUMERIC NOT NULL,
    price NUMERIC NOT NULL,
    txn_date TIMESTAMPTZ NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now()
);
-- +goose StatementBegin
SELECT 'up SQL query';
-- +goose StatementEnd
-- +goose Down
DROP TABLE IF EXISTS transactions;
-- +goose StatementBegin
SELECT 'down SQL query';
-- +goose StatementEnd