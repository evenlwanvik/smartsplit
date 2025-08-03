package web

import (
	"html/template"
	"net/http"
)

type WebClient interface {
	ShowIndex(w http.ResponseWriter) error
}

// WebService handles business logic for web pages.
// TODO: Rename to not use same name as package
type WebService struct {
	tpl *template.Template
}

// NewWebService creates a new WebService.
func NewWebService() *WebService {
	return &WebService{tpl: template.Must(
		template.ParseFS(htmlFS, "templates/*.html"),
	)}
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
