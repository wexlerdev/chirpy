-- +goose Up
CREATE TABLE chirps (
    id UUID PRIMARY KEY,
    created_at TIMESTAMP,
    updated_at TIMESTAMP,
    body VARCHAR(140) NOT NULL,
    user_id UUID REFERENCES users (id) ON DELETE CASCADE
);

-- +goose Down
DROP TABLE chirps;
