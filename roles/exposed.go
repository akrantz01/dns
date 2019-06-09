package roles

import (
	"github.com/akrantz01/krantz.dev/dns/util"
	bolt "go.etcd.io/bbolt"
	"net/http"
)

// Handle requests regarding roles
func AllRolesHandler(db *bolt.DB) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case "GET":
			list(w, r, db)
			break
		case "POST":
			create(w, r, db)
			break
		default:
			util.Responses.Error(w, http.StatusMethodNotAllowed, "method not allowed")
			return
		}
	}
}

// Handle requests for methods regarding singlar roles
func SingleRoleHandler(path string, db *bolt.DB) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case "GET":
			read(w, r, path, db)
			return
		case "PUT":
			update(w, r, path, db)
			return
		case "DELETE":
			deleteRole(w, r, path, db)
			return
		default:
			util.Responses.Error(w, http.StatusMethodNotAllowed, "method not allowed")
			return
		}
	}
}
