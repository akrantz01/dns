package main

import (
	bolt "go.etcd.io/bbolt"
	"net"
)

func setupDB(db *bolt.DB) error {
	return db.Batch(func(tx *bolt.Tx) error {
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
		return nil
	})
}

func getARecord(db *bolt.DB, qname string) (net.IP, error) {
	var addr net.IP

	err := db.View(func(tx *bolt.Tx) error {
		records := tx.Bucket([]byte("A"))

		value := records.Get([]byte(qname[:len(qname) - 1]))
		if len(value) != 0 {
			addr = net.ParseIP(string(value))
		}

		return nil
	})
	return addr, err
}

func getAAAARecord(db *bolt.DB, qname string) (net.IP, error) {
	var addr net.IP

	err := db.View(func(tx *bolt.Tx) error {
		records := tx.Bucket([]byte("AAAA"))

		value := records.Get([]byte(qname[:len(qname) - 1]))
		if len(value) != 0 {
			addr = net.ParseIP(string(value))
		}

		return nil
	})
	return addr, err
}

func getCNAMERecord(db *bolt.DB, qname string) (string, error) {
	var target string

	err := db.View(func(tx *bolt.Tx) error {
		records := tx.Bucket([]byte("CNAME"))

		value := records.Get([]byte(qname[:len(qname) - 1]))
		if len(value) != 0 {
			target = string(value)
		}

		return nil
	})
	return target, err
}
