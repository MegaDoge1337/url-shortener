-- +goose Up
CREATE TABLE url (
    id    SERIAL PRIMARY KEY,
    url   VARCHAR NOT NULL,
    alias VARCHAR NOT NULL
);

-- +goose Down
DROP TABLE url;
