-- +goose Up
-- Create operation_type enum
CREATE TYPE operation_type AS ENUM (
    'NORMAL_PURCHASE',
    'PURCHASE_WITH_INSTALLMENTS',
    'WITHDRAWAL',
    'CREDIT_VOUCHER'
);

-- Create accounts table
CREATE TABLE accounts (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    document_number VARCHAR(20) NOT NULL UNIQUE,
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP NOT NULL DEFAULT NOW()
);

-- Create index on document_number for faster lookups
CREATE INDEX idx_accounts_document_number ON accounts(document_number);

-- Create transactions table
CREATE TABLE transactions (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    account_id UUID NOT NULL,
    amount NUMERIC(15, 2) NOT NULL,
    event_date TIMESTAMP NOT NULL DEFAULT NOW(),
    operation_type operation_type NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    CONSTRAINT fk_account
        FOREIGN KEY (account_id)
        REFERENCES accounts(id)
        ON DELETE CASCADE
);

-- Create indexes for better query performance
CREATE INDEX idx_transactions_account_id ON transactions(account_id);
CREATE INDEX idx_transactions_event_date ON transactions(event_date);
CREATE INDEX idx_transactions_operation_type ON transactions(operation_type);

-- +goose Down
-- Drop tables (foreign keys will be dropped automatically)
DROP TABLE IF EXISTS transactions;
DROP TABLE IF EXISTS accounts;

-- Drop the enum type
DROP TYPE IF EXISTS operation_type;
