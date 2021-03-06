package utils

import (
	"encoding/json"
	"net/http"
)

func RespondJson(w http.ResponseWriter, r *http.Request, statusCode int, body interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	json.NewEncoder(w).Encode(body)
}
