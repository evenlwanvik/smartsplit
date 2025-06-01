package rest

import (
	"encoding/json"
	"errors"
	"net/http"
	"strconv"

	"github.com/evenlwanvik/smartsplit/internal/logging"
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

const (
	UnableToDecodeRequestBody = "unable to decode request body"
)

// logError logs an error that occurred while processing a request.
func logError(r *http.Request, err error) {
	ctx := r.Context()
	logger := logging.LoggerFromContext(ctx)

	logger.Error(
		"an error occurred while processing request",
		"method", r.Method,
		"url", r.URL.String(),
		"error", err,
	)
}

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

// GetPathParamInt retrieves a path parameter from the request URL.
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

// WriteJSONResponse writes a JSON response to the http.ResponseWriter with the specified
// status code.
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

// InternalServerError sends a 500 Internal Server Error response with a generic message and
// logs the error.
func InternalServerError(w http.ResponseWriter, r *http.Request, err error) {
	logError(r, err)
	serverErrorMessage := "the server encountered a problem and could not process your request"
	http.Error(w, serverErrorMessage, http.StatusInternalServerError)
}

// BadRequest sends a 400 Bad Request response with a custom message and logs the error.
func BadRequest(
	w http.ResponseWriter, r *http.Request, message string, err error,
) {
	if message == "" {
		message = "Bad request"
	}
	logError(r, err)
	http.Error(w, message, http.StatusBadRequest)
}

// UnableToGetPathParamFromRequest sends a 400 Bad Request response when a query parameter
// cannot be retrieved from the request, along with a custom error message.
func UnableToGetPathParamFromRequest(
	w http.ResponseWriter, r *http.Request, key string, err error,
) {
	message := "unable to get path parameter: " + key
	BadRequest(w, r, message, err)
}
