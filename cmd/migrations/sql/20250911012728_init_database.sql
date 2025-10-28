-- +goose Up
-- +goose StatementBegin
CREATE TABLE "users" (
    id UUID DEFAULT uuidv7() PRIMARY KEY,
    first_name VARCHAR NOT NULL,
    last_name VARCHAR NOT NULL,
    email TEXT NOT NULL UNIQUE,
    keycloak_id TEXT UNIQUE,
    created_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMPTZ
);

CREATE TABLE "companies" (
    id UUID DEFAULT uuidv7() PRIMARY KEY,
    name VARCHAR NOT NULL,
    keycloak_id TEXT UNIQUE,
    created_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMPTZ
);

CREATE TABLE "user_companies" (
    id UUID DEFAULT uuidv7() PRIMARY KEY,
    user_id UUID NOT NULL REFERENCES "users" (id),
    company_id UUID NOT NULL REFERENCES "companies" (id),
    created_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMPTZ
);

CREATE INDEX idx_user_companies_user_id ON "user_companies" (user_id);

CREATE INDEX idx_user_companies_company_id ON "user_companies" (company_id);

-- +goose StatementEnd
-- +goose Down
-- +goose StatementBegin
DROP INDEX IF EXISTS idx_user_companies_user_id;

DROP INDEX IF EXISTS idx_user_companies_company_id;

DROP TABLE IF EXISTS "user_companies";

DROP TABLE IF EXISTS "users";

DROP TABLE IF EXISTS "companies";

-- +goose StatementEnd