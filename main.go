package main

import (
	"flag"
	"github.com/akrantz01/krantz.dev/dns/db"
	"github.com/akrantz01/krantz.dev/dns/routes"
	"github.com/akrantz01/krantz.dev/dns/util"
	"github.com/gorilla/handlers"
	"github.com/miekg/dns"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
	bolt "go.etcd.io/bbolt"
	"log"
	"net/http"
	"os"
	"time"
)

var database *bolt.DB

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
			ip := db.GetARecord(database, q.Name)
			if ip != nil {
				r.Answer = append(r.Answer, &dns.A{Hdr: hdr, A: ip})
			}
		case dns.TypeAAAA:
			ip := db.GetAAAARecord(database, q.Name)
			if ip != nil {
				r.Answer = append(r.Answer, &dns.AAAA{Hdr: hdr, AAAA: ip})
			}
		case dns.TypeCNAME:
			target := db.GetCNAMERecord(database, q.Name)
			if target != "" {
				r.Answer = append(r.Answer, &dns.CNAME{Hdr: hdr, Target: target})
			}
		case dns.TypeMX:
			host, priority := db.GetMXRecord(database, q.Name)
			if host != "" {
				r.Answer = append(r.Answer, &dns.MX{Hdr: hdr, Preference: priority, Mx: host})
			}
		case dns.TypeLOC:
			vers, siz, hor, ver, lat, lon, alt := db.GetLOCRecord(database, q.Name)
			if lat != 0 && lon != 0 {
				r.Answer = append(r.Answer, &dns.LOC{Hdr: hdr, Version: vers, Size: siz, HorizPre: hor, VertPre: ver, Latitude: lat, Longitude: lon, Altitude: alt})
			}
		case dns.TypeSRV:
			priority, weight, port, target := db.GetSRVRecord(database, q.Name)
			if target != "" {
				r.Answer = append(r.Answer, &dns.SRV{Hdr: hdr, Priority: priority, Weight: weight, Port: port, Target: target})
			}
		case dns.TypeSPF:
			txt := db.GetSPFRecord(database, q.Name)
			if len(txt) != 0 {
				r.Answer = append(r.Answer, &dns.SPF{Hdr: hdr, Txt: txt})
			}
		case dns.TypeTXT:
			content := db.GetTXTRecord(database, q.Name)
			if len(content) != 0 {
				r.Answer = append(r.Answer, &dns.TXT{Hdr: hdr, Txt: content})
			}
		case dns.TypeNS:
			nameserver := db.GetNSRecord(database, q.Name)
			if nameserver != "" {
				r.Answer = append(r.Answer, &dns.NS{Hdr: hdr, Ns: nameserver})
			}
		case dns.TypeCAA:
			dflag, tag, content := db.GetCAARecord(database, q.Name)
			if tag != "" && content != "" {
				r.Answer = append(r.Answer, &dns.CAA{Hdr: hdr, Flag: dflag, Tag: tag, Value: content})
			}
		case dns.TypePTR:
			ptr := db.GetPTRRecord(database, q.Name)
			if ptr != "" {
				r.Answer = append(r.Answer, &dns.PTR{Hdr: hdr, Ptr: ptr})
			}
		case dns.TypeCERT:
			tpe, tag, algo, cert := db.GetCERTRecord(database, q.Name)
			if cert != "" {
				r.Answer = append(r.Answer, &dns.CERT{Hdr: hdr, Type: tpe, KeyTag: tag, Algorithm: algo, Certificate: cert})
			}
		case dns.TypeDNSKEY:
			flags, proto, algo, pub := db.GetDNSKEYRecord(database, q.Name)
			if pub != "" {
				r.Answer = append(r.Answer, &dns.DNSKEY{Hdr: hdr, Flags: flags, Protocol: proto, Algorithm: algo, PublicKey: pub})
			}
		case dns.TypeDS:
			ktag, algo, dtype, digest := db.GetDSRecord(database, q.Name)
			if digest != "" {
				r.Answer = append(r.Answer, &dns.DS{Hdr: hdr, KeyTag: ktag, Algorithm: algo, DigestType: dtype, Digest: digest})
			}
		case dns.TypeNAPTR:
			ord, pref, dflag, serv, reg, rep := db.GetNAPTRRecord(database, q.Name)
			if serv != "" && reg != "" && rep != "" {
				r.Answer = append(r.Answer, &dns.NAPTR{Hdr: hdr, Order: ord, Preference: pref, Flags: dflag, Service: serv, Regexp: reg, Replacement: rep})
			}
		case dns.TypeSMIMEA:
			usage, sel, match, cert := db.GetSMIMEARecord(database, q.Name)
			if cert != "" {
				r.Answer = append(r.Answer, &dns.SMIMEA{Hdr: hdr, Usage: usage, Selector: sel, MatchingType: match, Certificate: cert})
			}
		case dns.TypeSSHFP:
			algo, tpe, fingerprint := db.GetSSHFPRecord(database, q.Name)
			if fingerprint != "" {
				r.Answer = append(r.Answer, &dns.SSHFP{Hdr: hdr, Algorithm: algo, Type: tpe, FingerPrint: fingerprint})
			}
		case dns.TypeTLSA:
			usg, sel, mat, cert := db.GetTLSARecord(database, q.Name)
			if cert != "" {
				r.Answer = append(r.Answer, &dns.TLSA{Hdr: hdr, Usage: usg, Selector: sel, MatchingType: mat, Certificate: cert})
			}
		case dns.TypeURI:
			pri, wei, tar := db.GetURIRecord(database, q.Name)
			if tar != "" {
				r.Answer = append(r.Answer, &dns.URI{Hdr: hdr, Priority: pri, Weight: wei, Target: tar})
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

	util.LogResponse(w, r, start)
}

