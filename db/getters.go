package db

import (
	"bytes"
	"encoding/binary"
	"encoding/gob"
	bolt "go.etcd.io/bbolt"
	"log"
	"net"
)

func (g get) A(qname string) net.IP {
	var addr net.IP

	if err := g.Db.View(func(tx *bolt.Tx) error {
		records := tx.Bucket([]byte("A"))

		if value := records.Get([]byte(qname[:len(qname)-1])); len(value) != 0 {
			addr = net.ParseIP(string(value))
		}

		return nil
	}); err != nil {
		log.Printf("Failed to retrieve A record for '%s': %v", qname, err)
	}
	return addr
}

func (g get) AAAA(qname string) net.IP {
	var addr net.IP

	if err := g.Db.View(func(tx *bolt.Tx) error {
		records := tx.Bucket([]byte("AAAA"))

		if value := records.Get([]byte(qname[:len(qname)-1])); len(value) != 0 {
			addr = net.ParseIP(string(value))
		}

		return nil
	}); err != nil {
		log.Printf("Failed to retrieve AAAA record for '%s': %v", qname, err)
	}
	return addr
}

func (g get) CNAME(qname string) string {
	var target string

	if err := g.Db.View(func(tx *bolt.Tx) error {
		records := tx.Bucket([]byte("CNAME"))

		if value := records.Get([]byte(qname[:len(qname)-1])); len(value) != 0 {
			target = string(value)
		}

		return nil
	}); err != nil {
		log.Printf("Failed to retrieve CNAME record for '%s': %v", qname, err)
	}
	return target
}

func (g get) MX(qname string) (string, uint16) {
	var (
		host     string
		priority uint16
	)

	if err := g.Db.View(func(tx *bolt.Tx) error {
		records := tx.Bucket([]byte("MX"))
		shortenedName := qname[:len(qname)-1]

		if hostValue := records.Get([]byte(shortenedName + "*host")); len(hostValue) != 0 {
			host = string(hostValue)
		}
		if priorityValue := records.Get([]byte(shortenedName + "*priority")); len(priorityValue) != 0 {
			priority = binary.BigEndian.Uint16(priorityValue)
		}

		return nil
	}); err != nil {
		log.Printf("Failed to retrieve MX record for '%s': %v", qname, err)
	}
	return host, priority
}

func (g get) LOC(qname string) (uint8, uint8, uint8, uint8, uint32, uint32, uint32) {
	var (
		version uint8
		size    uint8
		horiz   uint8
		vert    uint8
		lat     uint32
		long    uint32
		alt     uint32
	)

	if err := g.Db.View(func(tx *bolt.Tx) error {
		records := tx.Bucket([]byte("LOC"))
		shortenedName := qname[:len(qname)-1]

		if versionValue := records.Get([]byte(shortenedName + "*version")); len(versionValue) != 0 {
			version = versionValue[0]
		}
		if sizeValue := records.Get([]byte(shortenedName + "*size")); len(sizeValue) != 0 {
			size = sizeValue[0]
		}
		if horizValue := records.Get([]byte(shortenedName + "*horiz")); len(horizValue) != 0 {
			horiz = horizValue[0]
		}
		if vertValue := records.Get([]byte(shortenedName + "*vert")); len(vertValue) != 0 {
			vert = vertValue[0]
		}
		if latValue := records.Get([]byte(shortenedName + "*lat")); len(latValue) != 0 {
			lat = binary.BigEndian.Uint32(latValue)
		}
		if longValue := records.Get([]byte(shortenedName + "*long")); len(longValue) != 0 {
			long = binary.BigEndian.Uint32(longValue)
		}
		if altValue := records.Get([]byte(shortenedName + "*alt")); len(altValue) != 0 {
			alt = binary.BigEndian.Uint32(altValue)
		}

		return nil
	}); err != nil {
		log.Printf("Failed to retrieve LOC record for '%s': %v", qname, err)
	}
	return version, size, horiz, vert, lat, long, alt
}

