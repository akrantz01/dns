package records

import (
	"github.com/akrantz01/krantz.dev/dns/db"
	"github.com/akrantz01/krantz.dev/dns/util"
	bolt "go.etcd.io/bbolt"
	"net/http"
	"strings"
)

func deleteRecord(w http.ResponseWriter, r *http.Request, path string, database *bolt.DB) {
	// Set database into operations
	db.Get.Db = database
	db.Set.Db = database
	db.Delete.Db = database

	if r.Method != "DELETE" {
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
	} else if r.Header.Get("Authorization") == "" {
		util.Responses.Error(w, http.StatusUnauthorized, "header 'Authorization' is required")
		return
	}

	// Verify JWT in headers
	_, err := db.TokenFromString(r.Header.Get("Authorization"), database)
	if err != nil {
		util.Responses.Error(w, http.StatusUnauthorized, "failed to authenticate: "+err.Error())
		return
	}

	// Accounts for extra dot and all lowercase in DNS query
	record := strings.ToLower(r.URL.Path[len(path):])

	switch r.URL.Query().Get("type") {
	case "A":
		err = db.Delete.A(record)
	case "AAAA":
		err = db.Delete.AAAA(record)
	case "CNAME":
		err = db.Delete.CNAME(record)
	case "MX":
		err = db.Delete.MX(record)
	case "LOC":
		err = db.Delete.LOC(record)
	case "SRV":
		err = db.Delete.SRV(record)
	case "SPF":
		err = db.Delete.SPF(record)
	case "TXT":
		err = db.Delete.TXT(record)
	case "NS":
		err = db.Delete.NS(record)
	case "CAA":
		err = db.Delete.CAA(record)
	case "PTR":
		err = db.Delete.PTR(record)
	case "CERT":
		err = db.Delete.CERT(record)
	case "DNSKEY":
		err = db.Delete.DNSKEY(record)
	case "DS":
		err = db.Delete.DS(record)
	case "NAPTR":
		err = db.Delete.NAPTR(record)
	case "SMIMEA":
		err = db.Delete.SMIMEA(record)
	case "SSHFP":
		err = db.Delete.SSHFP(record)
	case "TLSA":
		err = db.Delete.TLSA(record)
	case "URI":
		err = db.Delete.URI(record)
	default:
		util.Responses.Error(w, http.StatusBadRequest, "query parameter 'type' must be on of: A, AAAA, CNAME, MX, LOC, SRV, SPF, TXT, NS, CAA, PTR, CERT, DNSKEY, DS, NAPTR, SMIMEA, SSHFP, TLSA, URI")
		return
	}

	if err != nil {
		util.Responses.Error(w, http.StatusBadRequest, err.Error())
		return
	}
	util.Responses.Success(w)
}
