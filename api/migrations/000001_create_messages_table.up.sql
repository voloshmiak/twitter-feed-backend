CREATE TABLE IF NOT EXISTS messages (
        id UUID PRIMARY KEY,
        user_id STRING NOT NULL,
        content STRING NOT NULL,
        created_at TIMESTAMPTZ NOT NULL
);