CREATE TABLE IF NOT EXISTS workout.plan
(
    id      SERIAL PRIMARY KEY,
    user_id INT  NOT NULL REFERENCES identity.user(id),
    date    DATE NOT NULL,
    notes   TEXT
);