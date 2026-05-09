-- name: CreateLoanState :one
INSERT INTO loan_state (loan_id,
    date,
    payment,
    interest,
    other_payments,
    paydown,
    principal
)
VALUES ($1, $2, $3, $4, $5, $6, $7)
RETURNING *;

-- name: GetLoanStatesByLoanID :many
SELECT * FROM loan_state
WHERE loan_id = $1
ORDER BY date ASC;

-- name: DeleteLoanStatesByLoanID :exec
DELETE FROM loan_state
WHERE loan_id = $1;