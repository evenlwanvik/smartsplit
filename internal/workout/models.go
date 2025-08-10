package workout

import (
	"time"
)

type Filters struct {
	UserID *int `json:"user_id,omitempty"`
	PlanID *int `json:"plan,omitempty"`
}

type Muscle struct {
	ID          int    `json:"id"`
	Name        string `json:"name"`
	Group       string `json:"muscle_group"`
	Description string `json:"description,omitempty"`
}

type MuscleInput struct {
	Name        string `json:"name"`
	MuscleGroup string `json:"muscle_group"`
	Description string `json:"description,omitempty"`
}

type MuscleRank struct {
	ID        int       `json:"id"`
	UserID    int       `json:"user_id"`
	MuscleID  int       `json:"muscle_id"`
	Rank      *int      `json:"rank,omitempty"`
	UpdatedAt time.Time `json:"updated_at"`
}

type Plan struct {
	ID      int         `json:"id"`
	UserID  int         `json:"user_id"`
	Date    time.Time   `json:"date"`
	Notes   string      `json:"notes,omitempty"`
	Entries []PlanEntry `json:"entries,omitempty"`
}

type PlanEntry struct {
	ID        int       `json:"id"`
	PlanID    int       `json:"plan_id"`
	MuscleID  int       `json:"muscle_id"`
	Sets      int       `json:"sets"`
	CreatedAt time.Time `json:"created_at"`
}
