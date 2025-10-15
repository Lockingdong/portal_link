-- +goose Up
-- +goose StatementBegin
CREATE TABLE users (
    id SERIAL PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    email VARCHAR(255) NOT NULL UNIQUE,
    password VARCHAR(255) NOT NULL,
    created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP DEFAULT NOW()
);

-- Create indexes
CREATE UNIQUE INDEX idx_users_email ON users (email);
CREATE INDEX idx_users_created_at ON users (created_at);

-- Add comments
COMMENT ON TABLE users IS '使用者資料表';
COMMENT ON COLUMN users.id IS '使用者的唯一標識符';
COMMENT ON COLUMN users.name IS '使用者的全名';
COMMENT ON COLUMN users.email IS '使用者的電子郵件地址，必須是唯一的';
COMMENT ON COLUMN users.password IS '使用者的密碼';
COMMENT ON COLUMN users.created_at IS '建立時間 UTC';
COMMENT ON COLUMN users.updated_at IS '更新時間 UTC';
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP INDEX IF EXISTS idx_users_created_at;
DROP INDEX IF EXISTS idx_users_email;
DROP TABLE IF EXISTS users;
-- +goose StatementEnd
