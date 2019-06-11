package main

import (
	"flag"
	"github.com/akrantz01/krantz.dev/dns/db"
	"github.com/akrantz01/krantz.dev/dns/records"
	"github.com/akrantz01/krantz.dev/dns/roles"
	"github.com/akrantz01/krantz.dev/dns/users"
	"github.com/akrantz01/krantz.dev/dns/util"
	"github.com/gorilla/handlers"
	"github.com/miekg/dns"
	"github.com/rs/cors"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
	bolt "go.etcd.io/bbolt"
	"gopkg.in/hlandau/passlib.v1"
	"log"
	"net/http"
	"os"
	"time"
)

var database *bolt.DB

type handler struct {}
func (h *handler) ServeDNS(w dns.ResponseWriter, m *dns.Msg) {
	// Set database into getter and setter
	db.Get.Db = database
	db.Set.Db = database

	// Time request for logging
	start := time.Now()

	// Assemble response
	r := new(dns.Msg)
	r.SetReply(m)
	r.Authoritative = true

	// Iterate over all questions
	for _, q := range r.Question {
		hdr := dns.RR_Header{Name: q.Name, Rrtype: q.Qtype, Class: q.Qclass}

		// Do different things based on record type
		switch q.Qtype {
		case dns.TypeA:
			record := db.Get.A(q.Name)
			if record != nil {
				r.Answer = append(r.Answer, &dns.A{Hdr: hdr, A: record.Address})
			}
		case dns.TypeAAAA:
			record :=  db.Get.AAAA(q.Name)
			if record != nil {
				r.Answer = append(r.Answer, &dns.AAAA{Hdr: hdr, AAAA: record.Address})
			}
		case dns.TypeCNAME:
			record :=  db.Get.CNAME(q.Name)
			if record != nil {
				r.Answer = append(r.Answer, &dns.CNAME{Hdr: hdr, Target: record.Target})
			}
		case dns.TypeMX:
			record :=  db.Get.MX(q.Name)
			if record != nil {
				r.Answer = append(r.Answer, &dns.MX{Hdr: hdr, Preference: record.Priority, Mx: record.Host})
			}
		case dns.TypeLOC:
			record :=  db.Get.LOC(q.Name)
			if record != nil {
				locString, vers := record.ToParsable()
				r.Answer = append(r.Answer, util.ParseLOCString(locString, vers, hdr))
			}
		case dns.TypeSRV:
			record :=  db.Get.SRV(q.Name)
			if record != nil {
				r.Answer = append(r.Answer, &dns.SRV{Hdr: hdr, Priority: record.Priority, Weight: record.Weight, Port: record.Port, Target: record.Target})
			}
		case dns.TypeSPF:
			record :=  db.Get.SPF(q.Name)
			if record != nil {
				r.Answer = append(r.Answer, &dns.SPF{Hdr: hdr, Txt: record.Text})
			}
		case dns.TypeTXT:
			record :=  db.Get.TXT(q.Name)
			if record != nil {
				r.Answer = append(r.Answer, &dns.TXT{Hdr: hdr, Txt: record.Text})
			}
		case dns.TypeNS:
			record :=  db.Get.NS(q.Name)
			if record != nil {
				r.Answer = append(r.Answer, &dns.NS{Hdr: hdr, Ns: record.Nameserver})
			}
		case dns.TypeCAA:
			record :=  db.Get.CAA(q.Name)
			if record != nil {
				r.Answer = append(r.Answer, &dns.CAA{Hdr: hdr, Flag: record.Flag, Tag: record.Tag, Value: record.Content})
			}
		case dns.TypePTR:
			record := db.Get.PTR(q.Name)
			if record != nil {
				r.Answer = append(r.Answer, &dns.PTR{Hdr: hdr, Ptr: record.Domain})
			}
		case dns.TypeCERT:
			record :=  db.Get.CERT(q.Name)
			if record != nil {
				r.Answer = append(r.Answer, &dns.CERT{Hdr: hdr, Type: record.Type, KeyTag: record.KeyTag, Algorithm: record.Algorithm, Certificate: record.Certificate})
			}
		case dns.TypeDNSKEY:
			record :=  db.Get.DNSKEY(q.Name)
			if record != nil {
				r.Answer = append(r.Answer, &dns.DNSKEY{Hdr: hdr, Flags: record.Flags, Protocol: record.Protocol, Algorithm: record.Algorithm, PublicKey: record.PublicKey})
			}
		case dns.TypeDS:
			record :=  db.Get.DS(q.Name)
			if record != nil {
				r.Answer = append(r.Answer, &dns.DS{Hdr: hdr, KeyTag: record.KeyTag, Algorithm: record.Algorithm, DigestType: record.DigestType, Digest: record.Digest})
			}
		case dns.TypeNAPTR:
			record :=  db.Get.NAPTR(q.Name)
			if record != nil {
				r.Answer = append(r.Answer, &dns.NAPTR{Hdr: hdr, Order: record.Order, Preference: record.Preference, Flags: record.Flags, Service: record.Service, Regexp: record.Regexp, Replacement: record.Replacement})
			}
		case dns.TypeSMIMEA:
			record :=  db.Get.SMIMEA(q.Name)
			if record != nil {
				r.Answer = append(r.Answer, &dns.SMIMEA{Hdr: hdr, Usage: record.Usage, Selector: record.Selector, MatchingType: record.MatchingType, Certificate: record.Certificate})
			}
		case dns.TypeSSHFP:
			record :=  db.Get.SSHFP(q.Name)
			if record != nil {
				r.Answer = append(r.Answer, &dns.SSHFP{Hdr: hdr, Algorithm: record.Algorithm, Type: record.Type, FingerPrint: record.Fingerprint})
			}
		case dns.TypeTLSA:
			record :=  db.Get.TLSA(q.Name)
			if record != nil {
				r.Answer = append(r.Answer, &dns.TLSA{Hdr: hdr, Usage: record.Usage, Selector: record.Selector, MatchingType: record.MatchingType, Certificate: record.Certificate})
			}
		case dns.TypeURI:
			record :=  db.Get.URI(q.Name)
			if record != nil {
				r.Answer = append(r.Answer, &dns.URI{Hdr: hdr, Priority: record.Priority, Weight: record.Weight, Target: record.Target})
			}
		default:
			r.Rcode = dns.RcodeNameError
		}
	}

	// Throw error if no answers
	if len(r.Answer) == 0 {
		r.Rcode = dns.RcodeNameError
	}

	// Write response
	if err := w.WriteMsg(r); err != nil {
		log.Printf("Unable to send response: %v", err)
	}

	// Log to console
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
	flag.String("http.admin.name", "DNS Admin", "Name of the admin user")
	flag.String("http.admin.username", "admin", "Username of the admin user")
	flag.String("http.admin.password", "admin", "Password of the admin user")
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
	viper.SetDefault("http.admin.name", "DNS Admin")
	viper.SetDefault("http.admin.username", "admin")
	viper.SetDefault("http.admin.password", "admin")
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

	// Setup hashing
	if err := passlib.UseDefaults(passlib.DefaultsLatest); err != nil {
		log.Fatal("invalid hash configuration")
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

		// Allow CORS
		c := cors.AllowAll()

		// Setup API routes
		http.Handle("/api/records", c.Handler(handlers.LoggingHandler(os.Stdout, http.HandlerFunc(records.AllRecordsHandler(database)))))
		http.Handle("/api/records/", c.Handler(handlers.LoggingHandler(os.Stdout, http.HandlerFunc(records.SingleRecordHandler("/api/records/", database)))))
		http.Handle("/api/users", c.Handler(handlers.LoggingHandler(os.Stdout, http.HandlerFunc(users.AllUsersHandler(database)))))
		http.Handle("/api/users/login", c.Handler(handlers.LoggingHandler(os.Stdout, http.HandlerFunc(users.Login(database)))))
		http.Handle("/api/users/logout", c.Handler(handlers.LoggingHandler(os.Stdout, http.HandlerFunc(users.Logout(database)))))
		http.Handle("/api/roles", c.Handler(handlers.LoggingHandler(os.Stdout, http.HandlerFunc(roles.AllRolesHandler(database)))))
		http.Handle("/api/roles/", c.Handler(handlers.LoggingHandler(os.Stdout, http.HandlerFunc(roles.SingleRoleHandler("/api/roles/", database)))))
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
