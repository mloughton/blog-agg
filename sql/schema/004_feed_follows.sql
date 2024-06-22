-- +goose Up
CREATE TABLE feed_follows (
    id UUID PRIMARY KEY,
    feed_id UUID REFERENCES feeds(id) ON UPDATE CASCADE ON DELETE CASCADE NOT NULL,
    user_id UUID REFERENCES users(id) ON UPDATE CASCADE ON DELETE CASCADE NOT NULL,
    created_at TIMESTAMP NOT NULL,
    updated_at TIMESTAMP NOT NULL,
    UNIQUE (feed_id, user_id)
);

-- +goose Down
DROP TABLE feed_follows;