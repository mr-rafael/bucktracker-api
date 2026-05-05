-- +goose Up
ALTER TABLE savings
ADD total_deposited INT NOT NULL DEFAULT 0;

-- +goose Down
ALTER TABLE savings
DROP COLUMN total_deposited