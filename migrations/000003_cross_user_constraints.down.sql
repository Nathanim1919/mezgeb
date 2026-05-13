DROP TRIGGER IF EXISTS trg_transaction_product_owner ON transactions;
DROP FUNCTION IF EXISTS check_transaction_product_owner();
ALTER TABLE transactions DROP CONSTRAINT IF EXISTS fk_transactions_customer_owner;
ALTER TABLE products DROP CONSTRAINT IF EXISTS uq_products_id_user_id;
ALTER TABLE customers DROP CONSTRAINT IF EXISTS uq_customers_id_user_id;
