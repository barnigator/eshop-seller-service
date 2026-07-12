CREATE EXTENSION IF NOT EXISTS pgcrypto;

CREATE TABLE sellers (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id UUID NOT NULL,
    brand_name VARCHAR(120) NOT NULL,
    description TEXT NOT NULL DEFAULT '',
    status VARCHAR(32) NOT NULL DEFAULT 'pending',
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    deleted_at TIMESTAMPTZ,

    CONSTRAINT sellers_status_valid
        CHECK (
            status IN (
                'pending',
                'active',
                'blocked',
                'archived'
            )
        )
);

CREATE INDEX idx_sellers_user_id
    ON sellers(user_id)
    WHERE deleted_at IS NULL;

CREATE INDEX idx_sellers_status
    ON sellers(status)
    WHERE deleted_at IS NULL;

