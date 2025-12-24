package utils

import (
	"encoding/json"
	"net/http"
)

// Response represents a standard API response structure
type Response struct {
	Success bool        `json:"success"`
	Data    interface{} `json:"data,omitempty"`
	Error   *ErrorDetail `json:"error,omitempty"`
}

// ErrorDetail represents error information in responses
type ErrorDetail struct {
	Message string `json:"message"`
	Code    string `json:"code,omitempty"`
}

// SendJSON sends a JSON response with the given status code and data
func SendJSON(w http.ResponseWriter, statusCode int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)

	response := Response{
		Success: statusCode >= 200 && statusCode < 300,
		Data:    data,
	}

	json.NewEncoder(w).Encode(response)
}

// SendError sends an error response with the given status code and message
func SendError(w http.ResponseWriter, statusCode int, message string, code ...string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)

	errorDetail := &ErrorDetail{Message: message}
	if len(code) > 0 && code[0] != "" {
		errorDetail.Code = code[0]
	}

	response := Response{
		Success: false,
		Error:   errorDetail,
	}

	json.NewEncoder(w).Encode(response)
}

// SendSuccess sends a success response (200 OK) with data
func SendSuccess(w http.ResponseWriter, data interface{}) {
	SendJSON(w, http.StatusOK, data)
}

// SendCreated sends a created response (201 Created) with data
func SendCreated(w http.ResponseWriter, data interface{}) {
	SendJSON(w, http.StatusCreated, data)
}

// SendNoContent sends a no content response (204 No Content)
func SendNoContent(w http.ResponseWriter) {
	w.WriteHeader(http.StatusNoContent)
}

// Common error helpers

// BadRequest sends a 400 Bad Request error
func BadRequest(w http.ResponseWriter, message string) {
	SendError(w, http.StatusBadRequest, message, "bad_request")
}

// Unauthorized sends a 401 Unauthorized error
func Unauthorized(w http.ResponseWriter, message string) {
	SendError(w, http.StatusUnauthorized, message, "unauthorized")
}

// Forbidden sends a 403 Forbidden error
func Forbidden(w http.ResponseWriter, message string) {
	SendError(w, http.StatusForbidden, message, "forbidden")
}

// NotFound sends a 404 Not Found error
func NotFound(w http.ResponseWriter, message string) {
	SendError(w, http.StatusNotFound, message, "not_found")
}

// Conflict sends a 409 Conflict error
func Conflict(w http.ResponseWriter, message string) {
	SendError(w, http.StatusConflict, message, "conflict")
}

// InternalError sends a 500 Internal Server Error
func InternalError(w http.ResponseWriter, message string) {
	SendError(w, http.StatusInternalServerError, message, "internal_error")
}

// ServiceUnavailable sends a 503 Service Unavailable error
func ServiceUnavailable(w http.ResponseWriter, message string) {
	SendError(w, http.StatusServiceUnavailable, message, "service_unavailable")
}
