package main

import (
	"flag"
	"github.com/miekg/dns"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
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
			ip := getARecord(db, q.Name)
			if ip != nil {
				r.Answer = append(r.Answer, &dns.A{Hdr: hdr, A: ip})
			}
		case dns.TypeAAAA:
			ip := getAAAARecord(db, q.Name)
			if ip != nil {
				r.Answer = append(r.Answer, &dns.AAAA{Hdr: hdr, AAAA: ip})
			}
		case dns.TypeCNAME:
			target := getCNAMERecord(db, q.Name)
			if target != "" {
				r.Answer = append(r.Answer, &dns.CNAME{Hdr: hdr, Target: target})
			}
		case dns.TypeMX:
			host, priority := getMXRecord(db, q.Name)
			if host != "" {
				r.Answer = append(r.Answer, &dns.MX{Hdr: hdr, Preference: priority, Mx: host})
			}
		case dns.TypeLOC:
			vers, siz, hor, ver, lat, lon, alt := getLOCRecord(db, q.Name)
			if lat != 0 && lon != 0 {
				r.Answer = append(r.Answer, &dns.LOC{Hdr: hdr, Version: vers, Size: siz, HorizPre: hor, VertPre: ver, Latitude: lat, Longitude: lon, Altitude: alt})
			}
		case dns.TypeSRV:
			priority, weight, port, target := getSRVRecord(db, q.Name)
			if target != "" {
				r.Answer = append(r.Answer, &dns.SRV{Hdr: hdr, Priority: priority, Weight: weight, Port: port, Target: target})
			}
		case dns.TypeSPF:
			txt := getSPFRecord(db, q.Name)
			if len(txt) != 0 {
				r.Answer = append(r.Answer, &dns.SPF{Hdr: hdr, Txt: txt})
			}
		case dns.TypeTXT:
			content := getTXTRecord(db, q.Name)
			if len(content) != 0 {
				r.Answer = append(r.Answer, &dns.TXT{Hdr: hdr, Txt: content})
			}
		case dns.TypeNS:
			nameserver := getNSRecord(db, q.Name)
			if nameserver != "" {
				r.Answer = append(r.Answer, &dns.NS{Hdr: hdr, Ns: nameserver})
			}
		case dns.TypeCAA:
			dflag, tag, content := getCAARecord(db, q.Name)
			if tag != "" && content != "" {
				r.Answer = append(r.Answer, &dns.CAA{Hdr: hdr, Flag: dflag, Tag: tag, Value: content})
			}
		case dns.TypePTR:
			ptr := getPTRRecord(db, q.Name)
			if ptr != "" {
				r.Answer = append(r.Answer, &dns.PTR{Hdr: hdr, Ptr: ptr})
			}
		case dns.TypeCERT:
			tpe, tag, algo, cert := getCERTRecord(db, q.Name)
			if cert != "" {
				r.Answer = append(r.Answer, &dns.CERT{Hdr: hdr, Type: tpe, KeyTag: tag, Algorithm: algo, Certificate: cert})
			}
		case dns.TypeDNSKEY:
			flags, proto, algo, pub := getDNSKEYRecord(db, q.Name)
			if pub != "" {
				r.Answer = append(r.Answer, &dns.DNSKEY{Hdr: hdr, Flags: flags, Protocol: proto, Algorithm: algo, PublicKey: pub})
			}
		case dns.TypeDS:
			ktag, algo, dtype, digest := getDSRecord(db, q.Name)
			if digest != "" {
				r.Answer = append(r.Answer, &dns.DS{Hdr: hdr, KeyTag: ktag, Algorithm: algo, DigestType: dtype, Digest: digest})
			}
		case dns.TypeNAPTR:
			ord, pref, dflag, serv, reg, rep := getNAPTRRecord(db, q.Name)
			if serv != "" && reg != "" && rep != "" {
				r.Answer = append(r.Answer, &dns.NAPTR{Hdr: hdr, Order: ord, Preference: pref, Flags: dflag, Service: serv, Regexp: reg, Replacement: rep})
			}
		case dns.TypeSMIMEA:
			usage, sel, match, cert := getSMIMEARecord(db, q.Name)
			if cert != "" {
				r.Answer = append(r.Answer, &dns.SMIMEA{Hdr: hdr, Usage: usage, Selector: sel, MatchingType: match, Certificate: cert})
			}
		case dns.TypeSSHFP:
			algo, tpe, fingerprint := getSSHFPRecord(db, q.Name)
			if fingerprint != "" {
				r.Answer = append(r.Answer, &dns.SSHFP{Hdr: hdr, Algorithm: algo, Type: tpe, FingerPrint: fingerprint})
			}
		case dns.TypeTLSA:
			usg, sel, mat, cert := getTLSARecord(db, q.Name)
			if cert != "" {
				r.Answer = append(r.Answer, &dns.TLSA{Hdr: hdr, Usage: usg, Selector: sel, MatchingType: mat, Certificate: cert})
			}
		case dns.TypeURI:
			pri, wei, tar := getURIRecord(db, q.Name)
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

	logResponse(w, r, start)
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
	db, err = bolt.Open(viper.GetString("database"), 0666, nil)
	if err != nil {
		log.Fatalf("Failed to open database: %v", err)
	}
	defer func() { if err := db.Close(); err != nil { log.Fatalf("Failed to close database: %v", err) }}()

	// Setup database structure
	if err := setupDB(db); err != nil {
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

	var protocols string
	if !viper.GetBool("no-udp") && !viper.GetBool("no-tcp") {
		protocols = "TCP and UDP"
	} else if viper.GetBool("no-udp") {
		protocols = "TCP"
	} else {
		protocols = "UDP"
	}
	log.Printf("Listening on %s:%s with %s...", viper.GetString("host"), viper.GetString("port"), protocols)

	// Watch for errors
	select {
	case err := <- tcpErr:
		log.Fatalf("Failed to listen on 127.0.0.1:1052 with TCP: %v\n", err)
	case err := <- udpErr:
		log.Fatalf("Failed to listen on 127.0.0.1:1053 with UDP: %v\n", err)
	}
}
