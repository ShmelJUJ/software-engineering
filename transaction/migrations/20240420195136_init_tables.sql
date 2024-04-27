-- +goose Up
CREATE TABLE IF NOT EXISTS transaction_users (
    transaction_user_id UUID PRIMARY KEY NOT NULL,
    user_id UUID NOT NULL,
    created_at TIMESTAMP NOT NULL,
    updated_at TIMESTAMP NOT NULL
);

CREATE TABLE IF NOT EXISTS transactions (
    transaction_id UUID PRIMARY KEY NOT NULL,
    sender_id UUID NULL,
    receiver_id UUID NOT NULL,
    currency TEXT NOT NULL,
    amount BIGINT NOT NULL,
    status INT NOT NULL,
    method TEXT NOT NULL,
    canceled_reason TEXT,
    created_at TIMESTAMP NOT NULL,
    updated_at TIMESTAMP NOT NULL,

    FOREIGN KEY (sender_id) REFERENCES transaction_users(transaction_user_id) ON UPDATE CASCADE ON DELETE SET NULL,
    FOREIGN KEY (receiver_id) REFERENCES transaction_users(transaction_user_id) ON UPDATE CASCADE ON DELETE SET NULL
);

-- +goose StatementBegin
SELECT 'up SQL query';
-- +goose StatementEnd

-- +goose Down
DROP TABLE IF EXISTS transactions;

DROP TABLE IF EXISTS transaction_users;

-- +goose StatementBegin
SELECT 'down SQL query';
-- +goose StatementEnd
