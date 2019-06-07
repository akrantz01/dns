package db

import (
	"bytes"
	"encoding/binary"
	"encoding/gob"
	bolt "go.etcd.io/bbolt"
)

func (s set) A(name, host string) error {
	return s.Db.Update(func(tx *bolt.Tx) error {
		return tx.Bucket([]byte("A")).Put([]byte(name), []byte(host))
	})
}

func (s set) AAAA(name, host string) error {
	return s.Db.Update(func(tx *bolt.Tx) error {
		return tx.Bucket([]byte("AAAA")).Put([]byte(name), []byte(host))
	})
}

func (s set)CNAME(name, target string) error {
	return s.Db.Update(func(tx *bolt.Tx) error {
		return tx.Bucket([]byte("CNAME")).Put([]byte(name), []byte(target))
	})
}

func (s set) MX(name string, priority uint16, host string) error {
	return s.Db.Update(func(tx *bolt.Tx) error {
		records := tx.Bucket([]byte("MX"))

		// Convert uint16 to binary
		p := make([]byte, binary.MaxVarintLen16)
		binary.BigEndian.PutUint16(p, priority)

		// Write data to bucket
		if err := records.Put([]byte(name + "*priority"), p); err != nil {
			return err
		}
		if err := records.Put([]byte(name + "*host"), []byte(host)); err != nil {
			return err
		}

		return nil
	})
}

func (s set) LOC(name string, version, size, horizontal, vertical uint8, latitude, longitude, altitude uint32) error {
	return s.Db.Update(func(tx *bolt.Tx) error {
		records := tx.Bucket([]byte("LOC"))

		// Convert uint32s to binary
		lat := make([]byte, binary.MaxVarintLen32)
		binary.BigEndian.PutUint32(lat, latitude)
		lon := make([]byte, binary.MaxVarintLen32)
		binary.BigEndian.PutUint32(lon, longitude)
		alt := make([]byte, binary.MaxVarintLen32)
		binary.BigEndian.PutUint32(alt, altitude)

		// Write data to bucket
		if err := records.Put([]byte(name + "*version"), []byte{version}); err != nil {
			return err
		}
		if err := records.Put([]byte(name + "*size"), []byte{size}); err != nil {
			return err
		}
		if err := records.Put([]byte(name + "*horiz"), []byte{horizontal}); err != nil {
			return err
		}
		if err := records.Put([]byte(name + "*vert"), []byte{vertical}); err != nil {
			return err
		}
		if err := records.Put([]byte(name + "*lat"), lat); err != nil {
			return err
		}
		if err := records.Put([]byte(name + "*long"), lon); err != nil {
			return err
		}
		if err := records.Put([]byte(name + "*alt"), alt); err != nil {
			return err
		}

		return nil
	})
}

func (s set) SRV(name string, priority, weight, port uint16, target string) error {
	return s.Db.Update(func(tx *bolt.Tx) error {
		records := tx.Bucket([]byte("SRV"))

		// Convert uint16s to binary
		pri := make([]byte, binary.MaxVarintLen16)
		binary.BigEndian.PutUint16(pri, priority)
		wei := make([]byte, binary.MaxVarintLen16)
		binary.BigEndian.PutUint16(wei, weight)
		por := make([]byte, binary.MaxVarintLen16)
		binary.BigEndian.PutUint16(por, port)

		// Write data to bucket
		if err := records.Put([]byte(name + "*priority"), pri); err != nil {
			return err
		}
		if err := records.Put([]byte(name + "*weight"), wei); err != nil {
			return err
		}
		if err := records.Put([]byte(name + "*port"), por); err != nil {
			return err
		}
		if err := records.Put([]byte(name + "*target"), []byte(target)); err != nil {
			return err
		}

		return nil
	})
}

func (s set) SPF(name string, text []string) error {
	return s.Db.Update(func(tx *bolt.Tx) error {
		// Get proper size of buffer
		size := 0
		for _, t := range text {
			size += len(t)
		}

		// Encode into gob
		encoded := bytes.NewBuffer(make([]byte, 0, size))
		enc := gob.NewEncoder(encoded)
		if err := enc.Encode(text); err != nil {
			return err
		}

		// Write to bucket
		return tx.Bucket([]byte("SPF")).Put([]byte(name), encoded.Bytes())
	})
}

