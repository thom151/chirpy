-- name: CreateUser :one
INSERT INTO users (id, created_at, updated_at, email)
VALUES (
    gen_random_uuid(),
    NOW(),
    NOW(),
    $1
)
RETURNING *;

-- name: DeleteUsers :exec
DELETE FROM users;

-- name: SetPassword :exec
UPDATE users SET hashed_password = $1, updated_at = CURRENT_TIMESTAMP where id = $2;


-- name: GetUser :one
SELECT * FROM users WHERE email = $1;

-- name: GetUserByID :one
SELECT id, created_at, updated_at, email FROM users WHERE id = $1;

-- name: UpdateUser :exec
UPDATE users SET hashed_password = $1, email = $2, updated_at = CURRENT_TIMESTAMP WHERE id = $3;

-- name: UpgradeUser :exec
UPDATE users SET is_chirpy_red = TRUE, updated_at = CURRENT_TIMESTAMP WHERE id = $1;
