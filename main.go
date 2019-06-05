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
			ip := db.Get.A(q.Name)
			if ip != nil {
				r.Answer = append(r.Answer, &dns.A{Hdr: hdr, A: ip})
			}
		case dns.TypeAAAA:
			ip := db.Get.AAAA(q.Name)
			if ip != nil {
				r.Answer = append(r.Answer, &dns.AAAA{Hdr: hdr, AAAA: ip})
			}
		case dns.TypeCNAME:
			target := db.Get.CNAME(q.Name)
			if target != "" {
				r.Answer = append(r.Answer, &dns.CNAME{Hdr: hdr, Target: target})
			}
		case dns.TypeMX:
			host, priority := db.Get.MX(q.Name)
			if host != "" {
				r.Answer = append(r.Answer, &dns.MX{Hdr: hdr, Preference: priority, Mx: host})
			}
		case dns.TypeLOC:
			vers, siz, hor, ver, lat, lon, alt := db.Get.LOC(q.Name)
			if lat != 0 && lon != 0 {
				r.Answer = append(r.Answer, &dns.LOC{Hdr: hdr, Version: vers, Size: siz, HorizPre: hor, VertPre: ver, Latitude: lat, Longitude: lon, Altitude: alt})
			}
		case dns.TypeSRV:
			priority, weight, port, target := db.Get.SRV(q.Name)
			if target != "" {
				r.Answer = append(r.Answer, &dns.SRV{Hdr: hdr, Priority: priority, Weight: weight, Port: port, Target: target})
			}
		case dns.TypeSPF:
			txt := db.Get.SPF(q.Name)
			if len(txt) != 0 {
				r.Answer = append(r.Answer, &dns.SPF{Hdr: hdr, Txt: txt})
			}
		case dns.TypeTXT:
			content := db.Get.TXT(q.Name)
			if len(content) != 0 {
				r.Answer = append(r.Answer, &dns.TXT{Hdr: hdr, Txt: content})
			}
		case dns.TypeNS:
			nameserver := db.Get.NS(q.Name)
			if nameserver != "" {
				r.Answer = append(r.Answer, &dns.NS{Hdr: hdr, Ns: nameserver})
			}
		case dns.TypeCAA:
			dflag, tag, content := db.Get.CAA(q.Name)
			if tag != "" && content != "" {
				r.Answer = append(r.Answer, &dns.CAA{Hdr: hdr, Flag: dflag, Tag: tag, Value: content})
			}
		case dns.TypePTR:
			ptr := db.Get.PTR(q.Name)
			if ptr != "" {
				r.Answer = append(r.Answer, &dns.PTR{Hdr: hdr, Ptr: ptr})
			}
		case dns.TypeCERT:
			tpe, tag, algo, cert := db.Get.CERT(q.Name)
			if cert != "" {
				r.Answer = append(r.Answer, &dns.CERT{Hdr: hdr, Type: tpe, KeyTag: tag, Algorithm: algo, Certificate: cert})
			}
		case dns.TypeDNSKEY:
			flags, proto, algo, pub := db.Get.DNSKEY(q.Name)
			if pub != "" {
				r.Answer = append(r.Answer, &dns.DNSKEY{Hdr: hdr, Flags: flags, Protocol: proto, Algorithm: algo, PublicKey: pub})
			}
		case dns.TypeDS:
			ktag, algo, dtype, digest := db.Get.DS(q.Name)
			if digest != "" {
				r.Answer = append(r.Answer, &dns.DS{Hdr: hdr, KeyTag: ktag, Algorithm: algo, DigestType: dtype, Digest: digest})
			}
		case dns.TypeNAPTR:
			ord, pref, dflag, serv, reg, rep := db.Get.NAPTR(q.Name)
			if serv != "" && reg != "" && rep != "" {
				r.Answer = append(r.Answer, &dns.NAPTR{Hdr: hdr, Order: ord, Preference: pref, Flags: dflag, Service: serv, Regexp: reg, Replacement: rep})
			}
		case dns.TypeSMIMEA:
			usage, sel, match, cert := db.Get.SMIMEA(q.Name)
			if cert != "" {
				r.Answer = append(r.Answer, &dns.SMIMEA{Hdr: hdr, Usage: usage, Selector: sel, MatchingType: match, Certificate: cert})
			}
		case dns.TypeSSHFP:
			algo, tpe, fingerprint := db.Get.SSHFP(q.Name)
			if fingerprint != "" {
				r.Answer = append(r.Answer, &dns.SSHFP{Hdr: hdr, Algorithm: algo, Type: tpe, FingerPrint: fingerprint})
			}
		case dns.TypeTLSA:
			usg, sel, mat, cert := db.Get.TLSA(q.Name)
			if cert != "" {
				r.Answer = append(r.Answer, &dns.TLSA{Hdr: hdr, Usage: usg, Selector: sel, MatchingType: mat, Certificate: cert})
			}
		case dns.TypeURI:
			pri, wei, tar := db.Get.URI(q.Name)
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
	flag.String("dns.host", "127.0.0.1", "IP address to run the DNS server on")
	flag.Int("dns.port", 53, "Port for the DNS server to listen on")
	flag.String("dns.database", "./records.db", "Database file to use")
	flag.Bool("dns.disable-tcp", false, "Disable listening on TCP")
	flag.Bool("dns.disable-udp", false, "Disable listening on UDP")
	flag.String("http.host", "127.0.0.1", "IP address to run the API on")
	flag.Int("http.port", 8080, "Port for the API to listen on")
	flag.Bool("http.disabled", false, "Disable the API entirely")
	pflag.CommandLine.AddGoFlagSet(flag.CommandLine)
	pflag.Parse()
	if err := viper.BindPFlags(pflag.CommandLine); err != nil { log.Fatalf("Failed to setup command line arguments: %v", err) }

	// Set configuration defaults
	viper.SetDefault("dns.host", "127.0.0.1")
	viper.SetDefault("dns.port", 53)
	viper.SetDefault("dns.database", "./records.db")
	viper.SetDefault("dns.disable-tcp", false)
	viper.SetDefault("dns.disable-udp", false)

	viper.SetDefault("http.host", "127.0.0.1")
	viper.SetDefault("http.port", 8080)
	viper.SetDefault("http.disabled", false)

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
	database, err = bolt.Open(viper.GetString("dns.database"), 0666, nil)
	if err != nil {
		log.Fatalf("Failed to open database: %v", err)
	}
	defer func() { if err := database.Close(); err != nil { log.Fatalf("Failed to close database: %v", err) }}()

	// Setup database structure
	if err := db.Setup(database); err != nil {
		log.Fatalf("Failed setting up database structure: %v", err)
	}

	// Check config is valid
	if viper.GetBool("dns.disable-tcp") && viper.GetBool("dns.disable-udp") { log.Fatalf("Invalid configuration: tcp and/or udp must be enabled, got both as disabled") }

	// Handle TCP connections
	tcpErr := make(chan error)
	go func() {
		if viper.GetBool("dns.disable-tcp") { return }
		tcp := &dns.Server{Addr: viper.GetString("dns.host") + ":" + viper.GetString("dns.port"), Net: "tcp"}
		tcp.Handler = &handler{}

		if err := tcp.ListenAndServe(); err != nil { tcpErr <- err }
	}()

	// Handle UDP connections
	udpErr := make(chan error)
	go func() {
		if viper.GetBool("dns.disable-udp") { return }
		udp := &dns.Server{Addr: viper.GetString("dns.host") + ":" + viper.GetString("dns.port"), Net: "udp"}
		udp.Handler = &handler{}

		if err := udp.ListenAndServe(); err != nil { udpErr <- err }
	}()

	// Handle REST API
	httpErr := make(chan error)
	go func() {
		if viper.GetBool("http.disabled") { return }

		http.Handle("/records", handlers.LoggingHandler(os.Stdout, http.HandlerFunc(routes.AllRecordsHandler(database))))
		http.Handle("/records/", handlers.LoggingHandler(os.Stdout, http.HandlerFunc(routes.SingleRecordHandler("/records/", database))))
		if err := http.ListenAndServe(viper.GetString("http.host") + ":" + viper.GetString("http.port"), nil); err != nil { httpErr <- err }
	}()

	// Assemble log
	var protocols string
	if !viper.GetBool("dns.disable-udp") && !viper.GetBool("dns.disable-tcp") {
		protocols = "TCP and UDP"
	} else if viper.GetBool("dns.disable-udp") {
		protocols = "TCP"
	} else {
		protocols = "UDP"
	}
	log.Printf("DNS server listening on %s:%s with %s...", viper.GetString("dns.host"), viper.GetString("dns.port"), protocols)

	if !viper.GetBool("http.disabled") { log.Printf("HTTP server listening on %s:%s...", viper.GetString("http.host"), viper.GetString("http.port")) }

	// Watch for errors
	select {
	case err := <- tcpErr:
		log.Fatalf("DNS failed to listen on %s:%s with TCP: %v\n", viper.GetString("dns.host"), viper.GetString("dns.port"), err)
	case err := <- udpErr:
		log.Fatalf("DNS failed to listen on %s:%s with UDP: %v\n", viper.GetString("dns.host"), viper.GetString("dns.port"), err)
	case err := <- httpErr:
		log.Fatalf("API failed to listen on %s:%s: %v\n", viper.GetString("http.host"), viper.GetString("http.port"), err)
	}
}
