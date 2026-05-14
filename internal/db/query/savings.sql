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

-- name: GetSavingsOriginalData :one
SELECT starting_capital,
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
SET name = $1,
    starting_capital = $2,
    yearly_interest_rate = $3,
    interest_rate_type = $4,
    monthly_contribution = $5,
    duration_years = $6,
    tax_rate = $7,
    yearly_inflation_rate = $8,
    start_date = $9,
    monthly_interest_rate = $10,
    total_interest_earnings = $11,
    total_deposited = $12,
    rate_of_return = $13,
    inflation_adjusted_ror = $14
WHERE id = $1 AND user_id = $2
RETURNING *;
    

-- name: DeleteSavings :exec
DELETE FROM savings
WHERE id = $1 AND user_id = $2;