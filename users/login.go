package users

import (
	"encoding/json"
	"github.com/akrantz01/krantz.dev/dns/db"
	"github.com/akrantz01/krantz.dev/dns/util"
	bolt "go.etcd.io/bbolt"
	"gopkg.in/hlandau/passlib.v1"
	"net/http"
)

func Login(database *bolt.DB) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		// Validate initial request with request type, body exists, and content-type
		if r.Method != "POST" {
			util.Responses.Error(w, http.StatusMethodNotAllowed, "method not allowed")
			return
		} else if r.Body == nil {
			util.Responses.Error(w, http.StatusBadRequest, "body must be present")
			return
		} else if r.Header.Get("Content-Type") != "application/json" {
			util.Responses.Error(w, http.StatusBadRequest, "body must be of type JSON")
			return
		}

		// Validate body by decoding json, checking fields exist, and checking field type
		var body map[string]interface{}
		if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
			util.Responses.Error(w, http.StatusBadRequest, "failed to decode body: "+err.Error())
			return
		} else if err, _ := util.ValidateBody(body, []string{"username", "password"}, map[string]map[string]string{
			"username": {"type": "string", "required": "true"},
			"password": {"type": "string", "required": "true"},
		}); err != "" {
			util.Responses.Error(w, http.StatusBadRequest, err)
			return
		}

		// Check if user exists
		u, err := db.UserFromDatabase(body["username"].(string), database)
		if err != nil {
			util.Responses.Error(w, http.StatusUnauthorized, "invalid username or password")
			return
		}

		// Verify password
		if newHash, err := passlib.Verify(body["password"].(string), u.Password); err != nil {
			util.Responses.Error(w, http.StatusUnauthorized, "invalid username or password")
			return
		} else if newHash != "" {
			u.Password = newHash
			if err := u.Encode(database); err != nil {
				util.Responses.Error(w, http.StatusInternalServerError, "failed to rehash password: "+err.Error())
				return
			}
		}

		// Generate token
		token, err := db.NewToken(u, database)
		if err != nil {
			util.Responses.Error(w, http.StatusInternalServerError, "failed to save token to database: "+err.Error())
			return
		}

		util.Responses.SuccessWithData(w, map[string]string{"token": token})
	}
}
