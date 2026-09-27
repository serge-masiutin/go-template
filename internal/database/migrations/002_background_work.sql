-- +goose Up
CREATE TABLE note_emails (
    id bigint GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
    user_id bigint NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    job_id bigint UNIQUE,
    state text NOT NULL DEFAULT 'queued' CHECK (state IN ('queued', 'sending', 'sent', 'failed')),
    created_at timestamptz NOT NULL DEFAULT now(),
    finished_at timestamptz
);
CREATE INDEX note_emails_user_id_id_idx ON note_emails(user_id, id DESC);
CREATE UNIQUE INDEX note_emails_one_active_idx ON note_emails(user_id) WHERE state IN ('queued', 'sending');

CREATE TABLE assistant_runs (
    id bigint GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
    user_id bigint NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    job_id bigint UNIQUE,
    question text NOT NULL CHECK (length(btrim(question)) BETWEEN 1 AND 500),
    state text NOT NULL DEFAULT 'queued' CHECK (state IN ('queued', 'running', 'completed', 'failed')),
    answer text NOT NULL DEFAULT '',
    model text NOT NULL,
    prompt_version text NOT NULL,
    created_at timestamptz NOT NULL DEFAULT now(),
    finished_at timestamptz
);
CREATE INDEX assistant_runs_user_id_id_idx ON assistant_runs(user_id, id DESC);
CREATE UNIQUE INDEX assistant_runs_one_active_idx ON assistant_runs(user_id) WHERE state IN ('queued', 'running');

-- +goose Down
DROP TABLE assistant_runs;
DROP TABLE note_emails;
