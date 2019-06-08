package db

import (
	bolt "go.etcd.io/bbolt"
)

func (d delete) A(qname string) error {
	return  d.Db.Update(func(tx *bolt.Tx) error {
		return tx.Bucket([]byte("A")).Delete([]byte(qname))
	})
}

func (d delete) AAAA(qname string) error {
	return d.Db.Update(func(tx *bolt.Tx) error {
		return tx.Bucket([]byte("AAAA")).Delete([]byte(qname))
	})
}

func (d delete) CNAME(qname string) error {
	return d.Db.Update(func(tx *bolt.Tx) error {
		return tx.Bucket([]byte("CNAME")).Delete([]byte(qname))
	})
}

func (d delete) MX(qname string) error {
	return d.Db.Update(func(tx *bolt.Tx) error {
		records := tx.Bucket([]byte("MX"))

		if err := records.Delete([]byte(qname + "*host")); err != nil {
			return err
		}
		return records.Delete([]byte(qname + "*priority"))
	})
}

func (d delete) LOC(qname string) error {
	return d.Db.Update(func(tx *bolt.Tx) error {
		records := tx.Bucket([]byte("LOC"))

		if err := records.Delete([]byte(qname + "*version")); err != nil {
			return err
		}
		if err := records.Delete([]byte(qname + "*size")); err != nil {
			return err
		}
		if err := records.Delete([]byte(qname + "*horiz")); err != nil {
			return err
		}
		if err := records.Delete([]byte(qname + "*vert")); err != nil {
			return err
		}
		if err := records.Delete([]byte(qname + "*lat")); err != nil {
			return err
		}
		if err := records.Delete([]byte(qname + "*long")); err != nil {
			return err
		}
		return records.Delete([]byte(qname + "*alt"))
	})
}

func (d delete) SRV(qname string) error {
	return d.Db.Update(func(tx *bolt.Tx) error {
		records := tx.Bucket([]byte("SRV"))

		if err := records.Delete([]byte(qname + "*priority")); err != nil {
			return err
		}
		if err := records.Delete([]byte(qname + "*weight")); err != nil {
			return err
		}
		if err := records.Delete([]byte(qname + "*port")); err != nil {
			return err
		}
		return records.Delete([]byte(qname + "*target"))
	})
}

func (d delete) SPF(qname string) error {
	return d.Db.Update(func(tx *bolt.Tx) error {
		return tx.Bucket([]byte("SPF")).Delete([]byte(qname))
	})
}

func (d delete) TXT(qname string) error {
	return d.Db.Update(func(tx *bolt.Tx) error {
		return tx.Bucket([]byte("TXT")).Delete([]byte(qname))
	})
}

func (d delete) NS(qname string) error {
	return d.Db.Update(func(tx *bolt.Tx) error {
		return tx.Bucket([]byte("NS")).Delete([]byte(qname))
	})
}

func (d delete) CAA(qname string) error {
	return d.Db.Update(func(tx *bolt.Tx) error {
		records := tx.Bucket([]byte("CAA"))

		if err := records.Delete([]byte(qname + "*tag")); err != nil {
			return err
		}
		return records.Delete([]byte(qname + "*content"))
	})
}

func (d delete) PTR(qname string) error {
	return d.Db.Update(func(tx *bolt.Tx) error {
		return tx.Bucket([]byte("PTR")).Delete([]byte(qname))
	})
}

func (d delete) CERT(qname string) error {
	return d.Db.Update(func(tx *bolt.Tx) error {
		records := tx.Bucket([]byte("CERT"))

		if err := records.Delete([]byte(qname + "*type")); err != nil {
			return err
		}
		if err := records.Delete([]byte(qname + "*keytag")); err != nil {
			return err
		}
		if err := records.Delete([]byte(qname + "*algorithm")); err != nil {
			return err
		}
		return records.Delete([]byte(qname + "*certificate"))
	})
}

func (d delete) DNSKEY(qname string) error {
	return d.Db.Update(func(tx *bolt.Tx) error {
		records := tx.Bucket([]byte("DNSKEY"))

		if err := records.Delete([]byte(qname + "*flags")); err != nil {
			return err
		}
		if err := records.Delete([]byte(qname + "*protocol")); err != nil {
			return err
		}
		if err := records.Delete([]byte(qname + "*algorithm")); err != nil {
			return err
		}
		return records.Delete([]byte(qname + "*publickey"))
	})
}

func (d delete) DS(qname string) error {
	return d.Db.Update(func(tx *bolt.Tx) error {
		records := tx.Bucket([]byte("DS"))

		if err := records.Delete([]byte(qname + "*keytag")); err != nil {
			return err
		}
		if err := records.Delete([]byte(qname + "*algorithm")); err != nil {
			return err
		}
		if err := records.Delete([]byte(qname + "*digesttype")); err != nil {
			return err
		}
		return records.Delete([]byte(qname + "*digest"))
	})
}

func (d delete) NAPTR(qname string) error {
	return d.Db.Update(func(tx *bolt.Tx) error {
		records := tx.Bucket([]byte("NAPTR"))

		if err := records.Delete([]byte(qname + "*order")); err != nil {
			return err
		}
		if err := records.Delete([]byte(qname + "*preference")); err != nil {
			return err
		}
		if err := records.Delete([]byte(qname + "*flags")); err != nil {
			return err
		}
		if err := records.Delete([]byte(qname + "*service")); err != nil {
			return err
		}
		if err := records.Delete([]byte(qname + "*regexp")); err != nil {
			return err
		}
		return records.Delete([]byte(qname + "*replacement"))
	})
}

func (d delete) SMIMEA(qname string) error {
	return d.Db.Update(func(tx *bolt.Tx) error {
		records := tx.Bucket([]byte("SMIMEA"))

		if err := records.Delete([]byte(qname + "*usage")); err != nil {
			return err
		}
		if err := records.Delete([]byte(qname + "*selector")); err != nil {
			return err
		}
		if err := records.Delete([]byte(qname + "*matching")); err != nil {
			return err
		}
		return records.Delete([]byte(qname + "*certificate"))
	})
}

func (d delete) SSHFP(qname string) error {
	return d.Db.Update(func(tx *bolt.Tx) error {
		records := tx.Bucket([]byte("SSHFP"))

		if err := records.Delete([]byte(qname + "*algorithm")); err != nil {
			return err
		}
		if err := records.Delete([]byte(qname + "*type")); err != nil {
			return err
		}
		return records.Delete([]byte(qname + "*fingerprint"))
	})
}

func (d delete) TLSA(qname string) error {
	return d.Db.Update(func(tx *bolt.Tx) error {
		records := tx.Bucket([]byte("TLSA"))

		if err := records.Delete([]byte(qname + "*usage")); err != nil {
			return err
		}
		if err := records.Delete([]byte(qname + "*selector")); err != nil {
			return err
		}
		if err := records.Delete([]byte(qname + "*matching")); err != nil {
			return err
		}
		return records.Delete([]byte(qname + "*certificate"))
	})
}

func (d delete) URI(qname string) error {
	return d.Db.Update(func(tx *bolt.Tx) error {
		records := tx.Bucket([]byte("URI"))

		if err := records.Delete([]byte(qname + "*priority")); err != nil {
			return err
		}
		if err := records.Delete([]byte(qname + "*weight")); err != nil {
			return err
		}
		return records.Delete([]byte(qname + "*target"))
	})
}
