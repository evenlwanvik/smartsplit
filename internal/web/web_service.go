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

func (svc *WebService) Index(w http.ResponseWriter) error {
	return svc.tpl.ExecuteTemplate(w, "index.html", nil)
}

func (svc *WebService) Dashboard(w http.ResponseWriter) error {
	return svc.tpl.ExecuteTemplate(w, "dashboard.html", nil)
}

func (svc *WebService) Counters(w http.ResponseWriter) error {
	return svc.tpl.ExecuteTemplate(w, "counters.html", nil)
}

type MuscleInput struct {
	Muscles []*workout.Muscle
}

func (svc *WebService) Muscles(ctx context.Context, w http.ResponseWriter) error {
	muscles, err := svc.workout.GetMuscles(ctx)
	if err != nil {
		return err
	}

	tmplData := MuscleInput{Muscles: muscles}

	return svc.tpl.ExecuteTemplate(w, "muscles.html", tmplData)
}
