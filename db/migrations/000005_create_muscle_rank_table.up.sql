CREATE TABLE IF NOT EXISTS workout.muscle_rank
(
    id         SERIAL PRIMARY KEY,
    user_id    INT         NOT NULL REFERENCES identity.user (id),
    muscle_id  INT         NOT NULL REFERENCES workout.muscle (id),
    rank       INT         NULL,
    updated_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    UNIQUE (user_id, muscle_id)
);