package main

import (
	"github.com/miekg/dns"
	bolt "go.etcd.io/bbolt"
	"log"
	"time"
)

var db *bolt.DB

type handler struct {}
func (h *handler) ServeDNS(w dns.ResponseWriter, m *dns.Msg) {
	start := time.Now()

	r := new(dns.Msg)
	r.SetReply(m)
	r.Authoritative = true

	for _, q := range r.Question {
		hdr := dns.RR_Header{Name: q.Name, Rrtype: q.Qtype, Class: q.Qclass}

		switch q.Qtype {
		case dns.TypeA:
			ip, _ := getARecord(db, q.Name)
			if ip != nil {
				r.Answer = append(r.Answer, &dns.A{Hdr: hdr, A: ip})
			}
		case dns.TypeAAAA:
			ip, _ := getAAAARecord(db, q.Name)
			if ip != nil {
				r.Answer = append(r.Answer, &dns.AAAA{Hdr: hdr, AAAA: ip})
			}
		case dns.TypeCNAME:
			target, _ := getCNAMERecord(db, q.Name)
			if target != "" {
				r.Answer = append(r.Answer, &dns.CNAME{Hdr: hdr, Target: target})
			}
		case dns.TypeMX:
			host, priority, _ := getMXRecord(db, q.Name)
			if host != "" && priority != 0 {
				r.Answer = append(r.Answer, &dns.MX{Hdr: hdr, Preference: priority, Mx: host})
			}
		case dns.TypeLOC:
			vers, siz, hor, ver, lat, lon, alt, _ := getLOCRecord(db, q.Name)
			if vers != 0 && siz != 0 && hor != 0 && ver != 0 && lat != 0 && lon != 0 && alt != 0 {
				r.Answer = append(r.Answer, &dns.LOC{Hdr: hdr, Version: vers, Size: siz, HorizPre: hor, VertPre: ver, Latitude: lat, Longitude: lon, Altitude: alt})
			}
		case dns.TypeSRV:
			priority, weight, port, target, _ := getSRVRecord(db, q.Name)
			if priority != 0 && weight != 0 && port != 0 && target != "" {
				r.Answer = append(r.Answer, &dns.SRV{Hdr: hdr, Priority: priority, Weight: weight, Port: port, Target: target})
			}
		case dns.TypeSPF:
			txt, _ := getSPFRecord(db, q.Name)
			if len(txt) != 0 {
				r.Answer = append(r.Answer, &dns.SPF{Hdr: hdr, Txt: txt})
			}
		case dns.TypeTXT:
			content, _ := getTXTRecord(db, q.Name)
			if len(content) != 0 {
				r.Answer = append(r.Answer, &dns.TXT{Hdr: hdr, Txt: content})
			}
		default:
			r.Rcode = dns.RcodeNameError
		}
	}

	if len(r.Answer) == 0 {
		r.Rcode = dns.RcodeNameError
	}

	if err := w.WriteMsg(r); err != nil {
		log.Printf("Unable to send response: %v", err)
	}

	logResponse(w, r, start)
}

func main() {
	// Open database
	var err error
	db, err = bolt.Open("records.db", 0666, nil)
	if err != nil {
		log.Fatalf("Failed to open database: %v", err)
	}
	defer func() { if err := db.Close(); err != nil { log.Fatalf("Failed to close database: %v", err) }}()

	// Setup database structure
	if err := setupDB(db); err != nil {
		log.Fatalf("Failed setting up database structure: %v", err)
	}

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

	// Watch for errors
	select {
	case err := <- tcpErr:
		log.Fatalf("Failed to listen on 127.0.0.1:1052 with TCP: %v\n", err)
	case err := <- udpErr:
		log.Fatalf("Failed to listen on 127.0.0.1:1053 with UDP: %v\n", err)
	}
}
