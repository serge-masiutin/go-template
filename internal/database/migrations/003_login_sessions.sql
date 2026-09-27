-- +goose Up
CREATE TABLE login_sessions (
    token_hash text PRIMARY KEY CHECK (token_hash ~ '^[0-9a-f]{64}$'),
    user_id bigint NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    expires_at timestamptz NOT NULL
);
CREATE INDEX login_sessions_expires_at_idx ON login_sessions(expires_at);
CREATE INDEX login_sessions_user_id_idx ON login_sessions(user_id);

-- +goose Down
DROP TABLE login_sessions;
