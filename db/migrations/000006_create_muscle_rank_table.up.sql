CREATE TABLE IF NOT EXISTS workout.muscle_ranks
(
    id         SERIAL PRIMARY KEY,
    user_id    INT         NOT NULL REFERENCES auth.users (id),
    muscle_id  INT         NOT NULL REFERENCES workout.muscles (id),
    rank       INT         NULL,
    updated_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    UNIQUE (user_id, muscle_id)
);