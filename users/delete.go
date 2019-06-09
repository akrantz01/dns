package users

import (
	"github.com/akrantz01/krantz.dev/dns/db"
	"github.com/akrantz01/krantz.dev/dns/util"
	bolt "go.etcd.io/bbolt"
	"net/http"
	"strings"
)

func deleteUser(w http.ResponseWriter, r *http.Request, database *bolt.DB) {
	// Check request type and headers
	if r.Method != "DELETE" {
		util.Responses.Error(w, http.StatusMethodNotAllowed, "method not allowed")
		return
	} else if r.Header.Get("Authorization") == "" {
		util.Responses.Error(w, http.StatusUnauthorized, "header 'Authorization' is required")
		return
	}

	// Verify JWT in headers
	token, err := db.TokenFromString(r.Header.Get("Authorization"), database)
	if err != nil {
		util.Responses.Error(w, http.StatusUnauthorized, "failed to authenticate: "+err.Error())
		return
	}

	// Get user from database
	u, err := db.UserFromToken(token, database)
	if err != nil {
		util.Responses.Error(w, http.StatusInternalServerError, err.Error())
		return
	}

	// Operate differently if admin
	username := u.Username
	if u.Role == "admin" && r.URL.Query().Get("user") != "" {
		// Allow operating on different user if admin
		username = r.URL.Query().Get("user")
	}

	// Delete user
	if err := database.Update(func(tx *bolt.Tx) error {
		return tx.Bucket([]byte("users")).Delete([]byte(username))
	}); err != nil {
		util.Responses.Error(w, http.StatusInternalServerError, "failed to delete user from database: "+err.Error())
		return
	}

	// Delete all user tokens
	if err := database.Update(func(tx *bolt.Tx) error {
		tokens := tx.Bucket([]byte("tokens"))
		return tokens.ForEach(func(k, v []byte) error {
			if strings.Split(string(k), "-")[0] == username {
				if err := tokens.Delete([]byte(k)); err != nil {
					return err
				}
			}
			return nil
		})
	}); err != nil {
		util.Responses.Error(w, http.StatusInternalServerError, "failed to delete user tokens from database: "+err.Error())
		return
	}

	util.Responses.Success(w)
}
