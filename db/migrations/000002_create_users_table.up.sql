CREATE TABLE auth.users
(
    id            SERIAL PRIMARY KEY,
    email         TEXT        NOT NULL UNIQUE,
    first_name    TEXT        NOT NULL,
    last_name     TEXT        NOT NULL,
    username      TEXT        NOT NULL UNIQUE,
    password_hash TEXT        NOT NULL,
    created_at    TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at    TIMESTAMPTZ NOT NULL DEFAULT now()
);

-- Create an default user
INSERT INTO auth.users (email, first_name, last_name, username, password_hash)
VALUES ('a.a@a', 'A', 'A', 'a', '$2y$10$eImiTMZG8MrAwWj7v5b1uO9z3Z5f6k1F4m5Y5Z5F5Z5F5Z5F5Z5F5Z');