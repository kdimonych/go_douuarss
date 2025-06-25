-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS channels (
    id SERIAL PRIMARY KEY,
    title TEXT,
    link TEXT,
    description TEXT,
    language TEXT,
    last_build_date TIMESTAMPTZ
);

CREATE TABLE IF NOT EXISTS items (
    id SERIAL PRIMARY KEY,
    title TEXT NOT NULL,
    link TEXT,
    description TEXT NOT NULL,
    author TEXT,
    category TEXT,
    pub_date TIMESTAMPTZ,
    channel_id INTEGER NOT NULL REFERENCES channels(id) ON DELETE CASCADE
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE items;
DROP TABLE channels;
-- +goose StatementEnd
