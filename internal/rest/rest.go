package rest

import (
	"encoding/json"
	"errors"
	"fmt"
	"log/slog"
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

type ErrorMessage struct {
	Message any `json:"message"`
}

const (
	UnableToDecodeRequestBody  = "unable to decode request body"
	ResourceNotFoundMessage    = "the requested resource was not found"
	InternalServerErrorMessage = "the server encountered a problem and could not process your request"
	BadRequestMessage          = "the request was invalid or cannot be otherwise served"
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
func InternalServerErrorResponse(w http.ResponseWriter, r *http.Request, err error) {
	logError(r, err)
	http.Error(w, InternalServerErrorMessage, http.StatusInternalServerError)
}

// BadRequestResponse sends a 400 Bad Request response with a custom message and logs the error.
func BadRequestResponse(
	w http.ResponseWriter, r *http.Request, message string, err error,
) {
	if message == "" {
		message = BadRequestMessage
	}
	logError(r, err)
	http.Error(w, message, http.StatusBadRequest)
}

// NotFoundResponse sends a 404 Not Found response with a custom message and logs the error.
func NotFoundResponse(
	w http.ResponseWriter, r *http.Request, err error,
) {
	logError(r, err)
	http.Error(w, ResourceNotFoundMessage, http.StatusNotFound)
}

// UnableToGetPathParamFromRequest sends a 400 Bad Request response when a query parameter
// cannot be retrieved from the request, along with a custom error message.
func UnableToGetPathParamFromRequest(
	w http.ResponseWriter, r *http.Request, key string, err error,
) {
	message := "unable to get path parameter: " + key
	BadRequestResponse(w, r, message, err)
}

func LogError(r *http.Request, err error) {
	ctx := r.Context()
	logger := logging.LoggerFromContext(ctx)

	logger.Error(
		"an error occurred",
		"request_method", r.Method,
		"request_url", r.URL.String(),
		"error", err,
	)
}

func writeJSON(
	w http.ResponseWriter,
	status int,
	data any,
	headers http.Header,
) error {
	js, err := json.Marshal(data)
	if err != nil {
		return err
	}

	js = append(js, '\n')

	for key, values := range headers {
		for _, value := range values {
			w.Header().Add(key, value)
		}
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	if _, err = w.Write(js); err != nil {
		return err
	}

	return nil
}

func ErrorResponse(
	w http.ResponseWriter,
	r *http.Request,
	status int,
	message any,
) {
	ctx := r.Context()
	logger := logging.LoggerFromContext(ctx)

	logger.Info("writing error response", slog.Int("status", status), slog.Any("message", message))
	err := writeJSON(w, status, ErrorMessage{Message: message}, nil)
	if err != nil {
		logger.Error("error writing response", "error", err)
		LogError(r, err)
		w.WriteHeader(500)
	}
}

func ServerErrorResponse(w http.ResponseWriter, r *http.Request, err error) {
	logger := logging.LoggerFromContext(r.Context())

	LogError(r, err)
	const serverErrorMsg string = "the server encountered a problem and could not process your request"

	logger.Info(serverErrorMsg)
	ErrorResponse(w, r, http.StatusInternalServerError, serverErrorMsg)
}

func RespondWithJSON(
	w http.ResponseWriter,
	r *http.Request,
	status int,
	data any,
	headers http.Header,
) {
	logger := logging.LoggerFromContext(r.Context())

	logger.Info("marshalling data")
	js, err := json.Marshal(data)
	if err != nil {
		ServerErrorResponse(w, r, err)
	}

	js = append(js, '\n')

	logger.Info("adding headers")
	for key, values := range headers {
		for _, value := range values {
			w.Header().Add(key, value)
		}
	}

	logger.Info("writing response")
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	if _, err = w.Write(js); err != nil {
		ServerErrorResponse(w, r, err)
	}
}

func ReadIntParameter(key string, r *http.Request) (int, error) {
	s := r.PathValue(key)
	if s == "" {
		return 0, fmt.Errorf("empty string parameter")
	}
	i, err := strconv.Atoi(s)
	if err != nil {
		return 0, fmt.Errorf("invalid integer parameter")
	}
	return i, nil
}
