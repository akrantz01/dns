package roles

import (
	"github.com/akrantz01/krantz.dev/dns/db"
	"github.com/akrantz01/krantz.dev/dns/util"
	"github.com/dgrijalva/jwt-go"
	bolt "go.etcd.io/bbolt"
	"net/http"
)

func read(w http.ResponseWriter, r *http.Request, path string, database *bolt.DB) {
	// Validate initial request with type and header
	if r.Method != "GET" {
		util.Responses.Error(w, http.StatusMethodNotAllowed, "method not allowed")
		return
	} else if len(r.URL.Path[len(path):]) == 0 {
		util.Responses.Error(w, http.StatusBadRequest, "role must be specified in path")
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

	// Get user from database
	u, err := db.UserFromDatabase(claims["sub"].(string), database)
	if err != nil {
		util.Responses.Error(w, http.StatusUnauthorized, "failed to authenticate: "+err.Error())
		return
	}

	// Check user role
	if u.Role != "admin" {
		util.Responses.Error(w, http.StatusForbidden, "user must be of role 'admin'")
		return
	}

	// Get from database
	role := make(map[string]string)
	allow, deny, err := db.GetRole(r.URL.Path[len(path):], database)
	if err != nil {
		util.Responses.Error(w, http.StatusInternalServerError, "failed to retrieve record")
		return
	}
	role["allow"] = allow
	role["deny"] = deny

	util.Responses.SuccessWithData(w, role)
}
