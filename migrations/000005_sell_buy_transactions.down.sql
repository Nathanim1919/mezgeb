-- Remove quantity column
ALTER TABLE transactions DROP COLUMN quantity;

-- Restore original type constraint
ALTER TABLE transactions DROP CONSTRAINT transactions_type_check;
ALTER TABLE transactions ADD CONSTRAINT transactions_type_check
    CHECK (type IN ('debt', 'payment', 'purchase'));

-- Restore customer_id NOT NULL
ALTER TABLE transactions ALTER COLUMN customer_id SET NOT NULL;
