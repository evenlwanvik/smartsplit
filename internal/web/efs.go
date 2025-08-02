package web

import (
	"embed"
)

//go:embed templates/*.html
var htmlFS embed.FS

//go:embed static/*
var staticFS embed.FS
