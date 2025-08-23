package web

import (
	"context"
	"html/template"
	"net/http"
	"strconv"

	"github.com/evenlwanvik/smartsplit/internal/logging"
	"github.com/evenlwanvik/smartsplit/internal/rest"
	"github.com/evenlwanvik/smartsplit/internal/workout"
)

type Service struct {
	tpl     *template.Template
	workout workout.Client
}

// NewWebService creates a new WebService.
func NewService(workout workout.Client) Service {
	return Service{
		tpl:     template.Must(template.ParseFS(htmlFS, "templates/*.html")),
		workout: workout,
	}
}

// RegisterRoutes hooks up endpoints.
func (svc *Service) RegisterRoutes(ctx context.Context, mux *http.ServeMux) {
	logger := logging.LoggerFromContext(ctx)

	routeDefinitions := rest.RouteDefinitionList{
		{
			"GET /dashboard",
			svc.dashboardPage,
		},
		{
			"POST /plans/new",
			svc.newPlanPage,
		},
		{
			"DELETE /plans/{id}",
			svc.deletePlan,
		},
		{
			"POST /plans/entries",
			svc.planEntriesPage,
		},
	}

	for _, d := range routeDefinitions {
		logger.Info("adding route", "route", d.Path)
		mux.Handle(d.Path, d.Handler)
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

func (svc *Service) dashboardPage(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	ctx := r.Context()

	muscles, err := svc.workout.ReadMuscles(ctx)
	if err != nil {
		rest.InternalServerErrorResponse(w, r, err)
		return
	}

	vm := DashboardVM{
		Suggestion:    &SuggestionVM{ID: "seed", PrimaryLabel: "Upper Pull (back, biceps)", Accessories: "Core stability", Avoid: "Chest"},
		RecentMuscles: []MuscleVM{{Name: "Chest"}, {Name: "Quads"}},
		KPI:           KPI{Sessions: 3, UniqueMuscles: 8, RunKM: 18},
		Muscles:       muscles,
	}

	if err := svc.tpl.ExecuteTemplate(w, "dashboard.html", vm); err != nil {
		rest.InternalServerErrorResponse(w, r, err)
		return
	}
}

func parseIds(idInput []string) ([]int, error) {
	var parsedIds []int
	for _, id := range idInput {
		parsedId, err := strconv.Atoi(id)
		if err != nil {
			return nil, err
		}
		parsedIds = append(parsedIds, parsedId)
	}
	return parsedIds, nil
}

func (svc *Service) newPlanPage(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	w.Header().Set("Content-Type", "text/html; charset=utf-8")

	if err := r.ParseForm(); err != nil {
		http.Error(w, "bad form", http.StatusBadRequest)
		return
	}

	muscles := r.Form["muscles"]
	if len(muscles) == 0 {
		http.Error(w, "pick at least one muscle", http.StatusBadRequest)
		return
	}
	muscleIds, err := parseIds(muscles)
	if err != nil {
		rest.InternalServerErrorResponse(w, r, err)
		return
	}

	plan, err := svc.workout.CreatePlanWithEntries(ctx, "SomeNote", muscleIds)
	if err != nil {
		rest.InternalServerErrorResponse(w, r, err)
		return
	}

	err = svc.tpl.ExecuteTemplate(w, "_plan_entries_form.html", plan)
	if err != nil {
		rest.InternalServerErrorResponse(w, r, err)
		return
	}
}

func (svc *Service) deletePlan(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	w.Header().Set("Content-Type", "text/html; charset=utf-8")

	id, err := rest.ReadIntParameter("id", r)
	if err != nil {
		rest.BadRequestResponse(w, r, "invalid id", err)
		return
	}

	err = svc.workout.DeletePlan(ctx, id)
	if err != nil {
		rest.InternalServerErrorResponse(w, r, err)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

func (svc *Service) planEntriesPage(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	ctx := r.Context()
	logger := logging.LoggerFromContext(ctx)

	if err := r.ParseForm(); err != nil {
		http.Error(w, "bad form", http.StatusBadRequest)
		return
	}

	entryIDStr := r.PostForm["entry"]
	setsStr := r.PostForm["sets"]

	updateEntries := make([]*workout.PlanEntry, len(entryIDStr))
	for i := 0; i < len(entryIDStr); i++ {
		sets, err := strconv.Atoi(setsStr[i])
		if err != nil {
			rest.InternalServerErrorResponse(w, r, err)
			return
		}
		entryID, err := strconv.Atoi(entryIDStr[i])
		if err != nil {
			rest.InternalServerErrorResponse(w, r, err)
			return
		}
		updateEntries[i], err = svc.workout.UpdatePlanEntrySets(ctx, entryID, sets)
		if err != nil {
			rest.InternalServerErrorResponse(w, r, err)

		}
	}

	logger.Info("this is the form", "muscle", r.PostForm)
}
