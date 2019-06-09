package db

import (
	"encoding/binary"
	"encoding/json"
	bolt "go.etcd.io/bbolt"
	"log"
	"net"
)

func (g get) A(qname string) *A {
	a := &A{}

	if err := g.Db.View(func(tx *bolt.Tx) error {
		records := tx.Bucket([]byte("A"))

		if value := records.Get([]byte(qname[:len(qname)-1])); len(value) != 0 {
			a.Address = net.ParseIP(string(value))
		}

		return nil
	}); err != nil {
		log.Printf("Failed to retrieve A record for '%s': %v", qname, err)
        return nil
	} else if len(a.Address) == 0 {
		return nil
	}
	return a
}

func (g get) AAAA(qname string) *AAAA {
	a := &AAAA{}

	if err := g.Db.View(func(tx *bolt.Tx) error {
		records := tx.Bucket([]byte("AAAA"))

		if value := records.Get([]byte(qname[:len(qname)-1])); len(value) != 0 {
			a.Address = net.ParseIP(string(value))
		}

		return nil
	}); err != nil {
		log.Printf("Failed to retrieve AAAA record for '%s': %v", qname, err)
        return nil
	}
	return a
}

func (g get) CNAME(qname string) *CNAME {
	c := &CNAME{}

	if err := g.Db.View(func(tx *bolt.Tx) error {
		records := tx.Bucket([]byte("CNAME"))

		if value := records.Get([]byte(qname[:len(qname)-1])); len(value) != 0 {
			c.Target = string(value)
		}

		return nil
	}); err != nil {
		log.Printf("Failed to retrieve CNAME record for '%s': %v", qname, err)
        return nil
	}
	return c
}

func (g get) MX(qname string) *MX {
	m := &MX{}

	if err := g.Db.View(func(tx *bolt.Tx) error {
		records := tx.Bucket([]byte("MX"))
		shortenedName := qname[:len(qname)-1]

		if hostValue := records.Get([]byte(shortenedName + "*host")); len(hostValue) != 0 {
			m.Host = string(hostValue)
		}
		if priorityValue := records.Get([]byte(shortenedName + "*priority")); len(priorityValue) != 0 {
			m.Priority = binary.BigEndian.Uint16(priorityValue)
		}

		return nil
	}); err != nil {
		log.Printf("Failed to retrieve MX record for '%s': %v", qname, err)
        return nil
	}
	return m
}

func (g get) LOC(qname string) *LOC {
	l := &LOC{}

	if err := g.Db.View(func(tx *bolt.Tx) error {
		records := tx.Bucket([]byte("LOC"))
		shortenedName := qname[:len(qname)-1]

		if versionValue := records.Get([]byte(shortenedName + "*version")); len(versionValue) != 0 {
			l.Version = versionValue[0]
		}
		if sizeValue := records.Get([]byte(shortenedName + "*size")); len(sizeValue) != 0 {
			l.Size = sizeValue[0]
		}
		if horizValue := records.Get([]byte(shortenedName + "*horiz")); len(horizValue) != 0 {
			l.HorizontalPrecision = horizValue[0]
		}
		if vertValue := records.Get([]byte(shortenedName + "*vert")); len(vertValue) != 0 {
			l.VerticalPrecision = vertValue[0]
		}
		if latValue := records.Get([]byte(shortenedName + "*lat")); len(latValue) != 0 {
			l.Latitude = binary.BigEndian.Uint32(latValue)
		}
		if longValue := records.Get([]byte(shortenedName + "*long")); len(longValue) != 0 {
			l.Longitude = binary.BigEndian.Uint32(longValue)
		}
		if altValue := records.Get([]byte(shortenedName + "*alt")); len(altValue) != 0 {
			l.Altitude = binary.BigEndian.Uint32(altValue)
		}

		return nil
	}); err != nil {
		log.Printf("Failed to retrieve LOC record for '%s': %v", qname, err)
        return nil
	}
	return l
}

