CREATE TABLE IF NOT EXISTS orders (
    id BIGSERIAL PRIMARY KEY,

    user_id VARCHAR(255) NOT NULL,

    total_amount NUMERIC(12, 2) NOT NULL
        CHECK (total_amount > 0),

    currency VARCHAR(10) NOT NULL DEFAULT 'USD',

    status VARCHAR(20) NOT NULL DEFAULT 'PENDING',

    shipping_address TEXT NOT NULL,

    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),

    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),

    CONSTRAINT orders_status_check
        CHECK (
            status IN (
                'PENDING',
                'CONFIRMED',
                'PROCESSING',
                'SHIPPED',
                'DELIVERED',
                'CANCELLED'
            )
        )
);

CREATE INDEX IF NOT EXISTS idx_orders_user_id
    ON orders(user_id);

CREATE INDEX IF NOT EXISTS idx_orders_status
    ON orders(status);

CREATE INDEX IF NOT EXISTS idx_orders_created_at
    ON orders(created_at DESC);