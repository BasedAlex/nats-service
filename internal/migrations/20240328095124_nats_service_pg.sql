-- +goose Up
-- +goose StatementBegin
CREATE TABLE orders (
  order_id VARCHAR(255) PRIMARY KEY,
  order_data jsonb
);

-- +goose StatementEnd
    
-- +goose Down
-- +goose StatementBegin
DROP TABLE orders;
-- +goose StatementEnd