func (g get) SRV(qname string) (uint16, uint16, uint16, string) {
	var (
		priority uint16
		weight   uint16
		port     uint16
		target   string
	)

	if err := g.Db.View(func(tx *bolt.Tx) error {
		records := tx.Bucket([]byte("SRV"))
		shortenedName := qname[:len(qname)-1]

		if priorityValue := records.Get([]byte(shortenedName + "*priority")); len(priorityValue) != 0 {
			priority = binary.BigEndian.Uint16(priorityValue)
		}
		if weightValue := records.Get([]byte(shortenedName + "*weight")); len(weightValue) != 0 {
			weight = binary.BigEndian.Uint16(weightValue)
		}
		if portValue := records.Get([]byte(shortenedName + "*port")); len(portValue) != 0 {
			port = binary.BigEndian.Uint16(portValue)
		}
		if targetValue := records.Get([]byte(shortenedName + "*target")); len(targetValue) != 0 {
			target = string(targetValue)
		}

		return nil
	}); err != nil {
		log.Printf("Failed to retrieve SRV record for '%s': %v", qname, err)
	}
	return priority, weight, port, target
}

func (g get) SPF(qname string) []string {
	var txt []string

	if err := g.Db.View(func(tx *bolt.Tx) error {
		records := tx.Bucket([]byte("SPF"))

		if value := records.Get([]byte(qname[:len(qname)-1])); len(value) != 0 {
			decoded := bytes.NewBuffer(value)
			dec := gob.NewDecoder(decoded)
			if err := dec.Decode(&txt); err != nil {
			}
		}

		return nil
	}); err != nil {
		log.Printf("Failed to retrieve SPF record for '%s': %v", qname, err)
	}
	return txt
}

func (g get) TXT(qname string) []string {
	var content []string

	if err := g.Db.View(func(tx *bolt.Tx) error {
		records := tx.Bucket([]byte("TXT"))

		if value := records.Get([]byte(qname[:len(qname)-1])); len(value) != 0 {
			decoded := bytes.NewBuffer(value)
			dec := gob.NewDecoder(decoded)
			if err := dec.Decode(&content); err != nil {
			}
		}

		return nil
	}); err != nil {
		log.Printf("Failed to retrieve TXT record for '%s': %v", qname, err)
	}
	return content
}

func (g get) NS(qname string) string {
	var nameserver string

	if err := g.Db.View(func(tx *bolt.Tx) error {
		records := tx.Bucket([]byte("NS"))

		if value := records.Get([]byte(qname[:len(qname)-1])); len(value) != 0 {
			nameserver = string(value)
		}

		return nil
	}); err != nil {
		log.Printf("Failed to retrieve NS record for '%s': %v", qname, err)
	}
	return nameserver
}

func (g get) CAA(qname string) (uint8, string, string) {
	var (
		tag     string
		content string
	)

	if err := g.Db.View(func(tx *bolt.Tx) error {
		records := tx.Bucket([]byte("CAA"))
		shortenedName := qname[:len(qname)-1]

		if tagValue := records.Get([]byte(shortenedName + "*tag")); len(tagValue) != 0 {
			tag = string(tagValue)
		}
		if contentValue := records.Get([]byte(shortenedName + "*content")); len(contentValue) != 0 {
			content = string(contentValue)
		}

		return nil
	}); err != nil {
		log.Printf("Failed to retrieve CAA record for '%s': %v", qname, err)
	}
	return 0, tag, content
}

func (g get) PTR(qname string) string {
	var ptr string

	if err := g.Db.View(func(tx *bolt.Tx) error {
		records := tx.Bucket([]byte("PTR"))

		if value := records.Get([]byte(qname[:len(qname)-1])); len(value) != 0 {
			ptr = string(value)
		}

		return nil
	}); err != nil {
		log.Printf("Failed to retrieve PTR record for '%s': %v", qname, err)
	}
	return ptr
}

func (g get) CERT(qname string) (uint16, uint16, uint8, string) {
	var (
		tpe    uint16
		keyTag uint16
		algo   uint8
		cert   string
	)

	if err := g.Db.View(func(tx *bolt.Tx) error {
		records := tx.Bucket([]byte("CERT"))
		shortenedName := qname[:len(qname)-1]

		if typeValue := records.Get([]byte(shortenedName + "*type")); len(typeValue) != 0 {
			tpe = binary.BigEndian.Uint16(typeValue)
		}
		if keyTagValue := records.Get([]byte(shortenedName + "*keytag")); len(keyTagValue) != 0 {
			keyTag = binary.BigEndian.Uint16(keyTagValue)
		}
		if algoValue := records.Get([]byte(shortenedName + "*algorithm")); len(algoValue) != 0 {
			algo = algoValue[0]
		}
		if certValue := records.Get([]byte(shortenedName + "*certificate")); len(certValue) != 0 {
			cert = string(certValue)
		}

		return nil
	}); err != nil {
		log.Printf("Failed to retrieve CERT record for '%s': %v", qname, err)
	}
	return tpe, keyTag, algo, cert
}

