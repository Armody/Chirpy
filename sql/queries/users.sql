-- name: CreateUser :one
INSERT INTO users (id, created_at, updated_at, email, hashed_password)
VALUES (
    gen_random_uuid(),
    NOW(),
    NOW(),
    $1,
    $2
)
RETURNING *;

-- name: DeleteUsers :exec
DELETE FROM users;

-- name: GetUserByEmail :one
SELECT * FROM users
WHERE email = $1;

-- name: GetUserById :one
SELECT * FROM users
WHERE id = $1;

-- name: UpdateEmail :exec
UPDATE users
SET updated_at = NOW(),
    email = $2
WHERE id = $1;

-- name: UpdatePassword :exec
UPDATE users
SET updated_at = NOW(),
    hashed_password = $2
WHERE id = $1;

-- name: UpgradeToChirpyRed :one
UPDATE users
SET is_chirpy_red = true, updated_at = NOW()
WHERE id = $1
RETURNING *;