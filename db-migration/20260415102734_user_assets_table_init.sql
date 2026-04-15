-- +goose Up
CREATE TABLE user_assets (
    id BIGSERIAL PRIMARY KEY,
    user_id BIGINT NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    asset_id BIGINT NOT NULL REFERENCES assets(id) ON DELETE CASCADE,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    UNIQUE(user_id, asset_id)
);
-- +goose StatementBegin
SELECT 'up SQL query';
-- +goose StatementEnd
-- +goose Down
DROP TABLE IF EXISTS user_assets;
-- +goose StatementBegin
SELECT 'down SQL query';
-- +goose StatementEnd