func (s set) TXT(name string, text []string) error {
	return s.Db.Update(func(tx *bolt.Tx) error {
		// Get proper size of buffer
		size := 0
		for _, t := range text {
			size += len(t)
		}

		// Encode into gob
		encoded := bytes.NewBuffer(make([]byte, 0, size))
		enc := gob.NewEncoder(encoded)
		if err := enc.Encode(text); err != nil {
			return err
		}

		// Write to bucket
		return tx.Bucket([]byte("TXT")).Put([]byte(name), encoded.Bytes())
	})
}

func (s set) NS(name, nameserver string) error {
	return s.Db.Update(func(tx *bolt.Tx) error {
		return tx.Bucket([]byte("NS")).Put([]byte(name), []byte(nameserver))
	})
}

func (s set) CAA(name, tag, content string) error {
	return s.Db.Update(func(tx *bolt.Tx) error {
		records := tx.Bucket([]byte("CAA"))

		if err := records.Put([]byte(name + "*tag"), []byte(tag)); err != nil {
			return err
		}
		if err := records.Put([]byte(name + "*content"), []byte(content)); err != nil {
			return err
		}

		return nil
	})
}

func (s set) PTR(name, domain string) error {
	return s.Db.Update(func(tx *bolt.Tx) error {
		return tx.Bucket([]byte("PTR")).Put([]byte(name), []byte(domain))
	})
}

func (s set) CERT(name string, tpe, keytag uint16, algorithm uint8, certificate string) error {
	return s.Db.Update(func(tx *bolt.Tx) error {
		records := tx.Bucket([]byte("CERT"))

		// Convert uint16s to binary
		ty := make([]byte, binary.MaxVarintLen16)
		binary.BigEndian.PutUint16(ty, tpe)
		keta := make([]byte, binary.MaxVarintLen16)
		binary.BigEndian.PutUint16(keta, keytag)

		// Write data to bucket
		if err := records.Put([]byte(name + "*type"), ty); err != nil {
			return err
		}
		if err := records.Put([]byte(name + "*keytag"), keta); err != nil {
			return err
		}
		if err := records.Put([]byte(name + "*algorithm"), []byte{algorithm}); err != nil {
			return err
		}
		if err := records.Put([]byte(name + "*certificate"), []byte(certificate)); err != nil {
			return err
		}

		return nil
	})
}

func (s set) DNSKEY(name string, flags uint16, protocol, algorithm uint8, publickey string) error {
	return s.Db.Update(func(tx *bolt.Tx) error {
		records := tx.Bucket([]byte("DNSKEY"))

		// Convert uint16 to binary
		flgs := make([]byte, binary.MaxVarintLen16)
		binary.BigEndian.PutUint16(flgs, flags)

		// Write data to bucket
		if err := records.Put([]byte(name + "*flags"), flgs); err != nil {
			return err
		}
		if err := records.Put([]byte(name + "*protocol"), []byte{protocol}); err != nil {
			return err
		}
		if err := records.Put([]byte(name + "*algorithm"), []byte{algorithm}); err != nil {
			return err
		}
		if err := records.Put([]byte(name + "*publickey"), []byte(publickey)); err != nil {
			return err
		}

		return nil
	})
}

func (s set) DS(name string, keytag uint16, algorithm, digesttype uint8, digest string) error {
	return s.Db.Update(func(tx *bolt.Tx) error {
		records := tx.Bucket([]byte("DS"))

		// Convert uint16 to binary
		keta := make([]byte, binary.MaxVarintLen16)
		binary.BigEndian.PutUint16(keta, keytag)

		// Write data to bucket
		if err := records.Put([]byte(name + "*keytag"), keta); err != nil {
			return err
		}
		if err := records.Put([]byte(name + "*algorithm"), []byte{algorithm}); err != nil {
			return err
		}
		if err := records.Put([]byte(name + "*digesttype"), []byte{digesttype}); err != nil {
			return err
		}
		if err := records.Put([]byte(name + "*digest"), []byte(digest)); err != nil {
			return err
		}

		return nil
	})
}

