-- name: CreateChirp :one
INSERT INTO chirps (id, created_at, updated_at, user_id, body)
VALUES (
    gen_random_uuid(),
    NOW(),
    NOW(),
    $1,
    $2
)
RETURNING *;

-- name: GetAllChirpsAsc :many
SELECT * FROM chirps
ORDER BY created_at ASC;

-- name: GetAllChirpsDesc :many
SELECT * FROM chirps
ORDER BY created_at DESC;

-- name: GetChirpsByAuthorIdAsc :many
SELECT * FROM chirps
WHERE user_id = $1
ORDER BY created_at ASC;

-- name: GetChirpsByAuthorIdDesc :many
SELECT * FROM chirps
WHERE user_id = $1
ORDER BY created_at DESC;

-- name: GetChirp :one
SELECT * FROM chirps
WHERE id = $1;

-- name: DeleteChirp :one
DELETE FROM chirps
WHERE id = $1
RETURNING *;

