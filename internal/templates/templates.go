package templates

import (
	"embed"
	"fmt"
	"html/template"
	"io/fs"
	"path/filepath"
	"strings"
)

type TemplateCache map[string]*template.Template

// NewTemplateCache creates a new template cache from the provided embed.FS
//
// The pattern is used to search for templates in the provided embed.FS.
// Example: "html/*.tmpl"
func NewTemplateCache(files embed.FS, pattern string) (TemplateCache, error) {
	cache := TemplateCache{}

	pages, err := fs.Glob(files, pattern)
	if err != nil {
		return nil, err
	}

	// Define custom functions
	funcMap := template.FuncMap{
		"startsWith": func(s, prefix string) bool {
			return strings.HasPrefix(s, prefix)
		},
	}

	for _, page := range pages {
		name := filepath.Base(page)
		tmpl, err := template.New(name).Funcs(funcMap).ParseFS(files, pattern)
		if err != nil {
			return nil, err
		}
		cache[name] = tmpl
	}
	return cache, nil
}

// GetTemplate returns the template with the provided name from the template cache
func (tc *TemplateCache) GetTemplate(name string) (*template.Template, error) {
	if tc == nil {
		return nil, fmt.Errorf("invalid template cache provided")
	}
	tmpl, ok := (*tc)[name]
	if !ok {
		return nil, fmt.Errorf("template %s not found", name)
	}
	return tmpl, nil
}
