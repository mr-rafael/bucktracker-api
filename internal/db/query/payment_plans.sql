-- name: CreatePaymentPlan :one
INSERT INTO payment_plans (
    loan_id,
    name,
    duration_months,
    total_expenditure,
    total_paid,
    cost_of_credit
)
VALUES ($1, $2, $3, $4, $5, $6)
RETURNING *;

-- name: GetPaymentPlan :one
SELECT * FROM payment_plans
WHERE id = $1;

-- name: GetPaymentPlanByIDAndLoanID :one
SELECT * FROM payment_plans
WHERE id = $1 AND loan_id = $2;

-- name: GetPaymentPlansByLoanID :many
SELECT * FROM payment_plans
WHERE loan_id = $1;

-- name: UpdatePaymentPlan :one
UPDATE payment_plans
SET name = $2,
    duration_months = $3,
    total_expenditure = $4,
    total_paid = $5,
    cost_of_credit = $6
WHERE id = $1
RETURNING *;

-- name: DeletePaymentPlan :exec
DELETE FROM payment_plans
WHERE id = $1;
