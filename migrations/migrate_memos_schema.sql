-- +migrate Up
CREATE TABLE IF NOT EXISTS memos (
    id         INT GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
    title      VARCHAR(255) NOT NULL,
    context    TEXT,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
    )

-- +migrate Down
DROP TABLE memos;