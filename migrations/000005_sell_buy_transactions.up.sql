-- Add quantity column to transactions
ALTER TABLE transactions ADD COLUMN quantity BIGINT NOT NULL DEFAULT 0;

-- Allow sell/buy transaction types
ALTER TABLE transactions DROP CONSTRAINT transactions_type_check;
ALTER TABLE transactions ADD CONSTRAINT transactions_type_check
    CHECK (type IN ('debt', 'payment', 'purchase', 'sell', 'buy'));

-- Make customer_id nullable (sell/buy don't require a customer)
ALTER TABLE transactions ALTER COLUMN customer_id DROP NOT NULL;
