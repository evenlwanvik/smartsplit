CREATE TABLE IF NOT EXISTS workout.plan_entries
(
    id         SERIAL PRIMARY KEY,
    plan_id INT         NOT NULL REFERENCES workout.plans(id),
    muscle_id  INT         NOT NULL REFERENCES workout.muscles(id),
    sets       INT         NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    UNIQUE (plan_id, muscle_id)
);