CREATE TABLE IF NOT EXISTS users (
    id            BIGINT PRIMARY KEY,  -- Telegram user ID
    first_name    TEXT NOT NULL DEFAULT '',
    username      TEXT NOT NULL DEFAULT '',
    language_code TEXT NOT NULL DEFAULT 'en',
    created_at    TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at    TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE TABLE IF NOT EXISTS customers (
    id         BIGSERIAL PRIMARY KEY,
    user_id    BIGINT NOT NULL REFERENCES users(id),
    name       TEXT NOT NULL,
    phone      TEXT NOT NULL DEFAULT '',
    balance    BIGINT NOT NULL DEFAULT 0,  -- in cents (birr * 100), positive = they owe you
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    UNIQUE(user_id, LOWER(name))
);

CREATE TABLE IF NOT EXISTS products (
    id         BIGSERIAL PRIMARY KEY,
    user_id    BIGINT NOT NULL REFERENCES users(id),
    name       TEXT NOT NULL,
    price      BIGINT NOT NULL DEFAULT 0,  -- default price in cents
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    UNIQUE(user_id, LOWER(name))
);

-- Transaction types: 'debt' (they owe you), 'payment' (they paid you), 'purchase' (they bought something)
CREATE TABLE IF NOT EXISTS transactions (
    id          BIGSERIAL PRIMARY KEY,
    user_id     BIGINT NOT NULL REFERENCES users(id),
    customer_id BIGINT NOT NULL REFERENCES customers(id),
    product_id  BIGINT REFERENCES products(id),
    type        TEXT NOT NULL CHECK (type IN ('debt', 'payment', 'purchase')),
    amount      BIGINT NOT NULL,  -- always positive, in cents
    note        TEXT NOT NULL DEFAULT '',
    created_at  TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_transactions_user_id ON transactions(user_id);
CREATE INDEX idx_transactions_created_at ON transactions(created_at);
CREATE INDEX idx_customers_user_id ON customers(user_id);
CREATE INDEX idx_products_user_id ON products(user_id);
