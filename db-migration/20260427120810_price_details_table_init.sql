-- +goose Up
CREATE TABLE price_details (
    id BIGSERIAL PRIMARY KEY,
    asset_id BIGINT NOT NULL UNIQUE REFERENCES assets(id) ON DELETE CASCADE,
    curr_price NUMERIC NOT NULL,
    prev_price NUMERIC NOT NULL,
    updated_at TIMESTAMPTZ DEFAULT now()
);
-- +goose StatementBegin
SELECT 'up SQL query';
-- +goose StatementEnd
-- +goose Down
DROP TABLE IF EXISTS price_details;
-- +goose StatementBegin
SELECT 'down SQL query';
-- +goose StatementEnd