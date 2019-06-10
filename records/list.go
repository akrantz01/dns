package records

import (
	"github.com/akrantz01/krantz.dev/dns/db"
	"github.com/akrantz01/krantz.dev/dns/util"
	bolt "go.etcd.io/bbolt"
	"net/http"
	"strings"
)

// Handle the listing of all records
func list(w http.ResponseWriter, r *http.Request, database *bolt.DB) {
	if r.Method != "GET" {
		util.Responses.Error(w, http.StatusMethodNotAllowed, "method not allowed")
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

	// List all records of a type if query parameter given
	if len(r.URL.RawQuery) != 0 {
		if _, ok := r.URL.Query()["type"]; !ok {
			util.Responses.Error(w, http.StatusBadRequest, "query parameter 'type' is required for type filtering")
			return
		}

		var rawRecords []map[string]string
		for _, record := range r.URL.Query()["type"] {
			if err := database.View(func(tx *bolt.Tx) error {
				return tx.Bucket([]byte(record)).ForEach(func(k, v []byte) error {
					rawRecords = append(rawRecords, map[string]string{"name": strings.Split(string(k), "*")[0], "type": record})
					return nil
				})
			}); err != nil {
				util.Responses.Error(w, http.StatusInternalServerError, "failed to retrieve all records: "+err.Error())
				return
			}
		}

		// Remove duplicates from array
		encountered := map[string]bool{}
		var records []map[string]string
		for _, v := range rawRecords {
			if _, ok := encountered[v["name"]]; !ok {
				encountered[v["name"]] = true
				records = append(records, v)
			}
		}

		// Return empty array if none
		if records == nil {
			util.Responses.SuccessWithData(w, []string{})
			return
		}
		util.Responses.SuccessWithData(w, records)
		return
	}

	var rawRecords []map[string]string
	for _, record := range []string{"A", "AAAA", "CNAME", "MX", "LOC", "SRV", "SPF", "TXT", "NS", "CAA", "PTR", "CERT", "DNSKEY", "DS", "NAPTR", "SMIMEA", "SSHFP", "TLSA", "URI"} {
		if err := database.View(func(tx *bolt.Tx) error {
			return tx.Bucket([]byte(record)).ForEach(func(k, v []byte) error {
				rawRecords = append(rawRecords, map[string]string{"name": strings.Split(string(k), "*")[0], "type": record})
				return nil
			})
		}); err != nil {
			util.Responses.Error(w, http.StatusInternalServerError, "failed to retrieve all records: "+err.Error())
		}
	}

	// Remove duplicates from array
	encountered := map[string]bool{}
	var records []map[string]string
	for _, v := range rawRecords {
		if _, ok := encountered[v["name"]]; !ok {
			encountered[v["name"]] = true
			records = append(records, v)
		}
	}

	// Return empty array if none
	if records == nil {
		util.Responses.SuccessWithData(w, []string{})
		return
	}
	util.Responses.SuccessWithData(w, records)
}
