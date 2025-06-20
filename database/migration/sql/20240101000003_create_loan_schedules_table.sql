-- +goose Up
-- +goose StatementBegin
CREATE TABLE loan_schedules (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    loan_id UUID NOT NULL REFERENCES loans(id) ON DELETE CASCADE,
    due_date TIMESTAMP WITH TIME ZONE NOT NULL,
    basic_amount DECIMAL(15,2) NOT NULL,
    interest_amount DECIMAL(15,2) NOT NULL,
    total_payment DECIMAL(15,2) NOT NULL,
    status VARCHAR(50) NOT NULL DEFAULT 'pending',
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX idx_loan_schedules_loan_id ON loan_schedules(loan_id);
CREATE INDEX idx_loan_schedules_status ON loan_schedules(status);
CREATE INDEX idx_loan_schedules_due_date ON loan_schedules(due_date);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS loan_schedules;
-- +goose StatementEnd 