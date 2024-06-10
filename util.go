package main

import (
	"encoding/json"
	"net/http"
)

func (e *apiError) Error() string {
	return e.Err
}

func WriteJSON(w http.ResponseWriter, status int, data interface{}) error {
	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(status)
	return json.NewEncoder(w).Encode(data)
}
