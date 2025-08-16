package web

import (
	"context"
	"html/template"
	"net/http"

	"github.com/evenlwanvik/smartsplit/internal/workout"
)

type WebClient interface {
	ShowIndex(w http.ResponseWriter) error
}

// WebService handles business logic for web pages.
// TODO: Rename to not use same name as package
type WebService struct {
	tpl     *template.Template
	workout workout.Client
}

// NewWebService creates a new WebService.
func NewWebService(workout workout.Client) *WebService {
	return &WebService{
		tpl:     template.Must(template.ParseFS(htmlFS, "templates/*.html")),
		workout: workout,
	}
}

type DashboardVM struct {
	Suggestion    *SuggestionVM
	Muscles       []*workout.Muscle
	RecentMuscles []MuscleVM
	KPI           KPI
}

type SuggestionVM struct {
	ID           string
	PrimaryLabel string
	Accessories  string
	Avoid        string
}

type KPI struct {
	Sessions      int
	UniqueMuscles int
	RunKM         int
}

type PlanItemVM struct {
	Date    string
	Notes   string
	Entries []PlanEntryVM
}

type PlanEntryVM struct {
	MuscleName string
	Sets       int
}

type MuscleVM struct{ Name string }

func (svc *WebService) Dashboard(ctx context.Context, w http.ResponseWriter) error {
	muscles, err := svc.workout.ReadMuscles(ctx)
	if err != nil {
		return err
	}

	vm := DashboardVM{
		Suggestion:    &SuggestionVM{ID: "seed", PrimaryLabel: "Upper Pull (back, biceps)", Accessories: "Core stability", Avoid: "Chest"},
		RecentMuscles: []MuscleVM{{Name: "Chest"}, {Name: "Quads"}},
		KPI:           KPI{Sessions: 3, UniqueMuscles: 8, RunKM: 18},
		Muscles:       muscles,
	}
	return svc.tpl.ExecuteTemplate(w, "dashboard.html", vm)
}

type MuscleInput struct {
	Muscles []*workout.Muscle
}
