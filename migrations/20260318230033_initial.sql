-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS users (
    id UUID PRIMARY KEY,
    username TEXT NOT NULL
);

CREATE TABLE IF NOT EXISTS auth_identities (
    id UUID PRIMARY KEY,
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    provider INTEGER NOT NULL,
    provider_id TEXT NOT NULL DEFAULT '',
    email TEXT NOT NULL UNIQUE,
    password_hash TEXT NOT NULL,
    CONSTRAINT auth_identities_local_provider_only CHECK (provider = 0),
    CONSTRAINT auth_identities_local_provider_id_empty CHECK (provider_id = '')
);

CREATE INDEX IF NOT EXISTS idx_auth_identities_user_id ON auth_identities (user_id);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP INDEX IF EXISTS idx_auth_identities_user_id;
DROP TABLE IF EXISTS auth_identities;
DROP TABLE IF EXISTS users;
-- +goose StatementEnd
