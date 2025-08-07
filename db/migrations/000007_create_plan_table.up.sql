CREATE TABLE IF NOT EXISTS workout.plans
(
    id      SERIAL PRIMARY KEY,
    user_id INT  NOT NULL REFERENCES identity.users(id),
    date    DATE NOT NULL,
    notes   TEXT
);