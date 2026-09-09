-- name: CreateRefreshToken :one
INSERT INTO refresh_tokens (user_id, token_hash, expires_at, revoked)
VALUES ($1, $2, $3, $4)
RETURNING *;

-- name: GetTokenByHash :one
SELECT * FROM refresh_tokens
WHERE token_hash = $1;

-- name: RevokeTokenByUserID :exec
UPDATE refresh_tokens
SET revoked = TRUE
WHERE user_id = $1;