CREATE TABLE users (
    user_id UUID PRIMARY KEY,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
);
CREATE TABLE memos (
    memo_id    BIGSERIAL PRIMARY KEY,
    user_id    UUID NOT NULL REFERENCES users(user_id),
    body       TEXT NOT NULL DEFAULT '',
    mood       TEXT NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT now()
);