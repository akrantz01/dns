package routes

import (
	"fmt"
	"github.com/akrantz01/krantz.dev/dns/db"
	"github.com/akrantz01/krantz.dev/dns/util"
	bolt "go.etcd.io/bbolt"
	"net/http"
	"strings"
)

func read(w http.ResponseWriter, r *http.Request, path string, database *bolt.DB) {
	// Set database into getter and setter
	db.Get.Db = database
	db.Set.Db = database

	if r.Method != "GET" {
		util.Responses.Error(w, http.StatusMethodNotAllowed, "method not allowed")
		return
	} else if len(r.URL.Path[len(path):]) == 0 {
		util.Responses.Error(w, http.StatusBadRequest, "record must be specified in path")
		return
	} else if len(r.URL.RawQuery) == 0 {
		util.Responses.Error(w, http.StatusBadRequest, "query parameters are required")
		return
	} else if r.URL.Query().Get("type") == "" {
		util.Responses.Error(w, http.StatusBadRequest, "query parameter 'type' is required")
		return
	}

	// Accounts for extra dot and all lowercase in DNS request
	record := strings.ToLower(r.URL.Path[len(path):] + ".")
	var response db.Record

	switch r.URL.Query().Get("type") {
	case "A":
		response = db.Get.A(record)
	case "AAAA":
		response = db.Get.AAAA(record)
	case "CNAME":
		response = db.Get.CNAME(record)
	case "MX":
		response = db.Get.MX(record)
	case "LOC":
		response = db.Get.LOC(record)
	case "SRV":
		response = db.Get.SRV(record)
	case "SPF":
		response = db.Get.SPF(record)
	case "TXT":
		response = db.Get.TXT(record)
	case "NS":
		response = db.Get.NS(record)
	case "CAA":
		response = db.Get.CERT(record)
	case "PTR":
		response = db.Get.PTR(record)
	case "CERT":
		response = db.Get.CERT(record)
	case "DNSKEY":
		response = db.Get.DNSKEY(record)
	case "DS":
		response = db.Get.DS(record)
	case "NAPTR":
		response = db.Get.NAPTR(record)
	case "SMIMEA":
		response = db.Get.SMIMEA(record)
	case "SSHFP":
		response = db.Get.SSHFP(record)
	case "TLSA":
		response = db.Get.TLSA(record)
	case "URI":
		response = db.Get.URI(record)
	default:
		util.Responses.Error(w, http.StatusBadRequest, "query parameter 'type' must be on of: A, AAAA, CNAME, MX, LOC, SRV, SPF, TXT, NS, CAA, PTR, CERT, DNSKEY, DS, NAPTR, SMIMEA, SSHFP, TLSA, URI")
		return
	}

	// I know this is a hack, but its the best thing I could think of
	if fmt.Sprintf("%v", response) == "<nil>" {
		util.Responses.Error(w, http.StatusBadRequest, "record does not exist")
		return
	}

	util.Responses.SuccessWithData(w, response)
}
