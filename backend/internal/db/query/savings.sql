-- name: CreateSavings :one
INSERT INTO savings (user_id,
    name,
    starting_capital,
    yearly_interest_rate,
    interest_rate_type,
    monthly_contribution,
    duration_years,
    tax_rate,
    yearly_inflation_rate,
    start_date,
    monthly_interest_rate,
    total_interest_earnings,
    total_deposited,
    rate_of_return,
    inflation_adjusted_ror
)
VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15)
RETURNING *;

-- name: GetSavingsByUserID :many
SELECT id, name, starting_capital FROM savings
WHERE user_id = $1;

-- name: GetSavings :one
SELECT * FROM savings
WHERE id = $1 AND user_id = $2;

-- name: GetSavingsInitialData :one
SELECT name,
    starting_capital,
    yearly_interest_rate,
    interest_rate_type,
    monthly_contribution,
    duration_years,
    tax_rate,
    yearly_inflation_rate,
    start_date
FROM savings
WHERE id = $1 AND user_id = $2;

-- name: UpdateSavings :one
UPDATE savings
SET name = $3,
    starting_capital = $4,
    yearly_interest_rate = $5,
    interest_rate_type = $6,
    monthly_contribution = $7,
    duration_years = $8,
    tax_rate = $9,
    yearly_inflation_rate = $10,
    start_date = $11,
    monthly_interest_rate = $12,
    total_interest_earnings = $13,
    total_deposited = $14,
    rate_of_return = $15,
    inflation_adjusted_ror = $16
WHERE id = $1 AND user_id = $2
RETURNING *;
    

-- name: DeleteSavings :execrows
DELETE FROM savings
WHERE id = $1 AND user_id = $2;