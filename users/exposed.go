package users

import (
	"github.com/akrantz01/krantz.dev/dns/util"
	bolt "go.etcd.io/bbolt"
	"net/http"
)

// Handle requests for methods regarding specific users
func AllUsersHandler(db *bolt.DB) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case "GET":
			read(w, r, db)
			return
		case "POST":
			create(w, r, db)
			return
		case "PUT":
			update(w, r, db)
			return
		case "DELETE":
			deleteUser(w, r, db)
			return
		default:
			util.Responses.Error(w, http.StatusMethodNotAllowed, "method not allowed")
			return
		}
	}
}
