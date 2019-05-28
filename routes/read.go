package routes

import (
	"github.com/akrantz01/krantz.dev/dns/db"
	"github.com/akrantz01/krantz.dev/dns/util"
	bolt "go.etcd.io/bbolt"
	"net/http"
	"strings"
)

func read(w http.ResponseWriter, r *http.Request, path string, database *bolt.DB) {
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
	response := make(map[string]interface{})

	switch r.URL.Query().Get("type") {
	case "A":
		ip := db.GetARecord(database, record)
		response["host"] = ip.String()
	case "AAAA":
		response["host"] = db.GetAAAARecord(database, record).String()
	case "CNAME":
		response["domain"] = db.GetCNAMERecord(database, record)
	case "MX":
		server, priority := db.GetMXRecord(database, record)
		response["server"] = server
		response["priority"] = priority
	case "LOC":
		version, size, horizontal, vertical, latitude, longitude, altitude := db.GetLOCRecord(database, record)
		response["version"] = version
		response["size"] = size
		response["horizontal-precision"] = horizontal
		response["vertical-precision"] = vertical
		response["latitude"] = latitude
		response["longitude"] = longitude
		response["altitude"] = altitude
	case "SRV":
		priority, weight, port, target := db.GetSRVRecord(database, record)
		response["priority"] = priority
		response["weight"] = weight
		response["port"] = port
		response["target"] = target
	case "SPF":
		response["policy"] = db.GetSPFRecord(database, record)
	case "TXT":
		response["text"] = db.GetTXTRecord(database, record)
	case "NS":
		response["nameserver"] = db.GetNSRecord(database, record)
	case "CAA":
		flag, tag, value := db.GetCAARecord(database, record)
		response["flag"] = flag
		response["tag"] = tag
		response["value"] = value
	case "PTR":
		response["domain"] = db.GetPTRRecord(database, record)
	case "CERT":
		tpe, keyTag, algorithm, certificate := db.GetCERTRecord(database, record)
		response["type"] = tpe
		response["key-tag"] = keyTag
		response["algorithm"] = algorithm
		response["certificate"] = certificate
	case "DNSKEY":
		flags, protocol, algorithm, publicKey := db.GetDNSKEYRecord(database, record)
		response["flags"] = flags
		response["protocol"] = protocol
		response["algorithm"] = algorithm
		response["public-key"] = publicKey
	case "DS":
		keyTag, algorithm, digestType, digest := db.GetDSRecord(database, record)
		response["key-tag"] = keyTag
		response["algorithm"] = algorithm
		response["digest-type"] = digestType
		response["digest"] = digest
	case "NAPTR":
		order, preference, flags, service, regexp, replacement := db.GetNAPTRRecord(database, record)
		response["order"] = order
		response["preference"] = preference
		response["flags"] = flags
		response["service"] = service
		response["regexp"] = regexp
		response["replacement"] = replacement
	case "SMIMEA":
		usage, selector, matchingType, certificate := db.GetSMIMEARecord(database, record)
		response["usage"] = usage
		response["selector"] = selector
		response["matching-type"] = matchingType
		response["certificate"] = certificate
	case "SSHFP":
		algorithm, tpe, fingerprint := db.GetSSHFPRecord(database, record)
		response["algorithm"] = algorithm
		response["type"] = tpe
		response["fingerprint"] = fingerprint
	case "TLSA":
		usage, selector, matchingType, certificate := db.GetTLSARecord(database, record)
		response["usage"] = usage
		response["selector"] = selector
		response["matching-type"] = matchingType
		response["certificate"] = certificate
	case "URI":
		priority, weight, content := db.GetURIRecord(database, record)
		response["priority"] = priority
		response["weight"] = weight
		response["content"] = content
	default:
		util.Responses.Error(w, http.StatusBadRequest, "query parameter 'type' must be on of: A, AAAA, CNAME, MX, LOC, SRV, SPF, TXT, NS, CAA, PTR, CERT, DNSKEY, DS, NAPTR, SMIMEA, SSHFP, TLSA, URI")
		return
	}

	util.Responses.SuccessWithData(w, response)
}
