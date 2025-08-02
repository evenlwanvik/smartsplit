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

// ShowIndex executes the index template
func (svc *WebService) ShowIndex(w http.ResponseWriter) error {
	data := map[string]interface{}{
		// any global values
	}
	if err := svc.tpl.ExecuteTemplate(w, "index.html", data); err != nil {
		return err
	}
	return nil
}
