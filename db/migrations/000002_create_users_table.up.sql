CREATE TABLE identity.user
(
    id            SERIAL PRIMARY KEY,
    email         TEXT        NOT NULL UNIQUE,
    first_name    TEXT        NOT NULL,
    last_name     TEXT        NOT NULL,
    username      TEXT        NOT NULL UNIQUE,
    password_hash TEXT        NOT NULL,
    created_at    TIMESTAMPTZ NOT NULL DEFAULT now(),
    update_at     TIMESTAMPTZ NOT NULL DEFAULT now()
);