func (g get) DNSKEY(qname string) (uint16, uint8, uint8, string) {
	var (
		flags uint16
		proto uint8
		algo  uint8
		pub   string
	)

	if err := g.Db.View(func(tx *bolt.Tx) error {
		records := tx.Bucket([]byte("DNSKEY"))
		shortenedName := qname[:len(qname)-1]

		if flagsValue := records.Get([]byte(shortenedName + "*flags")); len(flagsValue) != 0 {
			flags = binary.BigEndian.Uint16(flagsValue)
		}
		if protoValue := records.Get([]byte(shortenedName + "*protocol")); len(protoValue) != 0 {
			proto = protoValue[0]
		}
		if algoValue := records.Get([]byte(shortenedName + "*algorithm")); len(algoValue) != 0 {
			algo = algoValue[0]
		}
		if pubValue := records.Get([]byte(shortenedName + "*publickey")); len(pubValue) != 0 {
			pub = string(pubValue)
		}

		return nil
	}); err != nil {
		log.Printf("Failed to retrieve DNSKEY record for '%s': %v", qname, err)
	}
	return flags, proto, algo, pub
}

func (g get) DS(qname string) (uint16, uint8, uint8, string) {
	var (
		ktag   uint16
		algo   uint8
		dtype  uint8
		digest string
	)

	if err := g.Db.View(func(tx *bolt.Tx) error {
		records := tx.Bucket([]byte("DS"))
		shortenedName := qname[:len(qname)-1]

		if ktagValue := records.Get([]byte(shortenedName + "*keytag")); len(ktagValue) != 0 {
			ktag = binary.BigEndian.Uint16(ktagValue)
		}
		if algoValue := records.Get([]byte(shortenedName + "*algorithm")); len(algoValue) != 0 {
			algo = algoValue[0]
		}
		if dtypeValue := records.Get([]byte(shortenedName + "*digesttype")); len(dtypeValue) != 0 {
			dtype = dtypeValue[0]
		}
		if digestValue := records.Get([]byte(shortenedName + "*digest")); len(digestValue) != 0 {
			digest = string(digestValue)
		}

		return nil
	}); err != nil {
		log.Printf("Failed to retrieve DS record for '%s': %v", qname, err)
	}
	return ktag, algo, dtype, digest
}

func (g get) NAPTR(qname string) (uint16, uint16, string, string, string, string) {
	var (
		order       uint16
		pref        uint16
		flags       string
		service     string
		regexp      string
		replacement string
	)

	if err := g.Db.View(func(tx *bolt.Tx) error {
		records := tx.Bucket([]byte("NAPTR"))
		shortenedName := qname[:len(qname)-1]

		if orderValue := records.Get([]byte(shortenedName + "*order")); len(orderValue) != 0 {
			order = binary.BigEndian.Uint16(orderValue)
		}
		if prefValue := records.Get([]byte(shortenedName + "*preference")); len(prefValue) != 0 {
			pref = binary.BigEndian.Uint16(prefValue)
		}
		if flagsValue := records.Get([]byte(shortenedName + "*flags")); len(flagsValue) != 0 {
			flags = string(flagsValue)
		}
		if serviceValue := records.Get([]byte(shortenedName + "*service")); len(serviceValue) != 0 {
			service = string(serviceValue)
		}
		if regexpValue := records.Get([]byte(shortenedName + "*regexp")); len(regexpValue) != 0 {
			regexp = string(regexpValue)
		}
		if replacementValue := records.Get([]byte(shortenedName + "*replacement")); len(replacementValue) != 0 {
			replacement = string(replacementValue)
		}

		return nil
	}); err != nil {
		log.Printf("Failed to retrieve NAPTR record for '%s': %v", qname, err)
	}
	return order, pref, flags, service, regexp, replacement
}

