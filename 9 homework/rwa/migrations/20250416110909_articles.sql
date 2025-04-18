-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS articles (
    article_id INTEGER PRIMARY KEY AUTOINCREMENT,
    slug TEXT,
    title TEXT NOT NULL,
    description TEXT NOT NULL,
    body TEXT NOT NULL,
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    favorites_count INTEGER DEFAULT 0,
    author_id INTEGER NOT NULL,
    FOREIGN KEY (author_id) REFERENCES users(user_id)
     ON DELETE CASCADE
);
CREATE INDEX IF NOT EXISTS articles_author_id ON articles(author_id);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP INDEX IF EXISTS articles_author_id;
DROP TABLE IF EXISTS articles;
-- +goose StatementEnd
