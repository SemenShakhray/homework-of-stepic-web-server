-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS users(
    user_id INTEGER PRIMARY KEY,
    username VARCHAR(128) NOT NULL ,
    email VARCHAR(128) NOT NULL,
    password VARCHAR(128) NOT NULL,
    bio TEXT,
    image TEXT,
    token VARCHAR(256),
    following INTEGER DEFAULT 0,
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP NOT NULL,
    updated_at DATETIME DEFAULT CURRENT_TIMESTAMP NOT NULL
); 
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS users;
-- +goose StatementEnd
