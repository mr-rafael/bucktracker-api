-- name: CreateLoan :one
INSERT INTO loans(user_id,
    name,
    starting_principal,
    yearly_interest_rate,
    monthly_payment,
    escrow_payment,
    start_date,
    monthly_interest_rate
)
VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
RETURNING *;

-- name: GetLoansByUserID :many
SELECT id, name, starting_principal FROM loans
WHERE user_id = $1;

-- name: GetLoan :one
SELECT * FROM loans
WHERE id = $1 AND user_id = $2;

-- name: GetLoanInitialData :one
SELECT name,
    starting_principal,
    yearly_interest_rate,
    monthly_payment,
    escrow_payment,
    start_date
FROM loans
WHERE id = $1 AND user_id = $2;

-- name: UpdateLoan :one
UPDATE loans
SET name = $3,
    starting_principal = $4,
    yearly_interest_rate = $5,
    monthly_payment = $6,
    escrow_payment = $7,
    start_date = $8,
    monthly_interest_rate = $9,
    default_payment_plan = $10
WHERE id = $1 AND user_id = $2
RETURNING *;

-- name: DeleteLoan :execrows
DELETE FROM loans
WHERE id = $1 AND user_id = $2;
