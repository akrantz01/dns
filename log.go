package main

import (
	"github.com/miekg/dns"
	"k8s.io/utils/strings"
	"log"
)

func logQuestion(q dns.Question) {
	var qtype string
	switch q.Qtype {
	// DS NAPTR SMIMEA SSHFP TLSA URI
	case dns.TypeA: qtype = "A"
	case dns.TypeAAAA: qtype = "AAAA"
	case dns.TypeCNAME: qtype = "CNAME"
	case dns.TypeMX: qtype = "MX"
	case dns.TypeLOC: qtype = "LOC"
	case dns.TypeSRV: qtype = "SRV"
	case dns.TypeSPF: qtype = "SPF"
	case dns.TypeTXT: qtype = "TXT"
	case dns.TypeNXT: qtype = "NXT"
	case dns.TypeCAA: qtype = "CAA"
	case dns.TypePTR: qtype = "PTR"
	case dns.TypeCERT: qtype = "CERT"
	case dns.TypeDNSKEY: qtype = "DNSKEY"
	case dns.TypeDS: qtype = "DS"
	case dns.TypeNAPTR: qtype = "NAPTR"
	case dns.TypeSMIMEA: qtype = "SMIMEA"
	case dns.TypeSSHFP: qtype = "SSHFP"
	case dns.TypeTLSA: qtype = "TLSA"
	case dns.TypeURI: qtype = "URI"
	}

	log.Printf("%s in %s", strings.ShortenString(q.Name, len(q.Name) - 1), qtype)
}
