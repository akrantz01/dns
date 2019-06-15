package db

import (
	"encoding/json"
	"fmt"
	bolt "go.etcd.io/bbolt"
	"regexp"
)

type Role struct {
	Name        string `json:"name"`
	Description string `json:"description"`
	Allow       string `json:"allow"`
	Deny        string `json:"deny"`
}

func CreateRole(name, description, allowFilter, denyFilter string, db *bolt.DB) error {
	if name == "admin" {
		return fmt.Errorf("cannot add permissions to role 'admin'")
	} else if _, err := regexp.Compile(allowFilter); err != nil {
		return fmt.Errorf("invalid regular expression for allow")
	} else if _, err := regexp.Compile(denyFilter); err != nil {
		return fmt.Errorf("invalid regular expression for deny")
	}

	r := Role{
		Name: name,
		Description: description,
		Allow: allowFilter,
		Deny: denyFilter,
	}
	data, err := json.Marshal(r)
	if err != nil {
		return err
	}

	return db.Update(func(tx *bolt.Tx) error {
		return tx.Bucket([]byte("roles")).Put([]byte(name), data)
	})
}

func GetRole(name string, db *bolt.DB) (*Role, error) {
	var r Role

	if err := db.Update(func(tx *bolt.Tx) error {
		if value := tx.Bucket([]byte("roles")).Get([]byte(name)); len(value) != 0 {
			return json.Unmarshal(value, &r)
		}
		return nil
	}); err != nil {
		return nil, err
	}

	return &r, nil
}

func DeleteRole(name string, db *bolt.DB) error {
	if name == "admin" {
		return fmt.Errorf("cannot delete role 'admin'")
	}

	return db.Update(func(tx *bolt.Tx) error {
		return tx.Bucket([]byte("roles")).Delete([]byte(name))
	})
}

func EvaluateRole(name, record string, db *bolt.DB) (bool, error) {
	// Allow by default
	approved := true

	// Retrieve role
	role, err := GetRole(name, db)
	if err != nil {
		return false, err
	}

	// Evaluate rules if they exist
	if role.Deny != "" {
		matched, err := regexp.Match(role.Deny, []byte(record))
		if err != nil {
			return false, err
		}
		approved = !matched
	}
	if role.Allow != "" {
		matched, err := regexp.Match(role.Allow, []byte(record))
		if err != nil {
			return false, err
		}
		approved = matched
	}

	return approved, nil
}
