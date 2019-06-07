package db

import "net"

type Record interface {
	Name() string
}

// Parts of an A record
type A struct {
	Address net.IP `json:"address"`
}
func (a A) Name() string { return "A" }

// Parts of an AAAA record
type AAAA struct {
	Address net.IP `json:"address"`
}
func (a AAAA) Name() string { return "AAAA" }

// Parts of a CNAME record
type CNAME struct {
	Target string `json:"target"`
}
func (c CNAME) Name() string { return "CNAME" }

// Parts of a MX record
type MX struct {
	Host     string `json:"host"`
	Priority uint16 `json:"priority"`
}
func (m MX) Name() string { return "MX" }

// Parts of a LOC record
type LOC struct {
	Version             uint8  `json:"version"`
	Size                uint8  `json:"size"`
	HorizontalPrecision uint8  `json:"horizontal-precision"`
	VerticalPrecision   uint8  `json:"vertical-precision"`
	Latitude            uint32 `json:"latitude"`
	Longitude           uint32 `json:"longitude"`
	Altitude            uint32 `json:"altitude"`
}
func (l LOC) Name() string { return "LOC" }

// Parts of a SRV record
type SRV struct {
	Priority uint16 `json:"priority"`
	Weight   uint16 `json:"weight"`
	Port     uint16 `json:"port"`
	Target   string `json:"target"`
}
func (s SRV) Name() string { return "SRV" }

// Parts of a SPF record
type SPF struct {
	Text []string `json:"text"`
}
func (s SPF) Name() string { return "SPF" }

// Parts of a TXT record
type TXT struct {
	Text []string `json:"text"`
}
func (t TXT) Name() string { return "TXT" }

// Parts of a NS record
type NS struct {
	Nameserver string `json:"nameserver"`
}
func (n NS) Name() string { return "NS" }

// Parts of a CAA record
type CAA struct {
	Flag    uint8  `json:"flag"`
	Tag     string `json:"tag"`
	Content string `json:"content"`
}
func (c CAA) Name() string { return "CAA" }

// Parts of a PTR record
type PTR struct {
	Domain string `json:"domain"`
}
func (p PTR) Name() string { return "PTR" }

// Parts of a CERT record
type CERT struct {
	Type        uint16 `json:"type"`
	KeyTag      uint16 `json:"key-tag"`
	Algorithm   uint8  `json:"algorithm"`
	Certificate string `json:"certificate"`
}
func (c CERT) Name() string { return "CERT" }

// Parts of a DNSKEY record
type DNSKEY struct {
	Flags     uint16 `json:"flags"`
	Protocol  uint8  `json:"protocol"`
	Algorithm uint8  `json:"algorithm"`
	PublicKey string `json:"public-key"`
}
func (d DNSKEY) Name() string { return "DNSKEY" }

// Parts of a DS record
type DS struct {
	KeyTag     uint16 `json:"key-tag"`
	Algorithm  uint8  `json:"algorithm"`
	DigestType uint8  `json:"digest-type"`
	Digest     string `json:"digest"`
}
func (d DS) Name() string { return "DS" }

// Parts of a NAPTR record
type NAPTR struct {
	Order       uint16 `json:"order"`
	Preference  uint16 `json:"preference"`
	Flags       string `json:"flags"`
	Service     string `json:"service"`
	Regexp      string `json:"regexp"`
	Replacement string `json:"replacement"`
}
func (n NAPTR) Name() string { return "NAPTR" }

// Parts of a SMIMEA record
type SMIMEA struct {
	Usage        uint8  `json:"usage"`
	Selector     uint8  `json:"selector"`
	MatchingType uint8  `json:"matching-type"`
	Certificate  string `json:"certificate"`
}
func (s SMIMEA) Name() string { return "SMIMEA" }

// Parts of a SSHFP record
type SSHFP struct {
	Algorithm   uint8  `json:"algorithm"`
	Type        uint8  `json:"type"`
	Fingerprint string `json:"fingerprint"`
}
func (s SSHFP) Name() string { return "SSHFP" }

// Parts of a TLSA record
type TLSA struct {
	Usage        uint8  `json:"usage"`
	Selector     uint8  `json:"selector"`
	MatchingType uint8  `json:"matching-type"`
	Certificate  string `json:"certificate"`
}
func (t TLSA) Name() string { return "TLSA" }

// Parts of a URI record
type URI struct {
	Priority uint16 `json:"priority"`
	Weight   uint16 `json:"weight"`
	Target   string `json:"target"`
}
func (u URI) Name() string { return "URI" }