func (g get) SRV(qname string) *SRV {
	s := &SRV{}

	if err := g.Db.View(func(tx *bolt.Tx) error {
		records := tx.Bucket([]byte("SRV"))
		shortenedName := qname[:len(qname)-1]

		if priorityValue := records.Get([]byte(shortenedName + "*priority")); len(priorityValue) != 0 {
			s.Priority = binary.BigEndian.Uint16(priorityValue)
		}
		if weightValue := records.Get([]byte(shortenedName + "*weight")); len(weightValue) != 0 {
			s.Weight = binary.BigEndian.Uint16(weightValue)
		}
		if portValue := records.Get([]byte(shortenedName + "*port")); len(portValue) != 0 {
			s.Port = binary.BigEndian.Uint16(portValue)
		}
		if targetValue := records.Get([]byte(shortenedName + "*target")); len(targetValue) != 0 {
			s.Target = string(targetValue)
		}

		return nil
	}); err != nil {
		log.Printf("Failed to retrieve SRV record for '%s': %v", qname, err)
        return nil
	}
	return s
}

func (g get) SPF(qname string) *SPF {
	var content []string

	if err := g.Db.View(func(tx *bolt.Tx) error {
		records := tx.Bucket([]byte("SPF"))

		if value := records.Get([]byte(qname[:len(qname)-1])); len(value) != 0 {
			if err := json.Unmarshal(value, &content); err != nil {
				return err
			}
		}

		return nil
	}); err != nil {
		log.Printf("Failed to retrieve SPF record for '%s': %v", qname, err)
        return nil
	}

	// Prune all empty strings
	var text []string
	for _, v := range content {
		if len(v) != 0 {
			text = append(text, v)
		}
	}

	return &SPF{Text: content}
}

func (g get) TXT(qname string) *TXT {
	var content []string

	if err := g.Db.View(func(tx *bolt.Tx) error {
		records := tx.Bucket([]byte("TXT"))

		if value := records.Get([]byte(qname[:len(qname)-1])); len(value) != 0 {
			if err := json.Unmarshal(value, &content); err != nil {
				return err
			}
		}

		return nil
	}); err != nil {
		log.Printf("Failed to retrieve TXT record for '%s': %v", qname, err)
        return nil
	}

	// Prune all empty strings
	var text []string
	for _, v := range content {
		if len(v) != 0 {
			text = append(text, v)
		}
	}

	return &TXT{Text: text}
}

func (g get) NS(qname string) *NS {
	n := &NS{}

	if err := g.Db.View(func(tx *bolt.Tx) error {
		records := tx.Bucket([]byte("NS"))

		if value := records.Get([]byte(qname[:len(qname)-1])); len(value) != 0 {
			n.Nameserver = string(value)
		}

		return nil
	}); err != nil {
		log.Printf("Failed to retrieve NS record for '%s': %v", qname, err)
        return nil
	}
	return n
}

func (g get) CAA(qname string) *CAA {
	c := &CAA{Flag: 0}

	if err := g.Db.View(func(tx *bolt.Tx) error {
		records := tx.Bucket([]byte("CAA"))
		shortenedName := qname[:len(qname)-1]

		if tagValue := records.Get([]byte(shortenedName + "*tag")); len(tagValue) != 0 {
			c.Tag = string(tagValue)
		}
		if contentValue := records.Get([]byte(shortenedName + "*content")); len(contentValue) != 0 {
			c.Content = string(contentValue)
		}

		return nil
	}); err != nil {
		log.Printf("Failed to retrieve CAA record for '%s': %v", qname, err)
        return nil
	}
	return c
}

func (g get) PTR(qname string) *PTR {
	p := &PTR{}

	if err := g.Db.View(func(tx *bolt.Tx) error {
		records := tx.Bucket([]byte("PTR"))

		if value := records.Get([]byte(qname[:len(qname)-1])); len(value) != 0 {
			p.Domain = string(value)
		}

		return nil
	}); err != nil {
		log.Printf("Failed to retrieve PTR record for '%s': %v", qname, err)
        return nil
	}
	return p
}

