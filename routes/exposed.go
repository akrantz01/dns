package routes

import (
	"github.com/akrantz01/krantz.dev/dns/util"
	bolt "go.etcd.io/bbolt"
	"net/http"
)

// Handle requests for methods regarding the entirety of the records
func AllRecordsHandler(db *bolt.DB) func(w http.ResponseWriter, r *http.Request) {
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

// Handle requests for methods regarding singular records
func SingleRecordHandler(path string, db *bolt.DB) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case "GET":
			read(w, r, path, db)
			return
		case "PUT":
			update(w, r, path,  db)
			return
		case "DELETE":
			deleteRecord(w, r, path, db)
			return
		default:
			util.Responses.Error(w, http.StatusMethodNotAllowed, "method not allowed")
			return
		}
	}
}
