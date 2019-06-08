package routes

import (
	"encoding/json"
	"github.com/akrantz01/krantz.dev/dns/db"
	"github.com/akrantz01/krantz.dev/dns/util"
	bolt "go.etcd.io/bbolt"
	"net"
	"net/http"
	"strings"
)

// Handle the updating of records
func update(w http.ResponseWriter, r *http.Request, path string, database *bolt.DB) {
	// Set database into getter and setter
	db.Get.Db = database
	db.Set.Db = database

	// Validate initial request with request type, body exists, and content type
	if r.Method != "PUT" {
		util.Responses.Error(w, http.StatusMethodNotAllowed, "method not allowed")
		return
	} else if r.Body == nil {
		util.Responses.Error(w, http.StatusBadRequest, "body must be present")
		return
	} else if r.Header.Get("Content-Type") != "application/json" {
		util.Responses.Error(w, http.StatusBadRequest, "body must be of type JSON")
		return
	} else if len(r.URL.Path[len(path):]) == 0 {
		util.Responses.Error(w, http.StatusBadRequest, "record must be specified in path")
	}

	// Validate body by decoding json, checking fields exists, and checking field type
	var body map[string]interface{}
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		util.Responses.Error(w, http.StatusBadRequest, "failed to decode body: "+err.Error())
		return
	} else if !util.Exists(body, "type") {
		util.Responses.Error(w, http.StatusBadRequest, "field 'type' is required")
		return
	} else if !util.Types.String(body["type"]) {
		util.Responses.Error(w, http.StatusBadRequest, "field 'type' must be a string")
		return
	}

	recordName := strings.ToLower(r.URL.Path[len(path):])

	// Parse out body by type
	switch strings.ToUpper(body["type"].(string)) {
	case "A":
		// Get original record from database
		record := db.Get.A(recordName + ".")
		if util.RecordDoesNotExist(record) {
			util.Responses.Error(w, http.StatusBadRequest, "specified record does not exist")
			return
		}

		// Update values if they exist in the body
		if util.Exists(body, "host") {
			if !util.Types.String(body["host"]) {
				util.Responses.Error(w, http.StatusBadRequest, "field 'host' must be a string")
				return
			} else if ip := net.ParseIP(body["host"].(string)); ip.To4().String() == "<nil>" {
				util.Responses.Error(w, http.StatusBadRequest, "field 'host' must be an IPv4 address")
				return
			}
			record.Address = net.ParseIP(body["host"].(string))
		}

		// Write updated values to the database
		if err := db.Set.A(recordName, body["host"].(string)); err != nil {
			util.Responses.Error(w, http.StatusInternalServerError, "failed to write record to database: "+err.Error())
			return
		}

	case "AAAA":
		// Get original record from database
		record := db.Get.AAAA(recordName + ".")
		if util.RecordDoesNotExist(record) {
			util.Responses.Error(w, http.StatusBadRequest, "specified record does not exist")
			return
		}

		// Update values if they exist in the body
		if util.Exists(body, "host") {
			if !util.Types.String(body["host"]) {
				util.Responses.Error(w, http.StatusBadRequest, "field 'host' must be a string")
				return
			} else if ip := net.ParseIP(body["host"].(string)); ip.To4().String() != "<nil>" {
				util.Responses.Error(w, http.StatusBadRequest, "field 'host' must be an IPv6 address")
				return
			}
			record.Address = net.ParseIP(body["host"].(string))
		}

		// Write updated values to the database
		if err := db.Set.AAAA(recordName, body["host"].(string)); err != nil {
			util.Responses.Error(w, http.StatusInternalServerError, "failed to write record to database: "+err.Error())
			return
		}

	case "CNAME":
		// Get original record from database
		record := db.Get.CNAME(recordName + ".")
		if util.RecordDoesNotExist(record) {
			util.Responses.Error(w, http.StatusBadRequest, "specified record does not exist")
			return
		}

		// Update values if they exist in the body
		if util.Exists(body, "target") {
			if !util.Types.String(body["target"]) {
				util.Responses.Error(w, http.StatusBadRequest, "field 'target' must be a string")
				return
			}
			record.Target = body["target"].(string)
		}

		// Write updated values to the database
		if err := db.Set.CNAME(recordName, record.Target); err != nil {
			util.Responses.Error(w, http.StatusInternalServerError, "failed to write record to database: "+err.Error())
			return
		}

	case "MX":
		// Get original record from database
		record := db.Get.MX(recordName + ".")
		if util.RecordDoesNotExist(record) {
			util.Responses.Error(w, http.StatusBadRequest, "specified record does not exist")
			return
		}

		// Update values if they exist in the body
		if util.Exists(body, "host") {
			if !util.Types.String(body["host"]) {
				util.Responses.Error(w, http.StatusBadRequest, "field 'host' must be a string")
				return
			}
			record.Host = body["host"].(string)
		}
		if util.Exists(body, "priority") {
			if !util.Types.Uint16(body["priority"]) {
				util.Responses.Error(w, http.StatusBadRequest, "field 'priority' must be an integer between 0 and 65535")
				return
			}
			record.Priority = uint16(body["priority"].(float64))
		}

		// Write updated values to the database
		if err := db.Set.MX(recordName, record.Priority, record.Host); err != nil {
			util.Responses.Error(w, http.StatusInternalServerError, "failed to write record to database: "+err.Error())
			return
		}

	case "LOC":
		// Get original record from database
		record := db.Get.LOC(recordName + ".")
		if util.RecordDoesNotExist(record) {
			util.Responses.Error(w, http.StatusBadRequest, "specified record does not exist")
			return
		}

		// Update values if they exist in body
		if util.Exists(body, "version") {
			if !util.Types.Uint8(body["version"]) {
				util.Responses.Error(w, http.StatusBadRequest, "field 'version' must be an integer between 0 and 255")
				return
			}
			record.Version = uint8(body["version"].(float64))
		}
		if util.Exists(body, "size") {
			if !util.Types.Uint8(body["size"]) {
				util.Responses.Error(w, http.StatusBadRequest, "field 'size' must be an integer between 0 and 255")
				return
			}
			record.Size = uint8(body["size"].(float64))
		}
		if util.Exists(body, "horizontal-precision") {
			if !util.Types.Uint8(body["horizontal-precision"]) {
				util.Responses.Error(w, http.StatusBadRequest, "field 'horizontal-precision' must be an integer between 0 and 255")
				return
			}
			record.HorizontalPrecision = uint8(body["horizontal-precision"].(float64))
		}
		if util.Exists(body, "vertical-precision") {
			if !util.Types.Uint8(body["vertical-precision"]) {
				util.Responses.Error(w, http.StatusBadRequest, "field 'vertical-precision' must be an integer between 0 and 255")
				return
			}
			record.VerticalPrecision = uint8(body["vertical-precision"].(float64))
		}
		if util.Exists(body, "latitude") {
			if !util.Types.Uint32(body["latitude"]) {
				util.Responses.Error(w, http.StatusBadRequest, "field 'latitude' must be an integer between 0 and 4294967295")
				return
			}
			record.Latitude = uint32(body["latitude"].(float64))
		}
		if util.Exists(body, "longitude") {
			if !util.Types.Uint32(body["longitude"]) {
				util.Responses.Error(w, http.StatusBadRequest, "field 'longitude' must be an integer between 0 and 4294967295")
				return
			}
			record.Longitude = uint32(body["longitude"].(float64))
		}
		if util.Exists(body, "altitude") {
			if !util.Types.Uint32(body["altitude"]) {
				util.Responses.Error(w, http.StatusBadRequest, "field 'altitude' must be an integer between 0 and 4294967295")
				return
			}
			record.Altitude = uint32(body["altitude"].(float64))
		}

		// Write updated values to database
		if err := db.Set.LOC(recordName, record.Version, record.Size, record.HorizontalPrecision, record.VerticalPrecision, record.Latitude, record.Longitude, record.Altitude); err != nil {
			util.Responses.Error(w, http.StatusInternalServerError, "failed to write record to database: "+err.Error())
			return
		}

	case "SRV":
		// Get original record from database
		record := db.Get.SRV(recordName + ".")
		if util.RecordDoesNotExist(record) {
			util.Responses.Error(w, http.StatusBadRequest, "specified record does not exist")
			return
		}

		// Update values if they exist in body
		if util.Exists(body, "priority") {
			if !util.Types.Uint16(body["priority"]) {
				util.Responses.Error(w, http.StatusBadRequest, "field 'priority' must be an integer between 0 and 65535")
				return
			}
			record.Priority = uint16(body["priority"].(float64))
		}
		if util.Exists(body, "weight") {
			if !util.Types.Uint16(body["weight"]) {
				util.Responses.Error(w, http.StatusBadRequest, "field 'weight' must be an integer between 0 and 65535")
				return
			}
			record.Weight = uint16(body["weight"].(float64))
		}
		if util.Exists(body, "port") {
			if !util.Types.Uint16(body["port"]) {
				util.Responses.Error(w, http.StatusBadRequest, "field 'port' must be an integer between 0 and 65535")
				return
			}
			record.Port = uint16(body["port"].(float64))
		}
		if util.Exists(body, "target") {
			if !util.Types.String(body["target"]) {
				util.Responses.Error(w, http.StatusBadRequest, "field 'target' must be a string")
				return
			}
			record.Target = body["target"].(string)
		}

		// Write updated values to database
		if err := db.Set.SRV(recordName, record.Priority, record.Weight, record.Port, record.Target); err != nil {
			util.Responses.Error(w, http.StatusInternalServerError, "failed to write record to database: "+err.Error())
			return
		}

	case "SPF":
		// Get original record from database
		record := db.Get.SPF(recordName + ".")
		if util.RecordDoesNotExist(record) {
			util.Responses.Error(w, http.StatusBadRequest, "specified record does not exist")
			return
		}

		// Update values if they exist in body
		if util.Exists(body, "text") {
			if !util.Types.StringArray(body["text"]) {
				util.Responses.Error(w, http.StatusBadRequest, "field 'text' must be an array of strings")
				return
			}
			text, _ := util.ConvertArrayToString(body["text"].([]interface{}))
			if len(text) < 1 {
				util.Responses.Error(w, http.StatusBadRequest, "field 'text' must be a length of 1")
				return
			}
			record.Text = text
		}

		// Write updated values to database
		if err := db.Set.SPF(recordName, record.Text); err != nil {
			util.Responses.Error(w, http.StatusInternalServerError, "failed to write record to database: "+err.Error())
			return
		}

	case "TXT":
		// Get original record from database
		record := db.Get.TXT(recordName + ".")
		if util.RecordDoesNotExist(record) {
			util.Responses.Error(w, http.StatusBadRequest, "specified record does not exist")
			return
		}

		// Update values if they exist in body
		if util.Exists(body, "text") {
			if !util.Types.StringArray(body["text"]) {
				util.Responses.Error(w, http.StatusBadRequest, "field 'text' must be an array of strings")
				return
			}
			text, _ := util.ConvertArrayToString(body["text"].([]interface{}))
			if len(text) < 1 {
				util.Responses.Error(w, http.StatusBadRequest, "field 'text' must be a length of 1")
				return
			}
			record.Text = text
		}

		// Write updated values to database
		if err := db.Set.TXT(recordName, record.Text); err != nil {
			util.Responses.Error(w, http.StatusInternalServerError, "failed to write record to database: "+err.Error())
			return
		}

	case "NS":
		// Get original record from database
		record := db.Get.NS(recordName + ".")
		if util.RecordDoesNotExist(record) {
			util.Responses.Error(w, http.StatusBadRequest, "specified record does not exist")
			return
		}

		// Update values if they exist in body
		if util.Exists(body, "nameserver") {
			if !util.Types.String(body["nameserver"]) {
				util.Responses.Error(w, http.StatusBadRequest, "field 'nameserver' must be a string")
				return
			}
			record.Nameserver = body["nameserver"].(string)
		}

		// Write updated values to database
		if err := db.Set.NS(recordName, record.Nameserver); err != nil {
			util.Responses.Error(w, http.StatusInternalServerError, "failed to write record to database: "+ err.Error())
			return
		}

	case "CAA":
		// Get original record from database
		record := db.Get.CAA(recordName + ".")
		if util.RecordDoesNotExist(record) {
			util.Responses.Error(w, http.StatusBadRequest, "specified record does not exist")
			return
		}

		// Update values if they exist in body
		if util.Exists(body, "tag") {
			if !util.Types.String(body["tag"]) {
				util.Responses.Error(w, http.StatusBadRequest, "field 'tag' must be a string")
				return
			}
			record.Tag = body["tag"].(string)
		}
		if util.Exists(body, "content") {
			if !util.Types.String(body["content"]) {
				util.Responses.Error(w, http.StatusBadRequest, "field 'content' must be a string")
				return
			}
			record.Tag = body["content"].(string)
		}

		// Write updated values to database
		if err := db.Set.CAA(recordName, record.Tag, record.Content); err != nil {
			util.Responses.Error(w, http.StatusInternalServerError, "failed to write record to database: "+ err.Error())
			return
		}

	case "PTR":
		// Get original record from database
		record := db.Get.PTR(recordName + ".")
		if util.RecordDoesNotExist(record) {
			util.Responses.Error(w, http.StatusBadRequest, "specified record does not exist")
			return
		}

		// Update values if they exist in body
		if util.Exists(body, "domain") {
			if !util.Types.String(body["domain"]) {
				util.Responses.Error(w, http.StatusBadRequest, "field 'domain' must be a string")
				return
			}
			record.Domain = body["domain"].(string)
		}

		// Write updated values to database
		if err := db.Set.PTR(recordName, record.Domain); err != nil {
			util.Responses.Error(w, http.StatusInternalServerError, "failed to write record to database: "+ err.Error())
			return
		}

	case "CERT":
		// Get original record from database
		record := db.Get.CERT(recordName + ".")
		if util.RecordDoesNotExist(record) {
			util.Responses.Error(w, http.StatusBadRequest, "specified record does not exist")
			return
		}

		// Update values if they exist in body
		if util.Exists(body, "c-type") {
			if !util.Types.Uint16(body["c-type"]) {
				util.Responses.Error(w, http.StatusBadRequest, "field 'c-type' must be an integer between 0 and 65535")
				return
			}
			record.Type = uint16(body["c-type"].(float64))
		}
		if util.Exists(body, "key-tag") {
			if !util.Types.Uint16(body["key-tag"]) {
				util.Responses.Error(w, http.StatusBadRequest, "field 'key-tag' must be an integer between 0 and 65535")
				return
			}
			record.KeyTag = uint16(body["key-tag"].(float64))
		}
		if util.Exists(body, "algorithm") {
			if !util.Types.Uint8(body["algorithm"]) {
				util.Responses.Error(w, http.StatusBadRequest, "field 'algorithm' must be an integer between 0 and 255")
				return
			}
			record.Algorithm = uint8(body["algorithm"].(float64))
		}
		if util.Exists(body, "certificate") {
			if !util.Types.String(body["certificate"]) {
				util.Responses.Error(w, http.StatusBadRequest, "field 'certificate' must be a string")
				return
			}
			record.Certificate = body["certificate"].(string)
		}

		// Write updated values to database
		if err := db.Set.CERT(recordName, record.Type, record.KeyTag, record.Algorithm, record.Certificate); err != nil {
			util.Responses.Error(w, http.StatusInternalServerError, "failed to write record to database: "+err.Error())
			return
		}

	case "DNSKEY":
		// Get original record from database
		record := db.Get.DNSKEY(recordName + ".")
		if util.RecordDoesNotExist(record) {
			util.Responses.Error(w, http.StatusBadRequest, "specified record does not exist")
			return
		}

		// Update values if they exist in body
		if util.Exists(body, "flags") {
			if !util.Types.Uint16(body["flags"]) {
				util.Responses.Error(w, http.StatusBadRequest, "field 'flags' must be an integer between 0 and 65535")
				return
			}
			record.Flags = uint16(body["flags"].(float64))
		}
		if util.Exists(body, "protocol") {
			if !util.Types.Uint8(body["protocol"]) {
				util.Responses.Error(w, http.StatusBadRequest, "field 'protocol' must be an integer between 0 and 65535")
				return
			}
			record.Protocol = uint8(body["protocol"].(float64))
		}
		if util.Exists(body, "algorithm") {
			if !util.Types.Uint8(body["algorithm"]) {
				util.Responses.Error(w, http.StatusBadRequest, "field 'algorithm' must be an integer between 0 and 255")
				return
			}
			record.Algorithm = uint8(body["algorithm"].(float64))
		}
		if util.Exists(body, "public-key") {
			if !util.Types.String(body["public-key"]) {
				util.Responses.Error(w, http.StatusBadRequest, "field 'public-key' must be a string")
				return
			}
			record.PublicKey = body["public-key"].(string)
		}

		// Write updated values to database
		if err := db.Set.DNSKEY(recordName, record.Flags, record.Protocol, record.Algorithm, record.PublicKey); err != nil {
			util.Responses.Error(w, http.StatusInternalServerError, "failed to write record to database: "+err.Error())
			return
		}

	case "DS":
		// Get original record from database
		record := db.Get.DS(recordName + ".")
		if util.RecordDoesNotExist(record) {
			util.Responses.Error(w, http.StatusBadRequest, "specified record does not exist")
			return
		}

		// Update values if they exist in body
		if util.Exists(body, "key-tag") {
			if !util.Types.Uint16(body["key-tag"]) {
				util.Responses.Error(w, http.StatusBadRequest, "field 'key-tag' must be an integer between 0 and 65535")
				return
			}
			record.KeyTag = uint16(body["key-tag"].(float64))
		}
		if util.Exists(body, "algorithm") {
			if !util.Types.Uint8(body["algorithm"]) {
				util.Responses.Error(w, http.StatusBadRequest, "field 'algorithm' must be an integer between 0 and 65535")
				return
			}
			record.Algorithm = uint8(body["algorithm"].(float64))
		}
		if util.Exists(body, "digest-type") {
			if !util.Types.Uint8(body["digest-type"]) {
				util.Responses.Error(w, http.StatusBadRequest, "field 'digest-type' must be an integer between 0 and 255")
				return
			}
			record.DigestType = uint8(body["digest-type"].(float64))
		}
		if util.Exists(body, "digest") {
			if !util.Types.String(body["digest"]) {
				util.Responses.Error(w, http.StatusBadRequest, "field 'digest' must be a string")
				return
			}
			record.Digest = body["digest"].(string)
		}

		// Write updated values to database
		if err := db.Set.DS(recordName, record.KeyTag, record.Algorithm, record.DigestType, record.Digest); err != nil {
			util.Responses.Error(w, http.StatusInternalServerError, "failed to write record to database: "+err.Error())
			return
		}

	case "NAPTR":
		// Get original record from database
		record := db.Get.NAPTR(recordName + ".")
		if util.RecordDoesNotExist(record) {
			util.Responses.Error(w, http.StatusBadRequest, "specified record does not exist")
			return
		}

		// Update values if they exist in body
		if util.Exists(body, "order") {
			if !util.Types.Uint16(body["order"]) {
				util.Responses.Error(w, http.StatusBadRequest, "field 'order' must be an integer between 0 and 65535")
				return
			}
			record.Order = uint16(body["order"].(float64))
		}
		if util.Exists(body, "preference") {
			if !util.Types.Uint16(body["preference"]) {
				util.Responses.Error(w, http.StatusBadRequest, "field 'preference' must be an integer between 0 and 65535")
				return
			}
			record.Preference = uint16(body["preference"].(float64))
		}
		if util.Exists(body, "flags") {
			if !util.Types.String(body["flags"]) {
				util.Responses.Error(w, http.StatusBadRequest, "field 'flags' must be a string")
				return
			}
			record.Flags = body["flags"].(string)
		}
		if util.Exists(body, "service") {
			if !util.Types.String(body["service"]) {
				util.Responses.Error(w, http.StatusBadRequest, "field 'service' must be a string")
				return
			}
			record.Service = body["service"].(string)
		}
		if util.Exists(body, "regexp") {
			if !util.Types.String(body["regexp"]) {
				util.Responses.Error(w, http.StatusBadRequest, "field 'regexp' must be a string")
				return
			}
			record.Regexp = body["regexp"].(string)
		}
		if util.Exists(body, "replacement") {
			if !util.Types.String(body["replacement"]) {
				util.Responses.Error(w, http.StatusBadRequest, "field 'replacement' must be a string")
				return
			}
			record.Replacement = body["replacement"].(string)
		}

		// Write updated values to database
		if err := db.Set.NAPTR(recordName, record.Order, record.Preference, record.Flags, record.Service, record.Regexp, record.Replacement); err != nil {
			util.Responses.Error(w, http.StatusInternalServerError, "failed to write record to database: "+err.Error())
			return
		}

	case "SMIMEA":
		// Get original record from database
		record := db.Get.SMIMEA(recordName + ".")
		if util.RecordDoesNotExist(record) {
			util.Responses.Error(w, http.StatusBadRequest, "specified record does not exist")
			return
		}

		// Update values if they exist in body
		if util.Exists(body, "usage") {
			if !util.Types.Uint8(body["usage"]) {
				util.Responses.Error(w, http.StatusBadRequest, "field 'usage' must be an integer between 0 and 255")
				return
			}
			record.Usage = uint8(body["usage"].(float64))
		}
		if util.Exists(body, "selector") {
			if !util.Types.Uint8(body["selector"]) {
				util.Responses.Error(w, http.StatusBadRequest, "field 'selector' must be an integer between 0 and 255")
				return
			}
			record.Selector = uint8(body["selector"].(float64))
		}
		if util.Exists(body, "matching-type") {
			if !util.Types.Uint8(body["matching-type"]) {
				util.Responses.Error(w, http.StatusBadRequest, "field 'matching-type' must be an integer between 0 and 255")
				return
			}
			record.MatchingType = uint8(body["matching-type"].(float64))
		}
		if util.Exists(body, "certificate") {
			if !util.Types.String(body["certificate"]) {
				util.Responses.Error(w, http.StatusBadRequest, "field 'certificate' must be a string")
				return
			}
			record.Certificate = body["certificate"].(string)
		}

		// Write updated values to database
		if err := db.Set.SMIMEA(recordName, record.Usage, record.Selector, record.MatchingType, record.Certificate); err != nil {
			util.Responses.Error(w, http.StatusInternalServerError, "failed to write record to database: "+err.Error())
			return
		}

	case "SSHFP":
		// Get original record from database
		record := db.Get.SSHFP(recordName + ".")
		if util.RecordDoesNotExist(record) {
			util.Responses.Error(w, http.StatusBadRequest, "specified record does not exist")
			return
		}

		// Update values if they exist in body
		if util.Exists(body, "algorithm") {
			if !util.Types.Uint8(body["algorithm"]) {
				util.Responses.Error(w, http.StatusBadRequest, "field 'algorithm' must be an integer between 0 and 255")
				return
			}
			record.Algorithm = uint8(body["algorithm"].(float64))
		}
		if util.Exists(body, "s-type") {
			if !util.Types.Uint8(body["s-type"]) {
				util.Responses.Error(w, http.StatusBadRequest, "field 's-type' must be an integer between 0 and 255")
				return
			}
			record.Type = uint8(body["s-type"].(float64))
		}
		if util.Exists(body, "fingerprint") {
			if !util.Types.String(body["fingerprint"]) {
				util.Responses.Error(w, http.StatusBadRequest, "field 'fingerprint' must be a string")
				return
			}
			record.Fingerprint = body["fingerprint"].(string)
		}

		// Write updated values to database
		if err := db.Set.SSHFP(recordName, record.Algorithm, record.Type, record.Fingerprint); err != nil {
			util.Responses.Error(w, http.StatusInternalServerError, "failed to write record to database: "+err.Error())
			return
		}

	case "TLSA":
		// Get original record from database
		record := db.Get.TLSA(recordName + ".")
		if util.RecordDoesNotExist(record) {
			util.Responses.Error(w, http.StatusBadRequest, "specified record does not exist")
			return
		}

		// Update values if they exist in body
		if util.Exists(body, "usage") {
			if !util.Types.Uint8(body["usage"]) {
				util.Responses.Error(w, http.StatusBadRequest, "field 'usage' must be an integer between 0 and 255")
				return
			}
			record.Usage = uint8(body["usage"].(float64))
		}
		if util.Exists(body, "selector") {
			if !util.Types.Uint8(body["selector"]) {
				util.Responses.Error(w, http.StatusBadRequest, "field 'selector' must be an integer between 0 and 255")
				return
			}
			record.Selector = uint8(body["selector"].(float64))
		}
		if util.Exists(body, "matching-type") {
			if !util.Types.Uint8(body["matching-type"]) {
				util.Responses.Error(w, http.StatusBadRequest, "field 'matching-type' must be an integer between 0 and 255")
				return
			}
			record.MatchingType = uint8(body["matching-type"].(float64))
		}
		if util.Exists(body, "certificate") {
			if !util.Types.String(body["certificate"]) {
				util.Responses.Error(w, http.StatusBadRequest, "field 'certificate' must be a string")
				return
			}
			record.Certificate = body["certificate"].(string)
		}

		// Write updated values to database
		if err := db.Set.TLSA(recordName, record.Usage, record.Selector, record.MatchingType, record.Certificate); err != nil {
			util.Responses.Error(w, http.StatusInternalServerError, "failed to write record to database: "+err.Error())
			return
		}

	case "URI":
		// Get original record from database
		record := db.Get.URI(recordName + ".")
		if util.RecordDoesNotExist(record) {
			util.Responses.Error(w, http.StatusBadRequest, "specified record does not exist")
			return
		}

		// Update values if they exist in body
		if util.Exists(body, "priority") {
			if !util.Types.Uint16(body["priority"]) {
				util.Responses.Error(w, http.StatusBadRequest, "field 'priority' must be an integer between 0 and 65535")
				return
			}
			record.Priority = uint16(body["priority"].(float64))
		}
		if util.Exists(body, "weight") {
			if !util.Types.Uint16(body["weight"]) {
				util.Responses.Error(w, http.StatusBadRequest, "field 'weight' must be an integer between 0 and 65535")
				return
			}
			record.Weight = uint16(body["weight"].(float64))
		}
		if util.Exists(body, "target") {
			if !util.Types.String(body["target"]) {
				util.Responses.Error(w, http.StatusBadRequest, "field 'target' must be a string")
				return
			}
			record.Target = body["target"].(string)
		}

		// Write updated values to database
		if err := db.Set.URI(recordName, record.Priority, record.Weight, record.Target); err != nil {
			util.Responses.Error(w, http.StatusInternalServerError, "failed to write record to database: "+err.Error())
			return
		}
	default:
		util.Responses.Error(w, http.StatusBadRequest, "field 'type' must be on of: A, AAAA, CNAME, MX, LOC, SRV, SPF, TXT, NS, CAA, PTR, CERT, DNSKEY, DS, NAPTR, SMIMEA, SSHFP, TLSA, URI")
		return
	}

	util.Responses.Success(w)
}
