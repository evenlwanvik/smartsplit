CREATE TABLE IF NOT EXISTS workout.plans
(
    id         SERIAL PRIMARY KEY,
    user_id    INT         NOT NULL REFERENCES auth.users (id),
    date       DATE        NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    notes      TEXT
);