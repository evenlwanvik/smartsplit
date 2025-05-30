package rest

import (
	"encoding/json"
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

func ReadJSONFromRequest(r *http.Request, v any) error {
	if r.Body == nil {
		return errors.New("request body is empty")
	}
	defer r.Body.Close()

	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(v)
	if err != nil {
		return err
	}
	return nil
}

func WriteJSONResponse(w http.ResponseWriter, status int, v any) error {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)

	if v == nil {
		return nil
	}

	encoder := json.NewEncoder(w)
	err := encoder.Encode(v)
	if err != nil {
		return err
	}
	return nil
}

type ErrorResponse struct {
	Message string `json:"message"`
	Status  int    `json:"status"`
}

func InternalServerError(w http.ResponseWriter, err error) {
	err = WriteJSONResponse(w, http.StatusInternalServerError, ErrorResponse{
		Message: err.Error(),
		Status:  http.StatusInternalServerError,
	})
	if err != nil {
		panic(err)
	}
}
