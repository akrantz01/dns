package users

import (
	"github.com/akrantz01/krantz.dev/dns/db"
	"github.com/akrantz01/krantz.dev/dns/util"
	bolt "go.etcd.io/bbolt"
	"net/http"
)

func read(w http.ResponseWriter, r *http.Request, database *bolt.DB) {
	// Validate initial request with request type
	if r.Method != "GET" {
		util.Responses.Error(w, http.StatusMethodNotAllowed, "method not allowed")
		return
	}

	// TODO: add user authentication checking

	// This is temporary, to be replaced with username from JWT
	if len(r.URL.RawQuery) == 0 {
		util.Responses.Error(w, http.StatusBadRequest, "query parameters are required")
		return
	} else if r.URL.Query().Get("username") == "" {
		util.Responses.Error(w, http.StatusBadRequest, "query paramter 'username' is required")
		return
	}

	// Retrieve user from database
	u, err := db.UserFromDatabase(r.URL.Query().Get("username"), database)
	if err != nil {
		util.Responses.Error(w, http.StatusBadRequest, err.Error())
		return
	}

	// Return user data
	util.Responses.SuccessWithData(w, u)
}
