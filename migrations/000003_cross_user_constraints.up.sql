-- Add composite unique constraints needed for composite foreign keys
ALTER TABLE customers ADD CONSTRAINT uq_customers_id_user_id UNIQUE (id, user_id);
ALTER TABLE products ADD CONSTRAINT uq_products_id_user_id UNIQUE (id, user_id);

-- Ensure transaction's customer belongs to the same user
ALTER TABLE transactions
    ADD CONSTRAINT fk_transactions_customer_owner
    FOREIGN KEY (customer_id, user_id) REFERENCES customers(id, user_id);

-- Ensure transaction's product (when set) belongs to the same user
-- We need a partial approach since product_id is nullable.
-- Create a function + trigger for this.
CREATE OR REPLACE FUNCTION check_transaction_product_owner()
RETURNS TRIGGER AS $$
BEGIN
    IF NEW.product_id IS NOT NULL THEN
        IF NOT EXISTS (
            SELECT 1 FROM products WHERE id = NEW.product_id AND user_id = NEW.user_id
        ) THEN
            RAISE EXCEPTION 'product does not belong to user';
        END IF;
    END IF;
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER trg_transaction_product_owner
    BEFORE INSERT OR UPDATE ON transactions
    FOR EACH ROW
    EXECUTE FUNCTION check_transaction_product_owner();
