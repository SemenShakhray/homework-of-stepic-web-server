-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS tags(
    tags_id INTEGER PRIMARY KEY AUTOINCREMENT,
    name VARCHAT(26) UNIQUE NOT NULL
);

CREATE TABLE IF NOT EXISTS article_tags(
    article_id INTEGER NOT NULL
    tags_id INTEGER NOT NULL,
    PRIMARY KEY (article_id, tag_id),
    FOREIGN KEY (article_id) REFERENCES articles(article_id) ON DELETE CASCADE,
    FOREIGN KEY (tag_id) REFERENCES tags(tag_id) ON DELETE CASCADE
)
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
SELECT 'down SQL query';
-- +goose StatementEnd
