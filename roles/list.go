package roles

import (
	"encoding/json"
	"github.com/akrantz01/krantz.dev/dns/db"
	"github.com/akrantz01/krantz.dev/dns/util"
	bolt "go.etcd.io/bbolt"
	"net/http"
)

func list(w http.ResponseWriter, r *http.Request, database *bolt.DB) {
	// Validate initial request with type and headers
	if r.Method != "GET" {
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
	user, err := db.UserFromToken(token, database)
	if err != nil {
		util.Responses.Error(w, http.StatusInternalServerError, err.Error())
		return
	}

	// Check role
	if user.Role != "admin" {
		util.Responses.Error(w, http.StatusForbidden, "user must be of role 'admin'")
		return
	}

	// Get from database
	var roles []db.Role
	if err := database.View(func(tx *bolt.Tx) error {
		return tx.Bucket([]byte("roles")).ForEach(func(k, v []byte) error {
			var role db.Role

			if err := json.Unmarshal(v, &role); err != nil {
				return err
			}

			roles = append(roles, role)

			return nil
		})
	}); err != nil {
		util.Responses.Error(w, http.StatusInternalServerError, "failed to retrieve all records: "+err.Error())
		return
	}

	util.Responses.SuccessWithData(w, roles)
}
