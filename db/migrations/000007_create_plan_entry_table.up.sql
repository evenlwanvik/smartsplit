CREATE TABLE IF NOT EXISTS workout.plan_entry
(
    id         SERIAL PRIMARY KEY,
    workout_id INT         NOT NULL REFERENCES workout.plan(id),
    muscle_id  INT         NOT NULL REFERENCES workout.muscle(id),
    sets       INT         NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    UNIQUE (workout_id, muscle_id)
);