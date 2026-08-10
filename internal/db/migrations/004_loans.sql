-- +goose Up
CREATE TABLE loans (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id UUID REFERENCES users(id) ON DELETE CASCADE,
    name TEXT NOT NULL,
    starting_principal INT NOT NULL,
    yearly_interest_rate TEXT NOT NULL,
    monthly_payment INT NOT NULL,
    escrow_payment INT NOT NULL,
    start_date TIMESTAMPTZ NOT NULL,
    monthly_interest_rate TEXT NOT NULL,
    default_payment_plan UUID,
    created_at TIMESTAMPTZ DEFAULT now()
);

CREATE TABLE payment_plans (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    loan_id UUID NOT NULL REFERENCES loans(id) ON DELETE CASCADE,
    name TEXT NOT NULL DEFAULT 'Default Payment Plan',
    duration_months INT NOT NULL,
    total_expenditure INT NOT NULL,
    total_paid INT NOT NULL,
    cost_of_credit TEXT NOT NULL
);

ALTER TABLE loans
    ADD CONSTRAINT loans_default_payment_plan_fkey
    FOREIGN KEY (default_payment_plan) REFERENCES payment_plans(id);

CREATE TABLE principal_payments (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    payment_plan_id UUID NOT NULL REFERENCES payment_plans(id) ON DELETE CASCADE,
    amount_paid INT NOT NULL,
    date TIMESTAMPTZ NOT NULL
);

CREATE TABLE loan_state (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    payment_plan_id UUID REFERENCES payment_plans(id) ON DELETE CASCADE,
    date TIMESTAMPTZ NOT NULL,
    payment INT NOT NULL,
    interest INT NOT NULL,
    other_payments INT NOT NULL,
    paydown INT NOT NULL,
    principal INT NOT NULL
);

-- +goose Down
DROP TABLE loan_state;
DROP TABLE principal_payments;
ALTER TABLE loans DROP CONSTRAINT loans_default_payment_plan_fkey;
DROP TABLE payment_plans;
DROP TABLE loans;
