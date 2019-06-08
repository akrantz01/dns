package db

import bolt "go.etcd.io/bbolt"

func Setup(db *bolt.DB) error {
	return db.Update(func(tx *bolt.Tx) error {
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
		return nil
	})
}
