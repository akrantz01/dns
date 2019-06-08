package routes

import (
	"encoding/json"
	"github.com/akrantz01/krantz.dev/dns/db"
	"github.com/akrantz01/krantz.dev/dns/util"
	bolt "go.etcd.io/bbolt"
	"net/http"
	"strings"
)

// Handle the creation of records
func create(w http.ResponseWriter, r *http.Request, database *bolt.DB) {
	// Set database into operations
	db.Get.Db = database
	db.Set.Db = database
	db.Delete.Db = database

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
		util.Responses.Error(w, http.StatusBadRequest, "failed to decode body: "+err.Error())
		return
	} else if err, _ := util.ValidateBody(body, []string{"type", "name"}, map[string]map[string]string{
		"type": {"required": "true", "type": "string"},
		"name": {"required": "true", "type": "string"},
	}); err != "" {
		util.Responses.Error(w, http.StatusBadRequest, err)
		return
	}

	// Parse out body by type
	switch strings.ToUpper(body["type"].(string)) {
	case "A":
		if err, _ := util.ValidateBody(body, []string{"host"}, map[string]map[string]string{"host": {"required": "true", "type": "ipv4"}}); err != "" {
			util.Responses.Error(w, http.StatusBadRequest, err)
			return
		} else if err := db.Set.A(body["name"].(string), body["host"].(string)); err != nil {
			util.Responses.Error(w, http.StatusInternalServerError, "failed to write record to database: "+err.Error())
			return
		}
	case "AAAA":
		if err, _ := util.ValidateBody(body, []string{"host"}, map[string]map[string]string{"host": {"required": "true", "type": "ipv6"}}); err != "" {
			util.Responses.Error(w, http.StatusBadRequest, err)
			return
		} else if err := db.Set.AAAA(body["name"].(string), body["host"].(string)); err != nil {
			util.Responses.Error(w, http.StatusInternalServerError, "failed to write record to database: "+err.Error())
			return
		}
	case "CNAME":
		if err, _ := util.ValidateBody(body, []string{"target"}, map[string]map[string]string{"target": {"required": "true", "type": "string"}}); err != "" {
			util.Responses.Error(w, http.StatusBadRequest, err)
			return
		} else if err := db.Set.CNAME(body["name"].(string), body["target"].(string)); err != nil {
			util.Responses.Error(w, http.StatusInternalServerError, "failed to write record to database: "+err.Error())
			return
		}
	case "MX":
		if err, _ := util.ValidateBody(body, []string{"priority", "host"}, map[string]map[string]string{
			"priority": {"type": "uint16", "required": "true"},
			"host": {"type": "string", "required": "true"},
		}); err != "" {
			util.Responses.Error(w, http.StatusBadRequest, err)
			return
		} else if err := db.Set.MX(body["name"].(string), uint16(body["priority"].(float64)), body["host"].(string)); err != nil {
			util.Responses.Error(w, http.StatusBadRequest, "failed to write record to database: "+err.Error())
			return
		}
	case "LOC":
		if err, _ := util.ValidateBody(body, []string{"version", "size", "horizontal-precision", "vertical-precision", "latitude", "longitude", "altitude"}, map[string]map[string]string{
			"version": {"type": "uint8", "required": "true"},
			"size": {"type": "uint8", "required": "true"},
			"horizontal-precision": {"type": "uint8", "required": "true"},
			"vertical-precision": {"type": "uint8", "required": "true"},
			"latitude": {"type": "uint32", "required": "true"},
			"longitude": {"type": "uint32", "required": "true"},
			"altitude": {"type": "uint32", "required": "true"},
		}); err != "" {
			util.Responses.Error(w, http.StatusBadRequest, err)
		} else if err := db.Set.LOC(body["name"].(string), uint8(body["version"].(float64)), uint8(body["size"].(float64)), uint8(body["horizontal-precision"].(float64)), uint8(body["vertical-precision"].(float64)), uint32(body["latitude"].(float64)), uint32(body["longitude"].(float64)), uint32(body["altitude"].(float64))); err != nil {
			util.Responses.Error(w, http.StatusInternalServerError, "failed to write record to database: "+err.Error())
			return
		}
	case "SRV":
		if err, _ := util.ValidateBody(body, []string{"priority", "weight", "port", "target"}, map[string]map[string]string{
			"priority": {"type": "uint16", "required": "true"},
			"weight": {"type": "uint16", "required": "true"},
			"port": {"type": "uint16", "required": "true"},
			"target": {"type": "string", "required": "true"},
		}); err != "" {
			util.Responses.Error(w, http.StatusBadRequest, err)
		} else if err := db.Set.SRV(body["name"].(string), uint16(body["priority"].(float64)), uint16(body["weight"].(float64)), uint16(body["port"].(float64)), body["target"].(string)); err != nil {
			util.Responses.Error(w, http.StatusInternalServerError, "failed to write record to database: "+err.Error())
			return
		}
	case "SPF":
		if err, _ := util.ValidateBody(body, []string{"text"}, map[string]map[string]string{"text": {"type": "stringarray", "required": "true"}}); err != "" {
			util.Responses.Error(w, http.StatusBadRequest, err)
			return
		}
		text, _ := util.ConvertArrayToString(body["text"].([]interface{}))
		if err := db.Set.SPF(body["name"].(string), text); err != nil {
			util.Responses.Error(w, http.StatusBadRequest, "failed to write record to database: "+err.Error())
			return
		}
	case "TXT":
		if err, _ := util.ValidateBody(body, []string{"text"}, map[string]map[string]string{"text": {"type": "stringarray", "required": "true"}}); err != "" {
			util.Responses.Error(w, http.StatusBadRequest, err)
			return
		}
		text, _ := util.ConvertArrayToString(body["text"].([]interface{}))
		if err := db.Set.TXT(body["name"].(string), text); err != nil {
			util.Responses.Error(w, http.StatusBadRequest, "failed to write record to database: "+err.Error())
			return
		}
	case "NS":
		if err, _ := util.ValidateBody(body, []string{"nameserver"}, map[string]map[string]string{"nameserver": {"type": "string", "required": "true"}}); err != "" {
			util.Responses.Error(w, http.StatusBadRequest, err)
			return
		} else if err := db.Set.NS(body["name"].(string), body["nameserver"].(string)); err != nil {
			util.Responses.Error(w, http.StatusInternalServerError, "failed to write record to database: "+err.Error())
			return
		}
	case "CAA":
		if err, _ := util.ValidateBody(body, []string{"content", "tag"}, map[string]map[string]string{
			"tag": {"type": "string", "required": "true"},
			"content": {"type": "string", "required": "true"},
		}); err != "" {
			util.Responses.Error(w, http.StatusBadRequest, err)
			return
		} else if err := db.Set.CAA(body["name"].(string), body["tag"].(string), body["content"].(string)); err != nil {
			util.Responses.Error(w, http.StatusInternalServerError, "failed to write record to database: "+err.Error())
			return
		}
	case "PTR":
		if err, _ := util.ValidateBody(body, []string{"domain"}, map[string]map[string]string{"domain": {"type": "string", "required": "true"}}); err != "" {
			util.Responses.Error(w, http.StatusBadRequest, err)
			return
		} else if err := db.Set.PTR(body["name"].(string), body["domain"].(string)); err != nil {
			util.Responses.Error(w, http.StatusInternalServerError, "failed to write record to database: "+err.Error())
			return
		}
	case "CERT":
		if err, _ := util.ValidateBody(body, []string{"c-type", "key-tag", "algorithm", "certificate"}, map[string]map[string]string{
			"c-type": {"type": "uint16", "required": "true"},
			"key-tag": {"type": "uint16", "required": "true"},
			"algorithm": {"type": "uint8", "required": "true"},
			"certificate": {"type": "string", "required": "true"},
		}); err != "" {
			util.Responses.Error(w, http.StatusBadRequest, err)
			return
		} else if err := db.Set.CERT(body["name"].(string), uint16(body["c-type"].(float64)), uint16(body["key-tag"].(float64)), uint8(body["algorithm"].(float64)), body["certificate"].(string)); err != nil {
			util.Responses.Error(w, http.StatusInternalServerError, "failed to write record to database: "+err.Error())
			return
		}
	case "DNSKEY":
		if err, _ := util.ValidateBody(body, []string{"flags", "protocol", "algorithm", "public-key"}, map[string]map[string]string{
			"flags": {"type": "uint16", "required": "true"},
			"protocol": {"type": "uint8", "required": "true"},
			"algorithm": {"type": "uint8", "required": "true"},
			"public-key": {"type": "string", "required": "true"},
		}); err != "" {
			util.Responses.Error(w, http.StatusBadRequest, err)
			return
		} else if err := db.Set.DNSKEY(body["name"].(string), uint16(body["flags"].(float64)), uint8(body["protocol"].(float64)), uint8(body["algorithm"].(float64)), body["public-key"].(string)); err != nil {
			util.Responses.Error(w, http.StatusInternalServerError, "failed to write record to database: "+err.Error())
			return
		}
	case "DS":
		if err, _ := util.ValidateBody(body, []string{"key-tag", "algorithm", "digest-type", "digest"}, map[string]map[string]string{
			"key-tag": {"type": "uint16", "required": "true"},
			"algorithm": {"type": "uint8", "required": "true"},
			"digest-type": {"type": "uint8", "required": "true"},
			"digest": {"type": "string", "required": "true"},
		}); err != "" {
			util.Responses.Error(w, http.StatusBadRequest, err)
			return
		} else if err := db.Set.DS(body["name"].(string), uint16(body["key-tag"].(float64)), uint8(body["algorithm"].(float64)), uint8(body["digest-type"].(float64)), body["digest"].(string)); err != nil {
			util.Responses.Error(w, http.StatusInternalServerError, "failed to write record to database: "+err.Error())
			return
		}
	case "NAPTR":
		if err, _ := util.ValidateBody(body, []string{"order", "preference", "flags", "service", "regexp", "replacement"}, map[string]map[string]string{
			"order": {"type": "uint16", "required": "true"},
			"preference": {"type": "uint16", "required": "true"},
			"flags": {"type": "string", "required": "true"},
			"service": {"type": "string", "required": "true"},
			"regexp": {"type": "string", "required": "true"},
			"replacement": {"type": "string", "required": "true"},
		}); err != "" {
			util.Responses.Error(w, http.StatusBadRequest, err)
			return
		} else if err := db.Set.NAPTR(body["name"].(string), uint16(body["order"].(float64)), uint16(body["preference"].(float64)), body["flags"].(string), body["service"].(string), body["regexp"].(string), body["replacement"].(string)); err != nil {
			util.Responses.Error(w, http.StatusInternalServerError, "failed to write record to database: "+err.Error())
			return
		}
	case "SMIMEA":
		if err, _ := util.ValidateBody(body, []string{"usage", "selector", "matching-type", "certificate"}, map[string]map[string]string{
			"usage": {"type": "uint8", "required": "true"},
			"selector": {"type": "uint8", "required": "true"},
			"matching-type": {"type": "uint8", "required": "true"},
			"certificate": {"type": "string", "required": "true"},
		}); err != "" {
			util.Responses.Error(w, http.StatusBadRequest, err)
			return
		} else if err := db.Set.SMIMEA(body["name"].(string), uint8(body["usage"].(float64)), uint8(body["selector"].(float64)), uint8(body["matching-type"].(float64)), body["certificate"].(string)); err != nil {
			util.Responses.Error(w, http.StatusInternalServerError, "failed to write record to database: "+err.Error())
			return
		}
	case "SSHFP":
		if err, _ := util.ValidateBody(body, []string{"algorithm", "s-type", "fingerprint"}, map[string]map[string]string{
			"algorithm": {"type": "uint8", "required": "true"},
			"s-type": {"type": "uint8", "required": "true"},
			"fingerprint": {"type": "string", "required": "true"},
		}); err != "" {
			util.Responses.Error(w, http.StatusBadRequest, err)
			return
		} else if err := db.Set.SSHFP(body["name"].(string), uint8(body["algorithm"].(float64)), uint8(body["s-type"].(float64)), body["fingerprint"].(string)); err != nil {
			util.Responses.Error(w, http.StatusInternalServerError, "failed to write record to database: "+err.Error())
			return
		}
	case "TLSA":
		if err, _ := util.ValidateBody(body, []string{"usage", "selector", "matching-type", "certificate"}, map[string]map[string]string{
			"usage": {"type": "uint8", "required": "true"},
			"selector": {"type": "uint8", "required": "true"},
			"matching-type": {"type": "uint8", "required": "true"},
			"certificate": {"type": "string", "required": "true"},
		}); err != "" {
			util.Responses.Error(w, http.StatusBadRequest, err)
			return
		} else if err := db.Set.TLSA(body["name"].(string), uint8(body["usage"].(float64)), uint8(body["selector"].(float64)), uint8(body["matching-type"].(float64)), body["certificate"].(string)); err != nil {
			util.Responses.Error(w, http.StatusInternalServerError, "failed to write record to database: "+err.Error())
			return
		}
	case "URI":
		if err, _ := util.ValidateBody(body, []string{"priority", "weight", "target"}, map[string]map[string]string{
			"priority": {"type": "uint16", "required": "true"},
			"weight": {"type": "uint16", "required": "true"},
			"target": {"type": "string", "required": "true"},
		}); err != "" {
			util.Responses.Error(w, http.StatusBadRequest, err)
			return
		} else if err := db.Set.URI(body["name"].(string), uint16(body["priority"].(float64)), uint16(body["weight"].(float64)), body["target"].(string)); err != nil {
			util.Responses.Error(w, http.StatusInternalServerError, "failed to write record to database: "+err.Error())
			return
		}
	default:
		util.Responses.Error(w, http.StatusBadRequest, "field 'type' must be on of: A, AAAA, CNAME, MX, LOC, SRV, SPF, TXT, NS, CAA, PTR, CERT, DNSKEY, DS, NAPTR, SMIMEA, SSHFP, TLSA, URI")
		return
	}

	util.Responses.Success(w)
}
