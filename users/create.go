package users

import (
	"encoding/json"
	"fmt"
	"github.com/akrantz01/krantz.dev/dns/db"
	"github.com/akrantz01/krantz.dev/dns/util"
	bolt "go.etcd.io/bbolt"
	"net/http"
)

func create(w http.ResponseWriter, r *http.Request, database *bolt.DB) {
	// Validate initial request with request, body exists, and content-type
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
	} else if err, _ := util.ValidateBody(body, []string{"name", "username", "password", "role"}, map[string]map[string]string{
		"name": {"type": "string", "required": "true"},
		"username": {"type": "string", "required": "true"},
		"password": {"type": "string", "required": "true"},
		"role": {"type": "string", "required": "true"},
	}); err != "" {
		util.Responses.Error(w, http.StatusBadRequest, err)
		return
	}

	// Check if already exists
	if err := database.View(func(tx *bolt.Tx) error {
		data := tx.Bucket([]byte("users")).Get([]byte(body["username"].(string)))
		if len(data) != 0 {
			return fmt.Errorf("user already exists")
		}
		return nil
	}); err != nil {
		util.Responses.Error(w, http.StatusBadRequest, err.Error())
		return
	}

	// Write to database
	u := db.NewUser(body["name"].(string), body["username"].(string), body["password"].(string), body["role"].(string))
	if err := u.Encode(database); err != nil {
		util.Responses.Error(w, http.StatusInternalServerError, "failed to write to database: "+err.Error())
		return
	}

	util.Responses.Success(w)
}
