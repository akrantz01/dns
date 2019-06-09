package users

import (
	"github.com/akrantz01/krantz.dev/dns/db"
	"github.com/akrantz01/krantz.dev/dns/util"
	"github.com/dgrijalva/jwt-go"
	bolt "go.etcd.io/bbolt"
	"net/http"
)

func read(w http.ResponseWriter, r *http.Request, database *bolt.DB) {
	// Validate initial request with request type
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

	// Get username from token
	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		util.Responses.Error(w, http.StatusBadRequest, "invalid JWT claims format")
		return
	}

	// Operate differently if admin
	username := claims["sub"].(string)
	u, err := db.UserFromDatabase(username, database)
	if err != nil {
		util.Responses.Error(w, http.StatusUnauthorized, "failed to retrieve user")
		return
	} else if u.Role == "admin" && r.URL.Query().Get("user") != "" {
		// Allow operating on different user if admin
		username = r.URL.Query().Get("user")
	}

	// Get list of all users if admin
	if username == "*" && u.Role == "admin" {
		var users []string

		if err := database.View(func(tx *bolt.Tx) error {
			return tx.Bucket([]byte("users")).ForEach(func(k, v []byte) error {
				users = append(users, string(k))
				return nil
			})
		}); err != nil {
			util.Responses.Error(w, http.StatusInternalServerError, "failed to retrieve all users: "+err.Error())
			return
		}

		util.Responses.SuccessWithData(w, users)
		return
	}

	// Retrieve user from database
	u, err = db.UserFromDatabase(username, database)
	if err != nil {
		util.Responses.Error(w, http.StatusBadRequest, err.Error())
		return
	}

	// Return user data
	util.Responses.SuccessWithData(w, u)
}
