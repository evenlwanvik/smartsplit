package rest

import (
	"errors"
	"net/http"
	"strconv"
)

type RouteDefinition struct {
	Path    string
	Handler http.HandlerFunc
}

type RouteDefinitionList []RouteDefinition

var (
	// ErrQueryParamNotFound If query parameter key is not found
	ErrQueryParamNotFound = errors.New("query parameter not found")
	// ErrInvalidQueryParam If query parameter value is invalid
	ErrInvalidQueryParam = errors.New("invalid query parameter value")

	// ErrPathParamNotFound If path parameter key is not found
	ErrPathParamNotFound = errors.New("path parameter not found")
	// ErrInvalidPathParam If path parameter value is invalid
	ErrInvalidPathParam = errors.New("invalid path parameter value")
)

// GetQueryParamInt retrieves a query parameter from the request URL.
func GetQueryParamInt(r *http.Request, key string) (int, error) {
	s := r.URL.Query().Get(key)
	if s == "" {
		return 0, ErrQueryParamNotFound
	}
	i, err := strconv.Atoi(s)
	if err != nil {
		return 0, ErrInvalidQueryParam
	}
	return i, nil
}

func GetPathParamInt(r *http.Request, key string) (int, error) {
	s := r.PathValue(key)
	if s == "" {
		return 0, ErrPathParamNotFound
	}
	i, err := strconv.Atoi(s)
	if err != nil {
		return 0, ErrInvalidPathParam
	}
	return i, nil
}
