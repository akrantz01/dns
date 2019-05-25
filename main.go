package main

import (
	"github.com/miekg/dns"
	"log"
	"net"
)

var records = map[string]string{
	"ipv4.test.": "127.0.0.1",
	"ipv6.test.": "::1",
}

type handler struct {}
func (h *handler) ServeDNS(w dns.ResponseWriter, m *dns.Msg) {
	r := dns.Msg{}
	r.SetReply(m)
	r.Authoritative = true

	for _, q := range r.Question {
		hdr := dns.RR_Header{Name: q.Name, Rrtype: q.Qtype, Class: q.Qclass}

		switch q.Qtype {
		case dns.TypeA:
			ip, ok := records[q.Name]
			if ok {
				r.Answer = append(r.Answer, &dns.A{Hdr: hdr, A: net.ParseIP(ip)})
			}
		case dns.TypeAAAA:
			ip, ok := records[q.Name]
			if ok {
				r.Answer = append(r.Answer, &dns.AAAA{Hdr: hdr, AAAA: net.ParseIP(ip)})
			}
		default:
			r.Rcode = dns.RcodeNameError
		}
	}

	if len(r.Answer) == 0 {
		r.Rcode = dns.RcodeNameError
	}

	if err := w.WriteMsg(&r); err != nil {
		log.Printf("Unable to send response: %v", err)
	}
}

func main() {
	server := &dns.Server{Addr: "127.0.0.1:1053", Net: "udp"}
	server.Handler = &handler{}
	if err := server.ListenAndServe(); err != nil {
		log.Fatalf("Failed to set udp listener: %v", err)
	}
}
