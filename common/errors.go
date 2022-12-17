package common

import (
	"errors"
	"net/http"
)

var (
	errorMessage      = Envolope{"error": true}
	ErrRecordNotFound = errors.New("record not found")
	ErrKafkaNotReady  = errors.New("kafka not ready")
)

func BadRequestResponse(w http.ResponseWriter, r *http.Request, err error) {
	errorMessage["message"] = err.Error()
	ErrorResponse(w, r, http.StatusBadRequest, errorMessage)
}

func ErrorResponse(w http.ResponseWriter, r *http.Request, status int, message any) {
	if err := WriteJSON(w, status, message, nil); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
	}
}

func FailedValidationResponse(w http.ResponseWriter, r *http.Request, errors map[string]string) {
	for key, error := range errors {
		errorMessage[key] = error
	}

	ErrorResponse(w, r, http.StatusUnprocessableEntity, errorMessage)
}

func KafkaErrorResponse(w http.ResponseWriter, r *http.Request, reffid string, err error) {
	errorMessage["reff_id"] = reffid
	errorMessage["message"] = err.Error()
	ErrorResponse(w, r, http.StatusInternalServerError, errorMessage)
}

func GenerateReffIDErrorResponse(w http.ResponseWriter, r *http.Request, err error) {
	errorMessage["message"] = err.Error()
	ErrorResponse(w, r, http.StatusInternalServerError, errorMessage)
}
