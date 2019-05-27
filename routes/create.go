package routes

import (
	"bytes"
	"encoding/binary"
	"encoding/gob"
	"encoding/json"
	"github.com/akrantz01/krantz.dev/dns/util"
	bolt "go.etcd.io/bbolt"
	"net"
	"net/http"
	"strings"
)

// Handle the creation of records
func create(w http.ResponseWriter, r *http.Request, db *bolt.DB) {
	// Validate initial request with request type, body exists, and content type
	if r.Method != "POST" {
		util.Responses.Error(w, http.StatusMethodNotAllowed, "method not allowed")
		return
	} else if r.Body == nil {
		util.Responses.Error(w, http.StatusBadRequest, "body must be present")
		return
	} else if r.Header.Get("Content-Type") != "application/json" {
		util.Responses.Error(w, http.StatusBadRequest, "body must be of type JSON")
		return
	}

	// Validate body by decoding json, checking fields exists, and checking field type
	var body map[string]interface{}
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		util.Responses.Error(w, http.StatusBadRequest, "failed to decode body: " + err.Error())
		return
	} else if !util.Exists(body, "type") {
		util.Responses.Error(w, http.StatusBadRequest, "field 'type' is required")
		return
	} else if !util.Exists(body, "name") {
		util.Responses.Error(w, http.StatusBadRequest, "field 'name' is required")
		return
	} else if !util.Types.String(body["type"]) {
		util.Responses.Error(w, http.StatusBadRequest, "field 'type' must be a string")
		return
	} else if !util.Types.String(body["name"]) {
		util.Responses.Error(w, http.StatusBadRequest, "field 'name' must be a string")
		return
	}

	// Parse out body by type
	switch strings.ToUpper(body["type"].(string)) {
	case "A":
		if !util.Exists(body, "host") {
			util.Responses.Error(w, http.StatusBadRequest, "field 'host' is required")
			return
		} else if !util.Types.String(body["host"]) {
			util.Responses.Error(w, http.StatusBadRequest, "field 'host' must be a string")
			return
		} else if ip := net.ParseIP(body["host"].(string)); ip.To4().String() == "<nil>" {
			util.Responses.Error(w, http.StatusBadRequest, "field 'host' must be an IPv4 address")
			return
		} else if err := db.Update(func(tx *bolt.Tx) error {
			records := tx.Bucket([]byte("A"))
			return records.Put([]byte(body["name"].(string)), []byte(body["host"].(string)))
		}); err != nil {
			util.Responses.Error(w, http.StatusInternalServerError, "failed to write record to database: " + err.Error())
			return
		}
	case "AAAA":
		if !util.Exists(body, "host") {
			util.Responses.Error(w, http.StatusBadRequest, "field 'host' is required")
			return
		} else if !util.Types.String(body["host"]) {
			util.Responses.Error(w, http.StatusBadRequest, "field 'host' must be a string")
			return
		} else if ip := net.ParseIP(body["host"].(string)); ip.To4().String() != "<nil>" {
			util.Responses.Error(w, http.StatusBadRequest, "field 'host' must be an IPv6 address")
			return
		} else if err := db.Update(func(tx *bolt.Tx) error {
			records := tx.Bucket([]byte("AAAA"))
			return records.Put([]byte(body["name"].(string)), []byte(body["host"].(string)))
		}); err != nil {
			util.Responses.Error(w, http.StatusInternalServerError, "failed to write record to database: " + err.Error())
			return
		}
	case "CNAME":
		if !util.Exists(body, "target") {
			util.Responses.Error(w, http.StatusBadRequest, "field 'target' is required")
			return
		} else if !util.Types.String(body["target"]) {
			util.Responses.Error(w, http.StatusBadRequest, "field 'target' must be a string")
			return
		} else if err := db.Update(func(tx *bolt.Tx) error {
			records := tx.Bucket([]byte("CNAME"))
			return records.Put([]byte(body["name"].(string)), []byte(body["target"].(string)))
		}); err != nil {
			util.Responses.Error(w, http.StatusInternalServerError, "failed to write record to database: " + err.Error())
			return
		}
	case "MX":
		if !util.Exists(body, "host") {
			util.Responses.Error(w, http.StatusBadRequest, "field 'host' is required")
			return
		} else if !util.Exists(body, "priority") {
			util.Responses.Error(w, http.StatusBadRequest, "field 'priority' is required")
			return
		} else if !util.Types.String(body["host"]) {
			util.Responses.Error(w, http.StatusBadRequest, "field 'host' must be a string")
			return
		} else if !util.Types.Uint16(body["priority"]) {
			util.Responses.Error(w, http.StatusBadRequest, "field 'priority' must be an integer between 0 and 65535")
			return
		} else if err := db.Update(func(tx *bolt.Tx) error {
			records := tx.Bucket([]byte("MX"))
			// Convert uint16 to binary
			priority := make([]byte, binary.MaxVarintLen16)
			binary.BigEndian.PutUint16(priority, uint16(body["priority"].(float64)))
			// Write data to bucket
			if err := records.Put([]byte(body["name"].(string) + "*priority"), priority); err != nil { return err }
			if err := records.Put([]byte(body["name"].(string) + "*host"), []byte(body["host"].(string))); err != nil { return err }
			return nil
		}); err != nil {
			util.Responses.Error(w, http.StatusBadRequest, "failed to write record to database: " + err.Error())
			return
		}
	case "LOC":
		if !util.Exists(body, "version") {
			util.Responses.Error(w, http.StatusBadRequest, "field 'version' is required")
			return
		} else if !util.Exists(body, "size") {
			util.Responses.Error(w, http.StatusBadRequest, "field 'size' is required")
			return
		} else if !util.Exists(body, "horizontal-precision") {
			util.Responses.Error(w, http.StatusBadRequest, "field 'horizontal-precision' is required")
			return
		} else if !util.Exists(body, "vertical-precision") {
			util.Responses.Error(w, http.StatusBadRequest, "field 'vertical-precision' is required")
			return
		} else if !util.Exists(body, "latitude") {
			util.Responses.Error(w, http.StatusBadRequest, "field 'latitude' is required")
			return
		} else if !util.Exists(body, "longitude") {
			util.Responses.Error(w, http.StatusBadRequest, "field 'longitude' is required")
			return
		} else if !util.Exists(body, "altitude") {
			util.Responses.Error(w, http.StatusBadRequest, "field 'altitude' is required")
			return
		} else if !util.Types.Uint8(body["version"]) {
			util.Responses.Error(w, http.StatusBadRequest, "field 'version' must be an integer between 0 and 255")
			return
		} else if !util.Types.Uint8(body["size"]) {
			util.Responses.Error(w, http.StatusBadRequest, "field 'size' must be an integer between 0 and 255")
			return
		} else if !util.Types.Uint8(body["horizontal-precision"]) {
			util.Responses.Error(w, http.StatusBadRequest, "field 'horizontal-precision' must be an integer between 0 and 255")
			return
		} else if !util.Types.Uint8(body["vertical-precision"]) {
			util.Responses.Error(w, http.StatusBadRequest, "field 'vertical-precision' must be an integer between 0 and 255")
			return
		} else if !util.Types.Uint32(body["latitude"]) {
			util.Responses.Error(w, http.StatusBadRequest, "field 'latitude' must be an integer between 0 and 4294967295")
			return
		} else if !util.Types.Uint32(body["longitude"]) {
			util.Responses.Error(w, http.StatusBadRequest, "field 'longitude' must be an integer between 0 and 4294967295")
			return
		} else if !util.Types.Uint32(body["altitude"]) {
			util.Responses.Error(w, http.StatusBadRequest, "field 'altitude' must be an integer between 0 and 4294967295")
		} else if err := db.Update(func(tx *bolt.Tx) error {
			records := tx.Bucket([]byte("LOC"))
			// Convert uint32s to binary
			latitude := make([]byte, binary.MaxVarintLen32)
			longitude := make([]byte, binary.MaxVarintLen32)
			altitude := make([]byte, binary.MaxVarintLen32)
			binary.BigEndian.PutUint32(latitude, uint32(body["latitude"].(float64)))
			binary.BigEndian.PutUint32(longitude, uint32(body["latitude"].(float64)))
			binary.BigEndian.PutUint32(altitude, uint32(body["altitude"].(float64)))
			// Write data to bucket
			if err := records.Put([]byte(body["name"].(string) + "*version"), []byte{uint8(body["version"].(float64))}); err != nil { return err }
			if err := records.Put([]byte(body["name"].(string) + "*size"), []byte{uint8(body["size"].(float64))}); err != nil { return err }
			if err := records.Put([]byte(body["name"].(string) + "*horiz"), []byte{uint8(body["horizontal-precision"].(float64))}); err != nil { return  err }
			if err := records.Put([]byte(body["name"].(string) + "*vert"), []byte{uint8(body["vertical-precision"].(float64))}); err != nil { return err }
			if err := records.Put([]byte(body["name"].(string) + "*lat"), latitude); err != nil { return err }
			if err := records.Put([]byte(body["name"].(string) + "*long"), longitude); err != nil { return err }
			if err := records.Put([]byte(body["name"].(string) + "*alt"), altitude); err != nil { return err }
			return nil
		}); err != nil {
			util.Responses.Error(w, http.StatusInternalServerError, "failed to write record to database: " + err.Error())
			return
		}
	case "SRV":
		if !util.Exists(body, "priority") {
			util.Responses.Error(w, http.StatusBadRequest, "field 'priority' is required")
			return
		} else if !util.Exists(body, "weight") {
			util.Responses.Error(w, http.StatusBadRequest, "field 'weight' is required")
			return
		} else if !util.Exists(body, "port") {
			util.Responses.Error(w, http.StatusBadRequest, "field 'port' is required")
			return
		} else if !util.Exists(body, "target") {
			util.Responses.Error(w, http.StatusBadRequest, "field 'target' is required")
			return
		} else if !util.Types.Uint16(body["priority"]) {
			util.Responses.Error(w, http.StatusBadRequest, "field 'priority' must be an integer between 0 and 65535")
			return
		} else if !util.Types.Uint16(body["weight"]) {
			util.Responses.Error(w, http.StatusBadRequest, "field 'weight' must be an integer between 0 and 65535")
			return
		} else if !util.Types.Uint16(body["port"]) {
			util.Responses.Error(w, http.StatusBadRequest, "field 'port' must be an integer between 0 and 65535")
			return
		} else if !util.Types.String(body["target"]) {
			util.Responses.Error(w, http.StatusBadRequest, "field 'target' must be a string")
			return
		} else if err := db.Update(func(tx *bolt.Tx) error {
			records := tx.Bucket([]byte("SRV"))
			// Convert uint16s to binary
			priority := make([]byte, binary.MaxVarintLen16)
			weight := make([]byte, binary.MaxVarintLen16)
			port := make([]byte, binary.MaxVarintLen16)
			binary.BigEndian.PutUint16(priority, uint16(body["priority"].(float64)))
			binary.BigEndian.PutUint16(weight, uint16(body["weight"].(float64)))
			binary.BigEndian.PutUint16(port, uint16(body["port"].(float64)))
			// Write data to bucket
			if err := records.Put([]byte(body["name"].(string) + "*priority"), priority); err != nil { return err }
			if err := records.Put([]byte(body["name"].(string) + "*weight"), weight); err != nil { return err }
			if err := records.Put([]byte(body["name"].(string) + "*port"), port); err != nil { return err }
			if err := records.Put([]byte(body["name"].(string) + "*target"), []byte(body["target"].(string))); err != nil { return err}
			return nil
		}); err != nil {
			util.Responses.Error(w, http.StatusInternalServerError, "failed to write record to database: " + err.Error())
			return
		}
	case "SPF":
		if !util.Exists(body, "text") {
			util.Responses.Error(w, http.StatusBadRequest, "field 'text' is required")
			return
		} else if !util.Types.StringArray(body["text"]) {
			util.Responses.Error(w, http.StatusBadRequest, "field 'text' must be an array of strings")
			return
		}
		text, _ := util.ConvertArrayToString(body["text"].([]interface{}))
		if len(text) < 1 {
			util.Responses.Error(w, http.StatusBadRequest, "field 'text' must a length of length 1")
			return
		} else if err := db.Update(func(tx *bolt.Tx) error {
			// Get proper size of buffer
			size := 0
			for _, t := range text {
				size += len(t)
			}
			// Encode into gob
			encoded := bytes.NewBuffer(make([]byte, 0, size))
			enc := gob.NewEncoder(encoded)
			if err := enc.Encode(text); err != nil { return err }
			// Write to bucket
			if err := tx.Bucket([]byte("SPF")).Put([]byte(body["name"].(string)), encoded.Bytes()); err != nil { return err}
			return nil
		}); err != nil {
			util.Responses.Error(w, http.StatusBadRequest, "failed to write record to database: " + err.Error())
			return
		}
	case "TXT":
		if !util.Exists(body, "text") {
			util.Responses.Error(w, http.StatusBadRequest, "field 'text' is required")
			return
		} else if !util.Types.StringArray(body["text"]) {
			util.Responses.Error(w, http.StatusBadRequest, "field 'text' must be an array of strings")
			return
		}
		text, _ := util.ConvertArrayToString(body["text"].([]interface{}))
		if len(text) < 1 {
			util.Responses.Error(w, http.StatusBadRequest, "field 'text' must a length of length 1")
			return
		} else if err := db.Update(func(tx *bolt.Tx) error {
			// Get proper size of buffer
			size := 0
			for _, t := range text {
				size += len(t)
			}
			// Encode into gob
			encoded := bytes.NewBuffer(make([]byte, 0, size))
			enc := gob.NewEncoder(encoded)
			if err := enc.Encode(text); err != nil { return err }
			// Write to bucket
			if err := tx.Bucket([]byte("TXT")).Put([]byte(body["name"].(string)), encoded.Bytes()); err != nil { return err}
			return nil
		}); err != nil {
			util.Responses.Error(w, http.StatusBadRequest, "failed to write record to database: " + err.Error())
			return
		}
	case "NS":
		if !util.Exists(body, "nameserver") {
			util.Responses.Error(w, http.StatusBadRequest, "field 'nameserver' is required")
			return
		} else if !util.Types.String(body["nameserver"]) {
			util.Responses.Error(w, http.StatusBadRequest, "field 'nameserver' must be a string")
			return
		} else if err := db.Update(func(tx *bolt.Tx) error {
			records := tx.Bucket([]byte("NS"))
			return records.Put([]byte(body["name"].(string)), []byte(body["nameserver"].(string)))
		}); err != nil {
			util.Responses.Error(w, http.StatusInternalServerError, "failed to write record to database: " + err.Error())
			return
		}
	case "CAA":
		if !util.Exists(body, "tag") {
			util.Responses.Error(w, http.StatusBadRequest, "field 'tag' is required")
			return
		} else if !util.Exists(body, "content") {
			util.Responses.Error(w, http.StatusBadRequest, "field 'content' is required")
			return
		} else if !util.Types.String(body["tag"]) {
			util.Responses.Error(w, http.StatusBadRequest, "field 'tag' must be a string")
			return
		} else if !util.Types.String(body["content"]) {
			util.Responses.Error(w, http.StatusBadRequest, "field 'content' must be a string")
			return
		} else if err := db.Update(func(tx *bolt.Tx) error {
			records := tx.Bucket([]byte("CAA"))
			if err := records.Put([]byte(body["name"].(string) + "*tag"), []byte(body["tag"].(string))); err != nil { return err }
			if err := records.Put([]byte(body["name"].(string) + "*content"), []byte(body["content"].(string))); err != nil { return err }
			return nil
		}); err != nil {
			util.Responses.Error(w, http.StatusInternalServerError, "failed to write record to database: " + err.Error())
			return
		}
	case "PTR":
		if !util.Exists(body, "domain") {
			util.Responses.Error(w, http.StatusBadRequest, "field 'domain' is required")
			return
		} else if !util.Types.String(body["domain"]) {
			util.Responses.Error(w, http.StatusBadRequest, "field 'domain' must be a string")
			return
		} else if err := db.Update(func(tx *bolt.Tx) error {
			records := tx.Bucket([]byte("PTR"))
			return records.Put([]byte(body["name"].(string)), []byte(body["domain"].(string)))
		}); err != nil {
			util.Responses.Error(w, http.StatusInternalServerError, "failed to write record to database: " + err.Error())
			return
		}
	case "CERT":
		if !util.Exists(body, "c-type") {
			util.Responses.Error(w, http.StatusBadRequest, "field 'c-type' is required")
			return
		} else if !util.Exists(body, "key-tag") {
			util.Responses.Error(w, http.StatusBadRequest, "field 'key-tag' is required")
			return
		} else if !util.Exists(body, "algorithm") {
			util.Responses.Error(w, http.StatusBadRequest, "field 'algorithm' is required")
			return
		} else if !util.Exists(body, "certificate") {
			util.Responses.Error(w, http.StatusBadRequest, "field 'certificate' is required")
			return
		} else if !util.Types.Uint16(body["c-type"]) {
			util.Responses.Error(w, http.StatusBadRequest, "field 'c-type' must be an integer between 0 and 65535")
			return
		} else if !util.Types.Uint16(body["key-tag"]) {
			util.Responses.Error(w, http.StatusBadRequest, "field 'key-tag' must be an integer between 0 and 65535")
			return
		} else if !util.Types.Uint8(body["algorithm"]) {
			util.Responses.Error(w, http.StatusBadRequest, "field 'algorithm' must be an integer between 0 and 255")
			return
		} else if !util.Types.String(body["certificate"]) {
			util.Responses.Error(w, http.StatusBadRequest, "field 'certificate' must be a string")
			return
		} else if err := db.Update(func(tx *bolt.Tx) error {
			records := tx.Bucket([]byte("CERT"))
			// Convert uint16s to binary
			ctype := make([]byte, binary.MaxVarintLen16)
			keytag := make([]byte, binary.MaxVarintLen16)
			binary.BigEndian.PutUint16(ctype, uint16(body["c-type"].(float64)))
			binary.BigEndian.PutUint16(keytag, uint16(body["key-tag"].(float64)))
			// Write data to bucket
			if err := records.Put([]byte(body["name"].(string) + "*type"), ctype); err != nil { return err }
			if err := records.Put([]byte(body["name"].(string) + "*keytag"), keytag); err != nil { return err }
			if err := records.Put([]byte(body["name"].(string) + "*algorithm"), []byte{uint8(body["algorithm"].(float64))}); err != nil { return err }
			if err := records.Put([]byte(body["name"].(string) + "*certificate"), []byte(body["certificate"].(string))); err != nil { return err }
			return nil
		}); err != nil {
			util.Responses.Error(w, http.StatusInternalServerError, "failed to write record to database: " + err.Error())
			return
		}
	case "DNSKEY":
		if !util.Exists(body, "flags") {
			util.Responses.Error(w, http.StatusBadRequest, "field 'flags' is required")
			return
		} else if !util.Exists(body, "protocol") {
			util.Responses.Error(w, http.StatusBadRequest, "field 'protocol' is required")
			return
		} else if !util.Exists(body, "algorithm") {
			util.Responses.Error(w, http.StatusBadRequest, "field 'algorithm' is required")
			return
		} else if !util.Exists(body, "public-key") {
			util.Responses.Error(w, http.StatusBadRequest, "field 'public-key' is required")
			return
		} else if !util.Types.Uint16(body["flags"]) {
			util.Responses.Error(w, http.StatusBadRequest, "field 'flags' must be a uint16")
			return
		} else if !util.Types.Uint8(body["protocol"]) {
			util.Responses.Error(w, http.StatusBadRequest, "field 'protocol' must be a uint8")
			return
		} else if !util.Types.Uint8(body["algorithm"]) {
			util.Responses.Error(w, http.StatusBadRequest, "field 'algorithm' must be a uint8")
			return
		} else if !util.Types.String(body["public-key"]) {
			util.Responses.Error(w, http.StatusBadRequest, "field 'public-key' must be a string")
			return
		} else if err := db.Update(func(tx *bolt.Tx) error {
			records := tx.Bucket([]byte("DNSKEY"))
			// Convert uint16 to binary
			flags := make([]byte, binary.MaxVarintLen16)
			binary.BigEndian.PutUint16(flags, uint16(body["flags"].(float64)))
			if err := records.Put([]byte(body["name"].(string) + "*flags"), flags); err != nil { return err }
			if err := records.Put([]byte(body["name"].(string) + "*protocol"), []byte{uint8(body["protocol"].(float64))}); err != nil { return err }
			if err := records.Put([]byte(body["name"].(string) + "*algorithm"), []byte{uint8(body["algorithm"].(float64))}); err != nil { return err }
			if err := records.Put([]byte(body["name"].(string) + "*publickey"), []byte(body["public-key"].(string))); err != nil { return err }
			return nil
		}); err != nil {
			util.Responses.Error(w, http.StatusInternalServerError, "failed to write record to database: " + err.Error())
			return
		}
	case "DS":
		if !util.Exists(body, "key-tag") {
			util.Responses.Error(w, http.StatusBadRequest, "field 'key-tag' is required")
			return
		} else if !util.Exists(body, "algorithm") {
			util.Responses.Error(w, http.StatusBadRequest, "field 'algorithm' is required")
			return
		} else if !util.Exists(body, "digest-type") {
			util.Responses.Error(w, http.StatusBadRequest, "field 'digest-type' is required")
			return
		} else if !util.Exists(body, "digest") {
			util.Responses.Error(w, http.StatusBadRequest, "field 'digest' is required")
			return
		} else if !util.Types.Uint16(body["key-tag"]) {
			util.Responses.Error(w, http.StatusBadRequest, "field 'key-tag' must be a uint16")
			return
		} else if !util.Types.Uint8(body["algorithm"]) {
			util.Responses.Error(w, http.StatusBadRequest, "field 'algorithm' must be a uint8")
			return
		} else if !util.Types.Uint8(body["digest-type"]) {
			util.Responses.Error(w, http.StatusBadRequest, "field 'digest-type' must be a uint8")
			return
		} else if !util.Types.String(body["digest"]) {
			util.Responses.Error(w, http.StatusBadRequest, "field 'digest' must be a string")
			return
		} else if err := db.Update(func(tx *bolt.Tx) error {
			records := tx.Bucket([]byte("DS"))
			// Convert uint16 to binary
			flags := make([]byte, binary.MaxVarintLen16)
			binary.BigEndian.PutUint16(flags, uint16(body["key-tag"].(float64)))
			if err := records.Put([]byte(body["name"].(string) + "*keytag"), flags); err != nil { return err }
			if err := records.Put([]byte(body["name"].(string) + "*algorithm"), []byte{uint8(body["algorithm"].(float64))}); err != nil { return err }
			if err := records.Put([]byte(body["name"].(string) + "*digesttype"), []byte{uint8(body["digest-type"].(float64))}); err != nil { return err }
			if err := records.Put([]byte(body["name"].(string) + "*digest"), []byte(body["digest"].(string))); err != nil { return err }
			return nil
		}); err != nil {
			util.Responses.Error(w, http.StatusInternalServerError, "failed to write record to database: " + err.Error())
			return
		}
	case "NAPTR":
		if !util.Exists(body, "order") {
			util.Responses.Error(w, http.StatusBadRequest, "field 'order' is required")
			return
		} else if !util.Exists(body, "preference") {
			util.Responses.Error(w, http.StatusBadRequest, "field 'preference' is required")
			return
		} else if !util.Exists(body, "flags") {
			util.Responses.Error(w, http.StatusBadRequest, "field 'flags' is required")
			return
		} else if !util.Exists(body, "service") {
			util.Responses.Error(w, http.StatusBadRequest, "field 'service' is required")
			return
		} else if !util.Exists(body, "regexp") {
			util.Responses.Error(w, http.StatusBadRequest, "field 'regexp' is required")
			return
		} else if !util.Exists(body, "replacement") {
			util.Responses.Error(w, http.StatusBadRequest, "field 'replacement' is required")
			return
		} else if !util.Types.Uint16(body["order"]) {
			util.Responses.Error(w, http.StatusBadRequest, "field 'order' must be an integer between 0 and 65535")
			return
		} else if !util.Types.Uint16(body["preference"]) {
			util.Responses.Error(w, http.StatusBadRequest, "field 'preference' must be an integer between 0 and 65535")
			return
		} else if !util.Types.String(body["flags"]) {
			util.Responses.Error(w, http.StatusBadRequest, "field 'flags' must be a string")
			return
		} else if !util.Types.String(body["service"]) {
			util.Responses.Error(w, http.StatusBadRequest, "field 'service' must be a string")
			return
		} else if !util.Types.String(body["regexp"]) {
			util.Responses.Error(w, http.StatusBadRequest, "field 'regexp' must be a string")
			return
		} else if !util.Types.String(body["replacement"]) {
			util.Responses.Error(w, http.StatusBadRequest, "field 'replacement' must be a string")
			return
		} else if err := db.Update(func(tx *bolt.Tx) error {
			records := tx.Bucket([]byte("NAPTR"))
			// Convert uint16s to binary
			order := make([]byte, binary.MaxVarintLen16)
			preference := make([]byte, binary.MaxVarintLen16)
			binary.BigEndian.PutUint16(order, uint16(body["order"].(float64)))
			binary.BigEndian.PutUint16(preference, uint16(body["preference"].(float64)))
			// Write data to bucket
			if err := records.Put([]byte(body["name"].(string) + "*order"), order); err != nil { return err }
			if err := records.Put([]byte(body["name"].(string) + "*preference"), preference); err != nil { return err }
			if err := records.Put([]byte(body["name"].(string) + "*flags"), []byte(body["flags"].(string))); err != nil { return err }
			if err := records.Put([]byte(body["name"].(string) + "*service"), []byte(body["service"].(string))); err != nil { return err }
			if err := records.Put([]byte(body["name"].(string) + "*regexp"), []byte(body["regexp"].(string))); err != nil { return err }
			if err := records.Put([]byte(body["name"].(string) + "*replacement"), []byte(body["replacement"].(string))); err != nil { return err }
			return nil
		}); err != nil {
			util.Responses.Error(w, http.StatusInternalServerError, "failed to write record to database: " + err.Error())
			return
		}
	case "SMIMEA":
		if !util.Exists(body, "usage") {
			util.Responses.Error(w, http.StatusBadRequest, "field 'usage' is required")
			return
		} else if !util.Exists(body, "selector") {
			util.Responses.Error(w, http.StatusBadRequest, "field 'selector' is required")
			return
		} else if !util.Exists(body, "matching-type") {
			util.Responses.Error(w, http.StatusBadRequest, "field 'matching-type' is required")
			return
		} else if !util.Exists(body, "certificate") {
			util.Responses.Error(w, http.StatusBadRequest, "field 'certificate' is required")
			return
		} else if !util.Types.Uint8(body["usage"]) {
			util.Responses.Error(w, http.StatusBadRequest, "field 'usage' must be a uint16")
			return
		} else if !util.Types.Uint8(body["selector"]) {
			util.Responses.Error(w, http.StatusBadRequest, "field 'selector' must be a uint8")
			return
		} else if !util.Types.Uint8(body["matching-type"]) {
			util.Responses.Error(w, http.StatusBadRequest, "field 'matching-type' must be a uint8")
			return
		} else if !util.Types.String(body["certificate"]) {
			util.Responses.Error(w, http.StatusBadRequest, "field 'certificate' must be a string")
			return
		} else if err := db.Update(func(tx *bolt.Tx) error {
			records := tx.Bucket([]byte("SMIMEA"))
			if err := records.Put([]byte(body["name"].(string) + "*usage"), []byte{uint8(body["usage"].(float64))}); err != nil { return err }
			if err := records.Put([]byte(body["name"].(string) + "*selector"), []byte{uint8(body["selector"].(float64))}); err != nil { return err }
			if err := records.Put([]byte(body["name"].(string) + "*matching"), []byte{uint8(body["matching-type"].(float64))}); err != nil { return err }
			if err := records.Put([]byte(body["name"].(string) + "*certificate"), []byte(body["certificate"].(string))); err != nil { return err }
			return nil
		}); err != nil {
			util.Responses.Error(w, http.StatusInternalServerError, "failed to write record to database: " + err.Error())
			return
		}
	case "SSHFP":
		if !util.Exists(body, "algorithm") {
			util.Responses.Error(w, http.StatusBadRequest, "field 'algorithm' is required")
			return
		} else if !util.Exists(body, "s-type") {
			util.Responses.Error(w, http.StatusBadRequest, "field 's-type' is required")
			return
		} else if !util.Exists(body, "fingerprint") {
			util.Responses.Error(w, http.StatusBadRequest, "field 'fingerprint' is required")
			return
		} else if !util.Types.Uint8(body["algorithm"]) {
			util.Responses.Error(w, http.StatusBadRequest, "field 'algorithm' must be a uint16")
			return
		} else if !util.Types.Uint8(body["s-type"]) {
			util.Responses.Error(w, http.StatusBadRequest, "field 's-type' must be a uint8")
			return
		} else if !util.Types.String(body["fingerprint"]) {
			util.Responses.Error(w, http.StatusBadRequest, "field 'fingerprint' must be a string")
			return
		} else if err := db.Update(func(tx *bolt.Tx) error {
			records := tx.Bucket([]byte("SSHFP"))
			if err := records.Put([]byte(body["name"].(string) + "*algorithm"), []byte{uint8(body["algorithm"].(float64))}); err != nil { return err }
			if err := records.Put([]byte(body["name"].(string) + "*type"), []byte{uint8(body["s-type"].(float64))}); err != nil { return err }
			if err := records.Put([]byte(body["name"].(string) + "*fingerprint"), []byte(body["fingerprint"].(string))); err != nil { return err }
			return nil
		}); err != nil {
			util.Responses.Error(w, http.StatusInternalServerError, "failed to write record to database: " + err.Error())
			return
		}
	case "TLSA":
		if !util.Exists(body, "usage") {
			util.Responses.Error(w, http.StatusBadRequest, "field 'usage' is required")
			return
		} else if !util.Exists(body, "selector") {
			util.Responses.Error(w, http.StatusBadRequest, "field 'selector' is required")
			return
		} else if !util.Exists(body, "matching-type") {
			util.Responses.Error(w, http.StatusBadRequest, "field 'matching-type' is required")
			return
		} else if !util.Exists(body, "certificate") {
			util.Responses.Error(w, http.StatusBadRequest, "field 'certificate' is required")
			return
		} else if !util.Types.Uint8(body["usage"]) {
			util.Responses.Error(w, http.StatusBadRequest, "field 'usage' must be a uint16")
			return
		} else if !util.Types.Uint8(body["selector"]) {
			util.Responses.Error(w, http.StatusBadRequest, "field 'selector' must be a uint8")
			return
		} else if !util.Types.Uint8(body["matching-type"]) {
			util.Responses.Error(w, http.StatusBadRequest, "field 'matching-type' must be a uint8")
			return
		} else if !util.Types.String(body["certificate"]) {
			util.Responses.Error(w, http.StatusBadRequest, "field 'certificate' must be a string")
			return
		} else if err := db.Update(func(tx *bolt.Tx) error {
			records := tx.Bucket([]byte("TLSA"))
			if err := records.Put([]byte(body["name"].(string) + "*usage"), []byte{uint8(body["usage"].(float64))}); err != nil { return err }
			if err := records.Put([]byte(body["name"].(string) + "*selector"), []byte{uint8(body["selector"].(float64))}); err != nil { return err }
			if err := records.Put([]byte(body["name"].(string) + "*matching"), []byte{uint8(body["matching-type"].(float64))}); err != nil { return err }
			if err := records.Put([]byte(body["name"].(string) + "*certificate"), []byte(body["certificate"].(string))); err != nil { return err }
			return nil
		}); err != nil {
			util.Responses.Error(w, http.StatusInternalServerError, "failed to write record to database: " + err.Error())
			return
		}
	case "URI":
		if !util.Exists(body, "priority") {
			util.Responses.Error(w, http.StatusBadRequest, "field 'priority' is required")
			return
		} else if !util.Exists(body, "weight") {
			util.Responses.Error(w, http.StatusBadRequest, "field 'weight' is required")
			return
		} else if !util.Exists(body, "target") {
			util.Responses.Error(w, http.StatusBadRequest, "field 'target' is required")
			return
		} else if !util.Types.Uint16(body["priority"]) {
			util.Responses.Error(w, http.StatusBadRequest, "field 'priority' must be an integer between 0 and 65535")
			return
		} else if !util.Types.Uint16(body["weight"]) {
			util.Responses.Error(w, http.StatusBadRequest, "field 'weight' must be an integer between 0 and 65535")
			return
		} else if !util.Types.String(body["target"]) {
			util.Responses.Error(w, http.StatusBadRequest, "field 'target' must be a string")
			return
		} else if err := db.Update(func(tx *bolt.Tx) error {
			records := tx.Bucket([]byte("URI"))
			// Convert uint16s to binary
			priority := make([]byte, binary.MaxVarintLen16)
			weight := make([]byte, binary.MaxVarintLen16)
			binary.BigEndian.PutUint16(priority, uint16(body["priority"].(float64)))
			binary.BigEndian.PutUint16(weight, uint16(body["weight"].(float64)))
			// Write data to bucket
			if err := records.Put([]byte(body["name"].(string) + "*priority"), priority); err != nil { return err }
			if err := records.Put([]byte(body["name"].(string) + "*weight"), weight); err != nil { return err }
			if err := records.Put([]byte(body["name"].(string) + "*target"), []byte(body["target"].(string))); err != nil { return err}
			return nil
		}); err != nil {
			util.Responses.Error(w, http.StatusInternalServerError, "failed to write record to database: " + err.Error())
			return
		}
	default:
		util.Responses.Error(w, http.StatusBadRequest, "field 'type' must be on of: A, AAAA, CNAME, MX, LOC, SRV, SPF, TXT, NS, CAA, PTR, CERT, DNSKEY, DS, NAPTR, SMIMEA, SSHFP, TLSA, URI")
		return
	}

	util.Responses.Success(w)
}