func (g get) SMIMEA(qname string) (uint8, uint8, uint8, string) {
	var (
		usage    uint8
		selector uint8
		matching uint8
		cert     string
	)

	if err := g.Db.View(func(tx *bolt.Tx) error {
		records := tx.Bucket([]byte("SMIMEA"))
		shortenedName := qname[:len(qname)-1]

		if usageValue := records.Get([]byte(shortenedName + "*usage")); len(usageValue) != 0 {
			usage = usageValue[0]
		}
		if selectorValue := records.Get([]byte(shortenedName + "*selector")); len(selectorValue) != 0 {
			selector = selectorValue[0]
		}
		if matchingValue := records.Get([]byte(shortenedName + "*matching")); len(matchingValue) != 0 {
			matching = matchingValue[0]
		}
		if certValue := records.Get([]byte(shortenedName + "*certificate")); len(certValue) != 0 {
			cert = string(certValue)
		}

		return nil
	}); err != nil {
		log.Printf("Failed to retrieve SMIMEA record for '%s': %v", qname, err)
	}
	return usage, selector, matching, cert
}

func (g get) SSHFP(qname string) (uint8, uint8, string) {
	var (
		algorithm   uint8
		tpe         uint8
		fingerprint string
	)

	if err := g.Db.View(func(tx *bolt.Tx) error {
		records := tx.Bucket([]byte("SSHFP"))
		shortenedName := qname[:len(qname)-1]

		if algorithmValue := records.Get([]byte(shortenedName + "*algorithm")); len(algorithmValue) != 0 {
			algorithm = algorithmValue[0]
		}
		if typeValue := records.Get([]byte(shortenedName + "*type")); len(typeValue) != 0 {
			tpe = typeValue[0]
		}
		if fingerprintValue := records.Get([]byte(shortenedName + "*fingerprint")); len(fingerprintValue) != 0 {
			fingerprint = string(fingerprintValue)
		}

		return nil
	}); err != nil {
		log.Printf("Failed to retrieve SSHFP record for '%s': %v", qname, err)
	}
	return algorithm, tpe, fingerprint
}

func (g get) TLSA(qname string) (uint8, uint8, uint8, string) {
	var (
		usage       uint8
		selector    uint8
		matching    uint8
		certificate string
	)

	if err := g.Db.View(func(tx *bolt.Tx) error {
		records := tx.Bucket([]byte("TLSA"))
		shortenedName := qname[:len(qname)-1]

		if usageValue := records.Get([]byte(shortenedName + "*usage")); len(usageValue) != 0 {
			usage = usageValue[0]
		}
		if selectorValue := records.Get([]byte(shortenedName + "*selector")); len(selectorValue) != 0 {
			selector = selectorValue[0]
		}
		if matchingValue := records.Get([]byte(shortenedName + "*matching")); len(matchingValue) != 0 {
			matching = matchingValue[0]
		}
		if certificateValue := records.Get([]byte(shortenedName + "*certificate")); len(certificateValue) != 0 {
			certificate = string(certificateValue)
		}

		return nil
	}); err != nil {
		log.Printf("Failed to retrieve TLSA record for '%s': %v", qname, err)
	}
	return usage, selector, matching, certificate
}

func (g get) URI(qname string) (uint16, uint16, string) {
	var (
		priority uint16
		weight   uint16
		target   string
	)

	if err := g.Db.View(func(tx *bolt.Tx) error {
		records := tx.Bucket([]byte("URI"))
		shortenedName := qname[:len(qname)-1]

		if priorityValue := records.Get([]byte(shortenedName + "*priority")); len(priorityValue) != 0 {
			priority = binary.BigEndian.Uint16(priorityValue)
		}
		if weightValue := records.Get([]byte(shortenedName + "*weight")); len(weightValue) != 0 {
			weight = binary.BigEndian.Uint16(weightValue)
		}
		if targetValue := records.Get([]byte(shortenedName + "*target")); len(targetValue) != 0 {
			target = string(targetValue)
		}

		return nil
	}); err != nil {
		log.Printf("Failed to retrieve URI record for '%s': %v", qname, err)
	}
	return priority, weight, target
}
