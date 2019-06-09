package db

import (
	"github.com/spf13/viper"
	bolt "go.etcd.io/bbolt"
	"gopkg.in/hlandau/passlib.v1"
	"log"
)

func Setup(db *bolt.DB) error {
	// Create buckets for data
	if err := db.Update(func(tx *bolt.Tx) error {
		// Setup records
		if _, err := tx.CreateBucketIfNotExists([]byte("A")); err != nil { return err }
		if _, err := tx.CreateBucketIfNotExists([]byte("AAAA")); err != nil { return err }
		if _, err := tx.CreateBucketIfNotExists([]byte("CNAME")); err != nil { return err }
		if _, err := tx.CreateBucketIfNotExists([]byte("MX")); err != nil { return err }
		if _, err := tx.CreateBucketIfNotExists([]byte("LOC")); err != nil { return err }
		if _, err := tx.CreateBucketIfNotExists([]byte("SRV")); err != nil { return err }
		if _, err := tx.CreateBucketIfNotExists([]byte("SPF")); err != nil { return err }
		if _, err := tx.CreateBucketIfNotExists([]byte("TXT")); err != nil { return err }
		if _, err := tx.CreateBucketIfNotExists([]byte("NS")); err != nil { return err }
		if _, err := tx.CreateBucketIfNotExists([]byte("CAA")); err != nil { return err }
		if _, err := tx.CreateBucketIfNotExists([]byte("PTR")); err != nil { return err }
		if _, err := tx.CreateBucketIfNotExists([]byte("CERT")); err != nil { return err }
		if _, err := tx.CreateBucketIfNotExists([]byte("DNSKEY")); err != nil { return err }
		if _, err := tx.CreateBucketIfNotExists([]byte("DS")); err != nil { return err }
		if _, err := tx.CreateBucketIfNotExists([]byte("NAPTR")); err != nil { return err }
		if _, err := tx.CreateBucketIfNotExists([]byte("SMIMEA")); err != nil { return err }
		if _, err := tx.CreateBucketIfNotExists([]byte("SSHFP")); err != nil { return err }
		if _, err := tx.CreateBucketIfNotExists([]byte("TLSA")); err != nil { return err }
		if _, err := tx.CreateBucketIfNotExists([]byte("URI")); err != nil { return err }

		// Setup authentication
		if _, err := tx.CreateBucketIfNotExists([]byte("users")); err != nil { return err }
		if _, err := tx.CreateBucketIfNotExists([]byte("tokens")); err != nil { return err }
		if _, err := tx.CreateBucketIfNotExists([]byte("roles")); err != nil { return err }
		return nil
	}); err != nil {
		return err
	}

	// Add default user if API is enabled
	if !viper.GetBool("http.disabled") {
		hash, err := passlib.Hash(viper.GetString("http.admin.password"))
		if err != nil {
			log.Fatalf("failed to hash admin password: %v", err)
		}

		u := NewUser(viper.GetString("http.admin.name"), viper.GetString("http.admin.username"), hash, "admin")
		if err := u.Encode(db); err != nil {
			return err
		}
	}

	return nil
}
