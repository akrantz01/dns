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
		case "POST":
			create(w, r, db)
			break
		default:
			util.Responses.Error(w, http.StatusMethodNotAllowed, "method not allowed")
			return
		}
	}
}
