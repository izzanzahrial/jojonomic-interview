package common

import "net/http"

var successMessage = Envolope{"error": false}

func ReffIDResponse(w http.ResponseWriter, r *http.Request, reffID string) {
	successMessage["reff_id"] = reffID
	SuccessResponse(w, r, http.StatusCreated, successMessage)
}

func SuccessResponse(w http.ResponseWriter, r *http.Request, status int, message any) {
	if err := WriteJSON(w, status, message, nil); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
	}
}
