-- name: CreateRefreshToken :one
INSERT INTO refresh_tokens (token, created_at, updated_at, user_id, expires_at, revoked_at)
VALUES (
    $1,
    NOW(),
    NOW(),
    $2,
    $3,
    NULL
)
RETURNING *;

-- name: GetRefreshToken :one
SELECT * FROM refresh_tokens
WHERE $1 = token;

-- name: GetUserFromRefreshToken :one
SELECT users.*
FROM users
INNER JOIN refresh_tokens ON users.id = refresh_tokens.user_id
WHERE refresh_tokens.token = $1;

-- name: RevokeRefreshToken :exec
UPDATE refresh_tokens
SET
updated_at = NOW(),
revoked_at = NOW()
WHERE $1 = token;


