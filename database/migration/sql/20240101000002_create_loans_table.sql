-- +goose Up
-- +goose StatementBegin
CREATE TABLE loans (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    borrower_id UUID NOT NULL REFERENCES borrowers(id) ON DELETE CASCADE,
    amount DECIMAL(15,2) NOT NULL,
    repayment_cadence_days INTEGER NOT NULL,
    interest_percentage DECIMAL(5,2) NOT NULL,
    interest_amount DECIMAL(15,2) NOT NULL,
    status VARCHAR(50) NOT NULL DEFAULT 'active',
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX idx_loans_borrower_id ON loans(borrower_id);
CREATE INDEX idx_loans_status ON loans(status);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS loans;
-- +goose StatementEnd 