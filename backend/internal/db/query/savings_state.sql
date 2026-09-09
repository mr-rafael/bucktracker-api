-- name: CreateSavingsState :one
INSERT INTO savings_state (savings_id,
    date,
    interest,
    tax,
    contribution,
    increase,
    capital
)
VALUES ($1, $2, $3, $4, $5, $6, $7)
RETURNING *;

-- name: GetSavingsStatesBySavingsID :many
SELECT * FROM savings_state
WHERE savings_id = $1
ORDER BY date ASC;

-- name: DeleteSavingsStatesBySavingsID :exec
DELETE FROM savings_state
WHERE savings_id = $1;