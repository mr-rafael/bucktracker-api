-- name: CreateLoanState :one
INSERT INTO loan_state (payment_plan_id,
    date,
    payment,
    interest,
    other_payments,
    paydown,
    principal
)
VALUES ($1, $2, $3, $4, $5, $6, $7)
RETURNING *;

-- name: GetLoanStatesByPaymentPlanID :many
SELECT * FROM loan_state
WHERE payment_plan_id = $1
ORDER BY date ASC;

-- name: DeleteLoanStatesByPaymentPlanID :exec
DELETE FROM loan_state
WHERE payment_plan_id = $1;
