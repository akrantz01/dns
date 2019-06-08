package users

import (
	"github.com/akrantz01/krantz.dev/dns/db"
	"github.com/akrantz01/krantz.dev/dns/util"
	bolt "go.etcd.io/bbolt"
	"net/http"
)

func Logout(database *bolt.DB) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		// Validate initial request with type and authorization
		if r.Method != "GET" {
			util.Responses.Error(w, http.StatusMethodNotAllowed, "method not allowed")
			return
		} else if r.Header.Get("Authorization") == "" {
			util.Responses.Error(w, http.StatusBadRequest, "header 'Authorization' is required")
			return
		}

		// Verify JWT in headers
		token, err := db.TokenFromString(r.Header.Get("Authorization"), database)
		if err != nil {
			util.Responses.Error(w, http.StatusUnauthorized, "failed to authenticate: "+err.Error())
			return
		}

		// Delete token
		if err := database.Update(func(tx *bolt.Tx) error {
			return tx.Bucket([]byte("tokens")).Delete([]byte(token.Header["kid"].(string)))
		}); err != nil {
			util.Responses.Error(w, http.StatusInternalServerError, "failed to delete token: "+err.Error())
			return
		}

		util.Responses.Success(w)
	}
}
