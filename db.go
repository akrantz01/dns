package main

import (
	"encoding/binary"
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

func getMXRecord(db *bolt.DB, qname string) (string, uint16, error) {
	var host string
	var priority uint16

	err := db.View(func(tx *bolt.Tx) error {
		records := tx.Bucket([]byte("MX"))

		hostValue := records.Get([]byte(qname[:len(qname) - 1] + "-host"))
		if len(hostValue) != 0 {
			host = string(hostValue)
		}

		priorityValue := records.Get([]byte(qname[:len(qname) - 1] + "-priority"))
		if len(priorityValue) != 0 {
			priority = binary.BigEndian.Uint16(priorityValue)
		}

		return nil
	})
	return host, priority, err
}

func getLOCRecord(db *bolt.DB, qname string) (uint8, uint8, uint8, uint8, uint32, uint32, uint32, error) {
	var (
		version	uint8
		size 	uint8
		horiz 	uint8
		vert 	uint8
		lat 	uint32
		long	uint32
		alt 	uint32
	)

	err := db.View(func(tx *bolt.Tx) error {
		records := tx.Bucket([]byte("LOC"))
		shortenedName := qname[:len(qname) - 1]

		versionValue := records.Get([]byte(shortenedName + "-version"))
		if len(versionValue) != 0 {
			version = versionValue[0]
		}
		sizeValue := records.Get([]byte(shortenedName + "-size"))
		if len(sizeValue) != 0 {
			size = sizeValue[0]
		}
		horizValue := records.Get([]byte(shortenedName + "-horiz"))
		if len(horizValue) != 0 {
			horiz = horizValue[0]
		}
		vertValue := records.Get([]byte(shortenedName + "-vert"))
		if len(vertValue) != 0 {
			vert = vertValue[0]
		}
		latValue := records.Get([]byte(shortenedName + "-lat"))
		if len(latValue) != 0 {
			lat = binary.BigEndian.Uint32(latValue)
		}
		longValue := records.Get([]byte(shortenedName + "-long"))
		if len(longValue) != 0 {
			long = binary.BigEndian.Uint32(longValue)
		}
		altValue := records.Get([]byte(shortenedName + "-alt"))
		if len(altValue) != 0 {
			alt = binary.BigEndian.Uint32(altValue)
		}

		return nil
	})

	return version, size, horiz, vert, lat, long, alt, err
}

func getSRVRecord(db *bolt.DB, qname string) (uint16, uint16, uint16, string, error) {
	var (
		priority uint16
		weight uint16
		port uint16
		target string
	)

	err := db.View(func(tx *bolt.Tx) error {
		records := tx.Bucket([]byte("SRV"))
		shortenedName := qname[:len(qname) - 1]

		priorityValue := records.Get([]byte(shortenedName + "-priority"))
		if len(priorityValue) != 0 {
			priority = binary.BigEndian.Uint16(priorityValue)
		}
		weightValue := records.Get([]byte(shortenedName + "-weight"))
		if len(weightValue) != 0 {
			weight = binary.BigEndian.Uint16(weightValue)
		}
		portValue := records.Get([]byte(shortenedName + "-port"))
		if len(portValue) != 0 {
			port = binary.BigEndian.Uint16(portValue)
		}
		targetValue := records.Get([]byte(shortenedName + "-target"))
		if len(targetValue) != 0 {
			target = string(targetValue)
		}

		return nil
	})

	return priority, weight, port, target, err
}
