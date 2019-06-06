package db

import "net"

// Parts of an A record
type A struct {
	Address net.IP `json:"address"`
}

// Parts of an AAAA record
type AAAA struct {
	Address net.IP `json:"address"`
}

// Parts of a CNAME record
type CNAME struct {
	Target string `json:"target"`
}

// Parts of a MX record
type MX struct {
	Host     string `json:"host"`
	Priority uint16 `json:"priority"`
}

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

// Parts of a SRV record
type SRV struct {
	Priority uint16 `json:"priority"`
	Weight   uint16 `json:"weight"`
	Port     uint16 `json:"port"`
	Target   string `json:"target"`
}

// Parts of a SPF record
type SPF struct {
	Text []string `json:"text"`
}

// Parts of a TXT record
type TXT struct {
	Text []string `json:"text"`
}

// Parts of a NS record
type NS struct {
	Nameserver string `json:"nameserver"`
}

// Parts of a CAA record
type CAA struct {
	Flag    uint8  `json:"flag"`
	Tag     string `json:"tag"`
	Content string `json:"content"`
}

// Parts of a PTR record
type PTR struct {
	Domain string `json:"domain"`
}

// Parts of a CERT record
type CERT struct {
	Type        uint16 `json:"type"`
	KeyTag      uint16 `json:"key-tag"`
	Algorithm   uint8  `json:"algorithm"`
	Certificate string `json:"certificate"`
}

// Parts of a DNSKEY record
type DNSKEY struct {
	Flags     uint16 `json:"flags"`
	Protocol  uint8  `json:"protocol"`
	Algorithm uint8  `json:"algorithm"`
	PublicKey string `json:"public-key"`
}

// Parts of a DS record
type DS struct {
	KeyTag     uint16 `json:"key-tag"`
	Algorithm  uint8  `json:"algorithm"`
	DigestType uint8  `json:"digest-type"`
	Digest     string `json:"digest"`
}

// Parts of a NAPTR record
type NAPTR struct {
	Order       uint16 `json:"order"`
	Preference  uint16 `json:"preference"`
	Flags       string `json:"flags"`
	Service     string `json:"service"`
	Regexp      string `json:"regexp"`
	Replacement string `json:"replacement"`
}

// Parts of a SMIMEA record
type SMIMEA struct {
	Usage        uint8  `json:"usage"`
	Selector     uint8  `json:"selector"`
	MatchingType uint8  `json:"matching-type"`
	Certificate  string `json:"certificate"`
}

// Parts of a SSHFP record
type SSHFP struct {
	Algorithm   uint8  `json:"algorithm"`
	Type        uint8  `json:"type"`
	Fingerprint string `json:"fingerprint"`
}

// Parts of a TLSA record
type TLSA struct {
	Usage        uint8  `json:"usage"`
	Selector     uint8  `json:"selector"`
	MatchingType uint8  `json:"matching-type"`
	Certificate  string `json:"certificate"`
}

// Parts of a URI record
type URI struct {
	Priority uint16 `json:"priority"`
	Weight   uint16 `json:"weight"`
	Target   string `json:"target"`
}
