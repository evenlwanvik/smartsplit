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

// ShowIndex executes the index template
func (svc *WebService) ShowIndex(w http.ResponseWriter) error {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	// assuming base.html includes `{{ template "index.html" . }}`
	return svc.tpl.ExecuteTemplate(w, "base.html", nil)
}
