package roles

import (
	"encoding/json"
	"github.com/akrantz01/krantz.dev/dns/db"
	"github.com/akrantz01/krantz.dev/dns/util"
	bolt "go.etcd.io/bbolt"
	"net/http"
)

// Handle the creation of roles
func create(w http.ResponseWriter, r *http.Request, database *bolt.DB) {
	// Validate initial request with request type, body exists, and content type
	if r.Method != "POST" {
		util.Responses.Error(w, http.StatusMethodNotAllowed, "method not allowed")
		return
	} else if r.Body == nil {
		util.Responses.Error(w, http.StatusBadRequest, "body must be present")
		return
	} else if r.Header.Get("Content-Type") != "application/json" {
		util.Responses.Error(w, http.StatusBadRequest, "body must be of type JSON")
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

	// Get u from token
	u, err := db.UserFromToken(token, database)
	if err != nil {
		util.Responses.Error(w, http.StatusInternalServerError, err.Error())
		return
	}

	// Check role
	if u.Role != "admin" {
		util.Responses.Error(w, http.StatusForbidden, "u must be of role 'admin'")
		return
	}

	// Validate body by decoding json, checking fields exist, and checking field type
	var body map[string]interface{}
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		util.Responses.Error(w, http.StatusBadRequest, "failed to decode  body: "+err.Error())
		return
	} else if err, _ := util.ValidateBody(body, []string{"name", "filter", "effect"}, map[string]map[string]string{
		"name": {"type": "string", "required": "true"},
		"description": {"type": "string", "required": "true"},
		"allow": {"type": "string", "required": "true"},
		"deny": {"type": "string", "required": "true"},
	}); err != "" {
		util.Responses.Error(w, http.StatusBadRequest, err)
		return
	}

	// Check if already exists
	role, err := db.GetRole(body["name"].(string), database)
	if err != nil {
		util.Responses.Error(w, http.StatusInternalServerError, "failed to retrieve existing roles")
		return
	} else if role.Name != "" {
		util.Responses.Error(w, http.StatusBadRequest, "role already exists")
		return
	}

	// Write role to database
	if err := db.CreateRole(body["name"].(string), body["description"].(string), body["allow"].(string), body["deny"].(string), database); err != nil {
		util.Responses.Error(w, http.StatusBadRequest, "failed to write role: "+err.Error())
		return
	}

	util.Responses.Success(w)
}
