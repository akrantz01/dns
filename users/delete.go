package users

import (
	"github.com/akrantz01/krantz.dev/dns/util"
	bolt "go.etcd.io/bbolt"
	"net/http"
)

func deleteUser(w http.ResponseWriter, r *http.Request, database *bolt.DB) {
	// Check request type
	if r.Method != "DELETE" {
		util.Responses.Error(w, http.StatusMethodNotAllowed, "method not allowed")
		return
	}

	// TODO: add user authentication checking

	// This is temporary, to be replaced with username from JWT
	if len(r.URL.RawQuery) == 0 {
		util.Responses.Error(w, http.StatusBadRequest, "query parameters are required")
		return
	} else if r.URL.Query().Get("username") == "" {
		util.Responses.Error(w, http.StatusBadRequest, "query parameter 'username' is required")
	}

	// Delete user
	if err := database.Update(func(tx *bolt.Tx) error {
		return tx.Bucket([]byte("users")).Delete([]byte(r.URL.Query().Get("username")))
	}); err != nil {
		util.Responses.Error(w, http.StatusInternalServerError, "failed to delete user from database: "+err.Error())
		return
	}

	util.Responses.Success(w)
}
