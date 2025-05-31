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

// Errors
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

// DecodeJSONFromRequest decodes a JSON request body into the provided struct.
func DecodeJSONFromRequest(r *http.Request, v any) error {
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

func InternalServerError(w http.ResponseWriter) {
	err := WriteJSONResponse(w, http.StatusInternalServerError, ErrorResponse{
		Message: "Internal Server Error",
		Status:  http.StatusInternalServerError,
	})
	if err != nil {
		panic(err)
	}
}

func UnableToGetPathParamFromRequest(w http.ResponseWriter, key string) {
	err := WriteJSONResponse(w, http.StatusBadRequest, ErrorResponse{
		Message: "Unable to get parameter from path: " + key,
		Status:  http.StatusBadRequest,
	})
	if err != nil {
		panic(err)
	}
}

func UnableToDecodeRequestBody(w http.ResponseWriter) {
	err := WriteJSONResponse(w, http.StatusBadRequest, ErrorResponse{
		Message: "Unable to decode request body",
		Status:  http.StatusBadRequest,
	})
	if err != nil {
		panic(err)
	}
}

func BadRequest(w http.ResponseWriter, message string) {
	err := WriteJSONResponse(w, http.StatusBadRequest, ErrorResponse{
		Message: message,
		Status:  http.StatusBadRequest,
	})
	if err != nil {
		panic(err)
	}
}
