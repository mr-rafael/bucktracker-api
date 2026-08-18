-- +goose Up
CREATE TABLE savings (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id UUID REFERENCES users(id) ON DELETE CASCADE,
    name TEXT NOT NULL,
    starting_capital INT NOT NULL,
    yearly_interest_rate TEXT NOT NULL,
    interest_rate_type TEXT NOT NULL,
    monthly_contribution INT NOT NULL,
    duration_years INT NOT NULL,
    tax_rate TEXT NOT NULL,
    yearly_inflation_rate TEXT,
    start_date DATE NOT NULL,
    monthly_interest_rate TEXT NOT NULL,
    total_interest_earnings INT NOT NULL,
    rate_of_return TEXT NOT NULL,
    inflation_adjusted_ror TEXT NOT NULL,
    total_deposited INT NOT NULL DEFAULT 0,
    created_at TIMESTAMPTZ DEFAULT now()
);

CREATE TABLE savings_state (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    savings_id UUID REFERENCES savings(id) ON DELETE CASCADE,
    date DATE NOT NULL,
    interest INT NOT NULL,
    tax INT NOT NULL,
    contribution INT NOT NULL,
    increase INT NOT NULL,
    capital INT NOT NULL
);

-- +goose Down
DROP TABLE savings_state;
DROP TABLE savings;