func main() {
	// Configuration setup
	viper.SetConfigName("config")
	viper.AddConfigPath(".")
	viper.AddConfigPath("$HOME")

	// Setup environment variables
	viper.SetEnvPrefix("dns")
	if err := viper.BindEnv("host", "port", "database", "tcp", "udp"); err != nil { log.Fatalf("Failed to setup environment variables: %v", err) }

	// Setup command line
	flag.String("host", "127.0.0.1", "IP address to run on")
	flag.Int("port", 53, "Port to listen on")
	flag.String("database", "./records.db", "Database file to use")
	flag.Bool("no-tcp", false, "Disable listening on TCP")
	flag.Bool("no-udp", false, "Disable listening on UDP")
	pflag.CommandLine.AddGoFlagSet(flag.CommandLine)
	pflag.Parse()
	if err := viper.BindPFlags(pflag.CommandLine); err != nil { log.Fatalf("Failed to setup command line arguments: %v", err) }

	// Set configuration defaults
	viper.SetDefault("host", "127.0.0.1")
	viper.SetDefault("port", 53)
	viper.SetDefault("database", "./records.db")
	viper.SetDefault("no-tcp", false)
	viper.SetDefault("no-udp", false)

	// Parse configuration
	if err := viper.ReadInConfig(); err != nil {
		switch err.(type) {
		case viper.ConfigFileNotFoundError:
			break
		default:
			log.Fatalf("Failed to read server configuration: %v", err)
		}
	}

	// Open database
	var err error
	database, err = bolt.Open(viper.GetString("database"), 0666, nil)
	if err != nil {
		log.Fatalf("Failed to open database: %v", err)
	}
	defer func() { if err := database.Close(); err != nil { log.Fatalf("Failed to close database: %v", err) }}()

	// Setup database structure
	if err := db.SetupDB(database); err != nil {
		log.Fatalf("Failed setting up database structure: %v", err)
	}

	// Check config is valid
	if viper.GetBool("no-tcp") && viper.GetBool("no-udp") { log.Fatalf("Invalid configuration: tcp and/or udp must be enabled, got both as disabled") }

	// Handle TCP connections
	tcpErr := make(chan error)
	go func() {
		if viper.GetBool("no-tcp") { return }
		tcp := &dns.Server{Addr: viper.GetString("host") + ":" + viper.GetString("port"), Net: "tcp"}
		tcp.Handler = &handler{}

		if err := tcp.ListenAndServe(); err != nil { tcpErr <- err }
	}()

	// Handle UDP connections
	udpErr := make(chan error)
	go func() {
		if viper.GetBool("no-udp") { return }
		udp := &dns.Server{Addr: viper.GetString("host") + ":" + viper.GetString("port"), Net: "udp"}
		udp.Handler = &handler{}

		if err := udp.ListenAndServe(); err != nil { udpErr <- err }
	}()

	// Handle REST API
	httpErr := make(chan error)
	go func() {
		http.Handle("/records", handlers.LoggingHandler(os.Stdout, http.HandlerFunc(routes.AllRecordsHandler(database))))
		http.Handle("/records/", handlers.LoggingHandler(os.Stdout, http.HandlerFunc(routes.SingleRecordHandler("/records/", database))))
		if err := http.ListenAndServe("127.0.0.1:8080", nil); err != nil { httpErr <- err }
	}()

	// Assemble log
	var protocols string
	if !viper.GetBool("no-udp") && !viper.GetBool("no-tcp") {
		protocols = "TCP and UDP"
	} else if viper.GetBool("no-udp") {
		protocols = "TCP"
	} else {
		protocols = "UDP"
	}
	log.Printf("DNS server listening on %s:%s with %s...", viper.GetString("host"), viper.GetString("port"), protocols)
	log.Printf("HTTP server listening on 127.0.0.1:8080...")

	// Watch for errors
	select {
	case err := <- tcpErr:
		log.Fatalf("DNS failed to listen on 127.0.0.1:1052 with TCP: %v\n", err)
	case err := <- udpErr:
		log.Fatalf("DNS failed to listen on 127.0.0.1:1053 with UDP: %v\n", err)
	case err := <- httpErr:
		log.Fatalf("API failed to listen on 127.0.0.1:8080: %v\n", err)
	}
}
