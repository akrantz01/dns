package records

import (
	"github.com/akrantz01/krantz.dev/dns/util"
	bolt "go.etcd.io/bbolt"
	"net/http"
	"strings"
)

// Handle the listing of all records
func list(w http.ResponseWriter, r *http.Request, db *bolt.DB) {
	if r.Method != "GET" {
		util.Responses.Error(w, http.StatusMethodNotAllowed, "method not allowed")
		return
	} else if len(r.URL.RawQuery) != 0 {
		if _, ok := r.URL.Query()["type"]; !ok {
			util.Responses.Error(w, http.StatusBadRequest, "query parameter 'type' is required for type filtering")
			return
		}

		records := make(map[string][]string)
		for _, record := range r.URL.Query()["type"] {
			records[record] = []string{}
			if err := db.View(func(tx *bolt.Tx) error {
				return tx.Bucket([]byte(record)).ForEach(func(k, v []byte) error {
					records[record] = append(records[record], strings.Split(string(k), "*")[0])
					return nil
				})
			}); err != nil {
				util.Responses.Error(w, http.StatusInternalServerError, "failed to retrieve all records: " + err.Error())
				return
			}
			records[record] = util.RemoveDuplicates(records[record])
		}

		util.Responses.SuccessWithData(w, records)
		return
	}

	records := map[string][]string{
		"A":      {},
		"AAAA":   {},
		"CNAME":  {},
		"MX":     {},
		"LOC":    {},
		"SRV":    {},
		"SPF":    {},
		"TXT":    {},
		"NS":     {},
		"CAA":    {},
		"PTR":    {},
		"CERT":   {},
		"DNSKEY": {},
		"DS":     {},
		"NAPTR":  {},
		"SMIMEA": {},
		"SSHFP":  {},
		"TLSA":   {},
		"URI":    {},
	}
	for record := range records {
		if err := db.View(func(tx *bolt.Tx) error {
			return tx.Bucket([]byte(record)).ForEach(func(k, v []byte) error {
				records[record] = append(records[record], strings.Split(string(k), "*")[0])
				return nil
			})
		}); err != nil {
			util.Responses.Error(w, http.StatusInternalServerError, "failed to retrieve all records: " + err.Error())
			return
		}
		records[record] = util.RemoveDuplicates(records[record])
	}

	util.Responses.SuccessWithData(w, records)
}
