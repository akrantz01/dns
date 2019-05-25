package main

import (
	"bytes"
	"encoding/binary"
	"encoding/gob"
	bolt "go.etcd.io/bbolt"
	"net"
)

func setupDB(db *bolt.DB) error {
	return db.Update(func(tx *bolt.Tx) error {
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
		shortenedName := qname[:len(qname) - 1]

		hostValue := records.Get([]byte(shortenedName + "-host"))
		if len(hostValue) != 0 {
			host = string(hostValue)
		}

		priorityValue := records.Get([]byte(shortenedName + "-priority"))
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

func getSPFRecord(db *bolt.DB, qname string) ([]string, error) {
	var txt []string

	err := db.View(func(tx *bolt.Tx) error {
		records := tx.Bucket([]byte("SPF"))

		value := records.Get([]byte(qname[:len(qname) - 1]))
		if len(value) != 0 {
			decoded := bytes.NewBuffer(value)
			dec := gob.NewDecoder(decoded)
			if err := dec.Decode(&txt); err != nil {}
		}

		return nil
	})

	return txt, err
}

func getTXTRecord(db *bolt.DB, qname string) ([]string, error) {
	var content []string

	err := db.View(func(tx *bolt.Tx) error {
		records := tx.Bucket([]byte("TXT"))

		value := records.Get([]byte(qname[:len(qname) - 1]))
		if len(value) != 0 {
			decoded := bytes.NewBuffer(value)
			dec := gob.NewDecoder(decoded)
			if err := dec.Decode(&content); err != nil {}
		}

		return nil
	})

	return content, err
}

func getNSRecord(db *bolt.DB, qname string) (string, error) {
	var nameserver string

	err := db.View(func(tx *bolt.Tx) error {
		records := tx.Bucket([]byte("NS"))

		value := records.Get([]byte(qname[:len(qname) - 1]))
		if len(value) != 0 {
			nameserver = string(value)
		}

		return nil
	})

	return nameserver, err
}

func getCAARecord(db *bolt.DB, qname string) (uint8, string, string, error) {
	var (
		tag string
		content string
	)

	err := db.View(func(tx *bolt.Tx) error {
		records := tx.Bucket([]byte("CAA"))
		shortenedName := qname[:len(qname) - 1]

		tagValue := records.Get([]byte(shortenedName + "-tag"))
		if len(tagValue) != 0 {
			tag = string(tagValue)
		}
		contentValue := records.Get([]byte(shortenedName + "-content"))
		if len(contentValue) != 0 {
			content = string(contentValue)
		}

		return nil
	})

	return 0, tag, content, err
}

func getPTRRecord(db *bolt.DB, qname string) (string, error) {
	var ptr string

	err := db.View(func(tx *bolt.Tx) error {
		records := tx.Bucket([]byte("PTR"))

		value := records.Get([]byte(qname[:len(qname) - 1]))
		if len(value) != 0 {
			ptr = string(value)
		}

		return nil
	})

	return ptr, err
}

func getCERTRecord(db *bolt.DB, qname string) (uint16, uint16, uint8, string, error) {
	var (
		tpe uint16
		keyTag uint16
		algo uint8
		cert string
	)

	err := db.View(func(tx *bolt.Tx) error {
		records := tx.Bucket([]byte("CERT"))
		shortenedName := qname[:len(qname) - 1]

		typeValue := records.Get([]byte(shortenedName + "-type"))
		if len(typeValue) != 0 {
			tpe = binary.BigEndian.Uint16(typeValue)
		}
		keyTagValue := records.Get([]byte(shortenedName + "-keytag"))
		if len(keyTagValue) != 0 {
			keyTag = binary.BigEndian.Uint16(keyTagValue)
		}
		algoValue := records.Get([]byte(shortenedName + "-algorithm"))
		if len(algoValue) != 0 {
			algo = algoValue[0]
		}
		certValue := records.Get([]byte(shortenedName + "-certificate"))
		if len(certValue) != 0 {
			cert = string(certValue)
		}

		return nil
	})

	return tpe, keyTag, algo, cert, err
}

func getDNSKEYRecord(db *bolt.DB, qname string) (uint16, uint8, uint8, string, error) {
	var (
		flags uint16
		proto uint8
		algo uint8
		pub string
	)

	err := db.View(func(tx *bolt.Tx) error {
		records := tx.Bucket([]byte("DNSKEY"))
		shortenedName := qname[:len(qname) - 1]

		flagsValue := records.Get([]byte(shortenedName + "-flags"))
		if len(flagsValue) != 0 {
			flags = binary.BigEndian.Uint16(flagsValue)
		}
		protoValue := records.Get([]byte(shortenedName + "-protocol"))
		if len(protoValue) != 0 {
			proto = protoValue[0]
		}
		algoValue := records.Get([]byte(shortenedName + "-algorithm"))
		if len(algoValue) != 0 {
			algo = algoValue[0]
		}
		pubValue := records.Get([]byte(shortenedName + "-publickey"))
		if len(pubValue) != 0 {
			pub = string(pubValue)
		}

		return nil
	})

	return flags, proto, algo, pub, err
}

func getDSRecord(db *bolt.DB, qname string) (uint16, uint8, uint8, string, error) {
	var (
		ktag uint16
		algo uint8
		dtype uint8
		digest string
	)

	err := db.View(func(tx *bolt.Tx) error {
		records := tx.Bucket([]byte("DS"))
		shortenedName := qname[:len(qname) - 1]

		ktagValue := records.Get([]byte(shortenedName + "-keytag"))
		if len(ktagValue) != 0 {
			ktag = binary.BigEndian.Uint16(ktagValue)
		}
		algoValue := records.Get([]byte(shortenedName + "-algorithm"))
		if len(algoValue) != 0 {
			algo = algoValue[0]
		}
		dtypeValue := records.Get([]byte(shortenedName + "-digesttype"))
		if len(dtypeValue) != 0 {
			dtype = dtypeValue[0]
		}
		digestValue := records.Get([]byte(shortenedName + "-digest"))
		if len(digestValue) != 0 {
			digest = string(digestValue)
		}

		return nil
	})

	return ktag, algo, dtype, digest, err
}

func getNAPTRRecord(db *bolt.DB, qname string) (uint16, uint16, string, string, string, string, error) {
	var (
		order uint16
		pref uint16
		flags string
		service string
		regexp string
		replacement string
	)

	err := db.View(func(tx *bolt.Tx) error {
		records := tx.Bucket([]byte("NAPTR"))
		shortenedName := qname[:len(qname) - 1]

		orderValue := records.Get([]byte(shortenedName + "-order"))
		if len(orderValue) != 0 {
			order = binary.BigEndian.Uint16(orderValue)
		}
		prefValue := records.Get([]byte(shortenedName + "-preference"))
		if len(prefValue) != 0 {
			pref = binary.BigEndian.Uint16(prefValue)
		}
		flagsValue := records.Get([]byte(shortenedName + "-flags"))
		if len(flagsValue) != 0 {
			flags = string(flagsValue)
		}
		serviceValue := records.Get([]byte(shortenedName + "-service"))
		if len(serviceValue) != 0 {
			service = string(serviceValue)
		}
		regexpValue := records.Get([]byte(shortenedName + "-regexp"))
		if len(regexpValue) != 0 {
			regexp = string(regexpValue)
		}
		replacementValue := records.Get([]byte(shortenedName + "-replacement"))
		if len(replacementValue) != 0 {
			replacement = string(replacementValue)
		}

		return nil
	})

	return order, pref, flags, service, regexp, replacement, err
}

func getSMIMEARecord(db *bolt.DB, qname string) (uint8, uint8, uint8, string, error) {
	var (
		usage uint8
		selector uint8
		matching uint8
		cert string
	)

	err := db.View(func(tx *bolt.Tx) error {
		records := tx.Bucket([]byte("SMIMEA"))
		shortenedName := qname[:len(qname) - 1]

		usageValue := records.Get([]byte(shortenedName + "-usage"))
		if len(usageValue) != 0 {
			usage = usageValue[0]
		}
		selectorValue := records.Get([]byte(shortenedName + "-selector"))
		if len(selectorValue) != 0 {
			selector = selectorValue[0]
		}
		matchingValue := records.Get([]byte(shortenedName + "-matching"))
		if len(matchingValue) != 0 {
			matching = matchingValue[0]
		}
		certValue := records.Get([]byte(shortenedName + "-certificate"))
		if len(certValue) != 0 {
			cert = string(certValue)
		}

		return nil
	})

	return usage, selector, matching, cert, err
}

func getSSHFPRecord(db *bolt.DB, qname string) (uint8, uint8, string, error) {
	var (
		algo uint8
		tpe uint8
		fingerprint string
	)

	err := db.View(func(tx *bolt.Tx) error {
		records := tx.Bucket([]byte("SSHFP"))
		shortenedName := qname[:len(qname) - 1]

		algoValue := records.Get([]byte(shortenedName + "-algorithm"))
		if len(algoValue) != 0 {
			algo = algoValue[0]
		}
		tpeValue := records.Get([]byte(shortenedName + "-type"))
		if len(tpeValue) != 0 {
			tpe = tpeValue[0]
		}
		fingerprintValue := records.Get([]byte(shortenedName + "-fingerprint"))
		if len(fingerprintValue) != 0 {
			fingerprint = string(fingerprintValue)
		}

		return nil
	})

	return algo, tpe, fingerprint, err
}

func getTLSARecord(db *bolt.DB, qname string) (uint8, uint8, uint8, string, error) {
	var (
		usage uint8
		selector uint8
		matching uint8
		cert string
	)

	err := db.View(func(tx *bolt.Tx) error {
		records := tx.Bucket([]byte("TLSA"))
		shortenedName := qname[:len(qname) - 1]

		usageValue := records.Get([]byte(shortenedName + "-usage"))
		if len(usageValue) != 0 {
			usage = usageValue[0]
		}
		selectorValue := records.Get([]byte(shortenedName + "-selector"))
		if len(selectorValue) != 0 {
			selector = selectorValue[0]
		}
		matchingValue := records.Get([]byte(shortenedName + "-matching"))
		if len(matchingValue) != 0 {
			matching = matchingValue[0]
		}
		certValue := records.Get([]byte(shortenedName + "-certificate"))
		if len(certValue) != 0 {
			cert = string(certValue)
		}

		return nil
	})

	return usage, selector, matching, cert, err
}

func getURIRecord(db *bolt.DB, qname string) (uint16, uint16, string, error) {
	var (
		priority uint16
		weight uint16
		target string
	)

	err := db.View(func(tx *bolt.Tx) error {
		records := tx.Bucket([]byte("URI"))
		shortenedName := qname[:len(qname) - 1]

		priorityValue := records.Get([]byte(shortenedName + "-priority"))
		if len(priorityValue) != 0 {
			priority = binary.BigEndian.Uint16(priorityValue)
		}
		weightValue := records.Get([]byte(shortenedName + "-weight"))
		if len(weightValue) != 0 {
			weight = binary.BigEndian.Uint16(weightValue)
		}
		targetValue := records.Get([]byte(shortenedName + "-target"))
		if len(targetValue) != 0 {
			target = string(targetValue)
		}

		return nil
	})

	return priority, weight, target, err
}