func (s set) NAPTR(name string, order, preference uint16, flags, service, regexp, replacement string) error {
	return s.Db.Update(func(tx *bolt.Tx) error {
		records := tx.Bucket([]byte("NAPTR"))

		// Convert uint16s to binary
		ordr := make([]byte, binary.MaxVarintLen16)
		binary.BigEndian.PutUint16(ordr, order)
		pref := make([]byte, binary.MaxVarintLen16)
		binary.BigEndian.PutUint16(pref, preference)

		// Write data to bucket
		if err := records.Put([]byte(name + "*order"), ordr); err != nil {
			return err
		}
		if err := records.Put([]byte(name + "*preference"), pref); err != nil {
			return err
		}
		if err := records.Put([]byte(name + "*flags"), []byte(flags)); err != nil {
			return err
		}
		if err := records.Put([]byte(name + "*service"), []byte(service)); err != nil {
			return err
		}
		if err := records.Put([]byte(name + "*regexp"), []byte(regexp)); err != nil {
			return err
		}
		if err := records.Put([]byte(name + "*replacement"), []byte(replacement)); err != nil {
			return err
		}

		return nil
	})
}

func (s set) SMIMEA(name string, usage, selector, matchingtype uint8, certificate string) error {
	return s.Db.Update(func(tx *bolt.Tx) error {
		records := tx.Bucket([]byte("SMIMEA"))

		// Write data to bucket
		if err := records.Put([]byte(name + "*usage"), []byte{usage}); err != nil {
			return err
		}
		if err := records.Put([]byte(name + "*selector"), []byte{selector}); err != nil {
			return err
		}
		if err := records.Put([]byte(name + "*matching"), []byte{matchingtype}); err != nil {
			return err
		}
		if err := records.Put([]byte(name + "*certificate"), []byte(certificate)); err != nil {
			return err
		}

		return nil
	})
}

func (s set) SSHFP(name string, algorithm, tpe uint8, fingerprint string) error {
	return s.Db.Update(func(tx *bolt.Tx) error {
		records := tx.Bucket([]byte("SSHFP"))

		// Write data to bucket
		if err := records.Put([]byte(name + "*algorithm"), []byte{algorithm}); err != nil {
			return err
		}
		if err := records.Put([]byte(name + "*type"), []byte{tpe}); err != nil {
			return err
		}
		if err := records.Put([]byte(name + "*fingerprint"), []byte(fingerprint)); err != nil {
			return err
		}

		return nil
	})
}

func (s set) TLSA(name string, usage, selector, matchingtype uint8, certificate string) error {
	return s.Db.Update(func(tx *bolt.Tx) error {
		records := tx.Bucket([]byte("TLSA"))

		// Write data to bucket
		if err := records.Put([]byte(name + "*usage"), []byte{usage}); err != nil {
			return err
		}
		if err := records.Put([]byte(name + "*selector"), []byte{selector}); err != nil {
			return err
		}
		if err := records.Put([]byte(name + "*matching"), []byte{matchingtype}); err != nil {
			return err
		}
		if err := records.Put([]byte(name + "*certificate"), []byte(certificate)); err != nil {
			return err
		}

		return nil
	})
}

func (s set) URI(name string, priority, weight uint16, target string) error {
	return s.Db.Update(func(tx *bolt.Tx) error {
		records := tx.Bucket([]byte("URI"))

		// Convert uint16s to binary
		pri := make([]byte, binary.MaxVarintLen16)
		binary.BigEndian.PutUint16(pri, priority)
		wei := make([]byte, binary.MaxVarintLen16)
		binary.BigEndian.PutUint16(wei, weight)

		// Write data to bucket
		if err := records.Put([]byte(name + "*priority"), pri); err != nil {
			return err
		}
		if err := records.Put([]byte(name + "*weight"), wei); err != nil {
			return err
		}
		if err := records.Put([]byte(name + "*target"), []byte(target)); err != nil {
			return err
		}

		return nil
	})
}
