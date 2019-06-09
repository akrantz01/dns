package db

import (
	"fmt"
	bolt "go.etcd.io/bbolt"
	"regexp"
)

func CreateRole(name, filter, effect string, db *bolt.DB) error {
	if _, err := regexp.Compile(filter); err != nil {
		return fmt.Errorf("invalid regular expression")
	}

	return db.Update(func(tx *bolt.Tx) error {
		return tx.Bucket([]byte("roles")).Put([]byte(name+"-"+effect), []byte(filter))
	})
}

func GetRole(name string, db *bolt.DB) (string, string, error) {
	var allow string
	var deny string

	if err := db.Update(func(tx *bolt.Tx) error {
		roles := tx.Bucket([]byte("roles"))
		if value := roles.Get([]byte(name+"-allow")); len(value) != 0 {
			allow = string(value)
		}
		if value := roles.Get([]byte(name+"-deny")); len(value) != 0 {
			deny = string(value)
		}
		return nil
	}); err != nil {
		return "", "", err
	}

	return allow, deny, nil
}

func DeleteRole(name, effect string, db *bolt.DB) error {
	return db.Update(func(tx *bolt.Tx) error {
		return tx.Bucket([]byte("roles")).Delete([]byte(name+"-"+effect))
	})
}

func EvaluateRole(name, record string, db *bolt.DB) (bool, error) {
	// Allow by default
	approved := true

	// Retrieve role
	allow, deny, err := GetRole(name, db)
	if err != nil {
		return false, err
	}

	// Evaluate rules if they exist
	if deny != "" {
		matched, err := regexp.Match(deny, []byte(record))
		if err != nil {
			return false, err
		}
		approved = !matched
	}
	if allow != "" {
		matched, err := regexp.Match(allow, []byte(record))
		if err != nil {
			return false, err
		}
		approved = matched
	}

	return approved, nil
}
