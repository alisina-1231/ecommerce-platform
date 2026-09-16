CREATE TABLE IF NOT EXISTS products (
    id BIGSERIAL PRIMARY KEY,

    name VARCHAR(255) NOT NULL,

    description TEXT NOT NULL DEFAULT '',

    price NUMERIC(12, 2) NOT NULL
        CHECK (price >= 0),

    currency CHAR(3) NOT NULL
        DEFAULT 'USD',

    created_at TIMESTAMPTZ NOT NULL
        DEFAULT NOW(),

    updated_at TIMESTAMPTZ NOT NULL
        DEFAULT NOW()
);

CREATE INDEX IF NOT EXISTS idx_products_name
    ON products (name);