package users

import (
	"encoding/json"
	"github.com/akrantz01/krantz.dev/dns/db"
	"github.com/akrantz01/krantz.dev/dns/util"
	"github.com/dgrijalva/jwt-go"
	bolt "go.etcd.io/bbolt"
	"net/http"
)

func update(w http.ResponseWriter, r *http.Request, database *bolt.DB) {
	// Validate initial request with request type, body exist, and content-type
	if r.Method != "PUT" {
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

	// Get username from token
	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		util.Responses.Error(w, http.StatusBadRequest, "invalid JWT claims format")
		return
	}

	// Validate body by decoding json, checking fields exist, and checking field types
	var body map[string]interface{}
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		util.Responses.Error(w, http.StatusBadRequest, "failed to decode body: "+err.Error())
		return
	}
	validationErr, valid := util.ValidateBody(body, []string{"name", "username", "password", "role"}, map[string]map[string]string{
		"name": {"type": "string", "required": "false"},
		"username": {"type": "string", "required": "false"},
		"password": {"type": "string", "required": "false"},
		"role": {"type": "string", "required": "false"},
	})
	if validationErr != "" {
		util.Responses.Error(w, http.StatusBadRequest, validationErr)
		return
	}

	// Get user from database
	u, err := db.UserFromDatabase(claims["sub"].(string), database)
	if err != nil {
		util.Responses.Error(w, http.StatusBadRequest, err.Error())
		return
	}

	// Update values if they exist in body
	if valid["name"] {
		u.Name = body["name"].(string)
	}
	if valid["username"] {
		u.Username = body["username"].(string)
	}
	if valid["password"] {
		u.Password = body["password"].(string)
	}
	if valid["role"] {
		u.Role = body["role"].(string)
	}

	// Write updates to database
	if err := u.Encode(database); err != nil {
		util.Responses.Error(w, http.StatusInternalServerError, "failed to write record to database: "+err.Error())
		return
	}

	util.Responses.Success(w)
}
