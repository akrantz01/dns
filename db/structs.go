package db

import bolt "go.etcd.io/bbolt"

var (
	// Getter object for "static" methods
	Get = get{Db: nil}
	// Setter object for "static" methods
	Set = set{Db: nil}
	// Delete object for "static" methods
	Delete = deleteRecord{Db: nil}
)

// Getters for different record types
type get struct {
	Db *bolt.DB
}

// Setters for different record types
type set struct {
	Db *bolt.DB
}

// Delete different record types
type deleteRecord struct {
	Db *bolt.DB
}