func (g get) CERT(qname string) *CERT {
	c := &CERT{}

	if err := g.Db.View(func(tx *bolt.Tx) error {
		records := tx.Bucket([]byte("CERT"))
		shortenedName := qname[:len(qname)-1]

		if typeValue := records.Get([]byte(shortenedName + "*type")); len(typeValue) != 0 {
			c.Type = binary.BigEndian.Uint16(typeValue)
		}
		if keyTagValue := records.Get([]byte(shortenedName + "*keytag")); len(keyTagValue) != 0 {
			c.KeyTag = binary.BigEndian.Uint16(keyTagValue)
		}
		if algoValue := records.Get([]byte(shortenedName + "*algorithm")); len(algoValue) != 0 {
			c.Algorithm = algoValue[0]
		}
		if certValue := records.Get([]byte(shortenedName + "*certificate")); len(certValue) != 0 {
			c.Certificate = string(certValue)
		}

		return nil
	}); err != nil {
		log.Printf("Failed to retrieve CERT record for '%s': %v", qname, err)
        return nil
	}
	return c
}

func (g get) DNSKEY(qname string) *DNSKEY {
	d := &DNSKEY{}

	if err := g.Db.View(func(tx *bolt.Tx) error {
		records := tx.Bucket([]byte("DNSKEY"))
		shortenedName := qname[:len(qname)-1]

		if flagsValue := records.Get([]byte(shortenedName + "*flags")); len(flagsValue) != 0 {
			d.Flags = binary.BigEndian.Uint16(flagsValue)
		}
		if protoValue := records.Get([]byte(shortenedName + "*protocol")); len(protoValue) != 0 {
			d.Protocol = protoValue[0]
		}
		if algoValue := records.Get([]byte(shortenedName + "*algorithm")); len(algoValue) != 0 {
			d.Algorithm = algoValue[0]
		}
		if pubValue := records.Get([]byte(shortenedName + "*publickey")); len(pubValue) != 0 {
			d.PublicKey = string(pubValue)
		}

		return nil
	}); err != nil {
		log.Printf("Failed to retrieve DNSKEY record for '%s': %v", qname, err)
        return nil
	}
	return d
}

func (g get) DS(qname string) *DS {
	d := &DS{}

	if err := g.Db.View(func(tx *bolt.Tx) error {
		records := tx.Bucket([]byte("DS"))
		shortenedName := qname[:len(qname)-1]

		if ktagValue := records.Get([]byte(shortenedName + "*keytag")); len(ktagValue) != 0 {
			d.KeyTag = binary.BigEndian.Uint16(ktagValue)
		}
		if algoValue := records.Get([]byte(shortenedName + "*algorithm")); len(algoValue) != 0 {
			d.Algorithm = algoValue[0]
		}
		if dtypeValue := records.Get([]byte(shortenedName + "*digesttype")); len(dtypeValue) != 0 {
			d.DigestType = dtypeValue[0]
		}
		if digestValue := records.Get([]byte(shortenedName + "*digest")); len(digestValue) != 0 {
			d.Digest = string(digestValue)
		}

		return nil
	}); err != nil {
		log.Printf("Failed to retrieve DS record for '%s': %v", qname, err)
        return nil
	}
	return d
}

func (g get) NAPTR(qname string) *NAPTR {
	n := &NAPTR{}

	if err := g.Db.View(func(tx *bolt.Tx) error {
		records := tx.Bucket([]byte("NAPTR"))
		shortenedName := qname[:len(qname)-1]

		if orderValue := records.Get([]byte(shortenedName + "*order")); len(orderValue) != 0 {
			n.Order = binary.BigEndian.Uint16(orderValue)
		}
		if prefValue := records.Get([]byte(shortenedName + "*preference")); len(prefValue) != 0 {
			n.Preference = binary.BigEndian.Uint16(prefValue)
		}
		if flagsValue := records.Get([]byte(shortenedName + "*flags")); len(flagsValue) != 0 {
			n.Flags = string(flagsValue)
		}
		if serviceValue := records.Get([]byte(shortenedName + "*service")); len(serviceValue) != 0 {
			n.Service = string(serviceValue)
		}
		if regexpValue := records.Get([]byte(shortenedName + "*regexp")); len(regexpValue) != 0 {
			n.Regexp = string(regexpValue)
		}
		if replacementValue := records.Get([]byte(shortenedName + "*replacement")); len(replacementValue) != 0 {
			n.Replacement = string(replacementValue)
		}

		return nil
	}); err != nil {
		log.Printf("Failed to retrieve NAPTR record for '%s': %v", qname, err)
        return nil
	}
	return n
}

