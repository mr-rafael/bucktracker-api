-- name: CreatePrincipalPayment :one
INSERT INTO principal_payments (
    payment_plan_id,
    amount_paid,
    date
)
VALUES ($1, $2, $3)
RETURNING *;

-- name: GetPrincipalPaymentsByPaymentPlanID :many
SELECT * FROM principal_payments
WHERE payment_plan_id = $1
ORDER BY date ASC;

-- name: DeletePrincipalPaymentsByPaymentPlanID :exec
DELETE FROM principal_payments
WHERE payment_plan_id = $1;
