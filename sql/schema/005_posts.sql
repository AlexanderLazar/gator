-- +goose up

CREATE TABLE posts (
    id uuid PRIMARY KEY NOT NULL,
    created_at TIMESTAMP NOT NULL,
    updated_at TIMESTAMP NOT NULL,
    title TEXT NOT NULL,
    url TEXT UNIQUE NOT NULL,
    description TEXT NOT NULL,
    published_at TIMESTAMP NOT NULL,
    feed_id uuid NOT NULL,
    FOREIGN KEY (feed_id) REFERENCES feeds (id)
);

-- +goose down
DROP TABLE posts;
