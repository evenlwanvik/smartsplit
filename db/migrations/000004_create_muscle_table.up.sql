CREATE TABLE IF NOT EXISTS workout.muscle
(
    id           SERIAL PRIMARY KEY,
    name         TEXT NOT NULL,
    muscle_group TEXT NOT NULL,
    description  TEXT
);