func (g get) SMIMEA(qname string) *SMIMEA {
	s := &SMIMEA{}

	if err := g.Db.View(func(tx *bolt.Tx) error {
		records := tx.Bucket([]byte("SMIMEA"))
		shortenedName := qname[:len(qname)-1]

		if usageValue := records.Get([]byte(shortenedName + "*usage")); len(usageValue) != 0 {
			s.Usage = usageValue[0]
		}
		if selectorValue := records.Get([]byte(shortenedName + "*selector")); len(selectorValue) != 0 {
			s.Selector = selectorValue[0]
		}
		if matchingValue := records.Get([]byte(shortenedName + "*matching")); len(matchingValue) != 0 {
			s.MatchingType = matchingValue[0]
		}
		if certValue := records.Get([]byte(shortenedName + "*certificate")); len(certValue) != 0 {
			s.Certificate = string(certValue)
		}

		return nil
	}); err != nil {
		log.Printf("Failed to retrieve SMIMEA record for '%s': %v", qname, err)
        return nil
	}
	return s
}

func (g get) SSHFP(qname string) *SSHFP {
	s := &SSHFP{}

	if err := g.Db.View(func(tx *bolt.Tx) error {
		records := tx.Bucket([]byte("SSHFP"))
		shortenedName := qname[:len(qname)-1]

		if algorithmValue := records.Get([]byte(shortenedName + "*algorithm")); len(algorithmValue) != 0 {
			s.Algorithm = algorithmValue[0]
		}
		if typeValue := records.Get([]byte(shortenedName + "*type")); len(typeValue) != 0 {
			s.Type = typeValue[0]
		}
		if fingerprintValue := records.Get([]byte(shortenedName + "*fingerprint")); len(fingerprintValue) != 0 {
			s.Fingerprint = string(fingerprintValue)
		}

		return nil
	}); err != nil {
		log.Printf("Failed to retrieve SSHFP record for '%s': %v", qname, err)
        return nil
	}
	return s
}

func (g get) TLSA(qname string) *TLSA {
	t := &TLSA{}

	if err := g.Db.View(func(tx *bolt.Tx) error {
		records := tx.Bucket([]byte("TLSA"))
		shortenedName := qname[:len(qname)-1]

		if usageValue := records.Get([]byte(shortenedName + "*usage")); len(usageValue) != 0 {
			t.Usage = usageValue[0]
		}
		if selectorValue := records.Get([]byte(shortenedName + "*selector")); len(selectorValue) != 0 {
			t.Selector = selectorValue[0]
		}
		if matchingValue := records.Get([]byte(shortenedName + "*matching")); len(matchingValue) != 0 {
			t.MatchingType = matchingValue[0]
		}
		if certificateValue := records.Get([]byte(shortenedName + "*certificate")); len(certificateValue) != 0 {
			t.Certificate = string(certificateValue)
		}

		return nil
	}); err != nil {
		log.Printf("Failed to retrieve TLSA record for '%s': %v", qname, err)
        return nil
	}
	return t
}

func (g get) URI(qname string) *URI {
	u := &URI{}

	if err := g.Db.View(func(tx *bolt.Tx) error {
		records := tx.Bucket([]byte("URI"))
		shortenedName := qname[:len(qname)-1]

		if priorityValue := records.Get([]byte(shortenedName + "*priority")); len(priorityValue) != 0 {
			u.Priority = binary.BigEndian.Uint16(priorityValue)
		}
		if weightValue := records.Get([]byte(shortenedName + "*weight")); len(weightValue) != 0 {
			u.Weight = binary.BigEndian.Uint16(weightValue)
		}
		if targetValue := records.Get([]byte(shortenedName + "*target")); len(targetValue) != 0 {
			u.Target = string(targetValue)
		}

		return nil
	}); err != nil {
		log.Printf("Failed to retrieve URI record for '%s': %v", qname, err)
        return nil
	}
	return u
}
