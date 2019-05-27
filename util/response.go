package util

import (
	"fmt"
	"log"
	"net/http"
)

var Responses = responses{}
type responses struct {}

// Return a generic success
func (r responses) Success(w http.ResponseWriter) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	if _, err := w.Write([]byte(`{"status": "success"}`)); err != nil {
		log.Printf("Failed to write responses: %v", err)
	}
}

// Return error with reason
func (r responses) Error(w http.ResponseWriter, status int, reason string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	if _, err := w.Write([]byte(fmt.Sprintf(`{"status": "error", "reason": "%s"}`, reason))); err != nil {
		log.Printf("Failed to write responses: %v", err)
	}
}
