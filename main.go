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
		logQuestion(q)
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
	// Handle TCP connections
	tcpErr := make(chan error)
	go func() {
		tcp := &dns.Server{Addr: "127.0.0.1:1053", Net: "tcp"}
		tcp.Handler = &handler{}

		if err := tcp.ListenAndServe(); err != nil { tcpErr <- err }
	}()

	// Handle UDP connections
	udpErr := make(chan error)
	go func() {
		udp := &dns.Server{Addr: "127.0.0.1:1053", Net: "udp"}
		udp.Handler = &handler{}

		if err := udp.ListenAndServe(); err != nil { udpErr <- err }
	}()

	log.Println("Listening on 127.0.0.1:1053 with TCP and UDP...")

	select {
	case err := <- tcpErr:
		log.Fatalf("Failed to listen on 127.0.0.1:1052 with TCP: %v\n", err)
	case err := <- udpErr:
		log.Fatalf("Failed to listen on 127.0.0.1:1053 with UDP: %v\n", err)
	}
}
