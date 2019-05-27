package routes

import (
	"github.com/akrantz01/krantz.dev/dns/util"
	bolt "go.etcd.io/bbolt"
	"net/http"
)

// Handle GET and POST requests to same route
func RecordsHandler(db *bolt.DB) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case "GET":
			list(w, r, db)
			return
		case "POST":
			create(w, r, db)
			return
		default:
			util.Responses.Error(w, http.StatusMethodNotAllowed, "method not allowed")
			return
		}
	}
}
