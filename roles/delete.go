package roles

import (
	"github.com/akrantz01/krantz.dev/dns/db"
	"github.com/akrantz01/krantz.dev/dns/util"
	"github.com/dgrijalva/jwt-go"
	bolt "go.etcd.io/bbolt"
	"net/http"
)

func deleteRole(w http.ResponseWriter, r *http.Request, path string, database *bolt.DB) {
	// Validate initial request with type, path, and header
	if r.Method != "DELETE" {
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

	roleName := r.URL.Path[len(path):]

	// Check if query parameter to delete specific effect
	if r.URL.Query().Get("effect") != "" {
		switch r.URL.Query().Get("effect") {
		case "allow":
			if err := db.DeleteRole(roleName, "allow", database); err != nil {
				util.Responses.Error(w, http.StatusInternalServerError, "failed to delete role: "+err.Error())
				return
			}
		case "deny":
			if err := db.DeleteRole(roleName, "deny", database); err != nil {
				util.Responses.Error(w, http.StatusInternalServerError, "failed to delete role: "+err.Error())
				return
			}
		default:
			util.Responses.Error(w, http.StatusBadRequest, "query parameter 'effect' must be 'allow' or 'deny'")
			return
		}

		util.Responses.Success(w)
		return
	}

	// Delete both effects
	if err := db.DeleteRole(roleName, "allow", database); err != nil {
		util.Responses.Error(w, http.StatusInternalServerError, "failed to delete role: "+err.Error())
		return
	}
	if err := db.DeleteRole(roleName, "deny", database); err != nil {
		util.Responses.Error(w, http.StatusInternalServerError, "failed to delete role: "+err.Error())
		return
	}

	util.Responses.Success(w)
}