package util

// from https://gist.github.com/jgrahamc/9807839

import (
	"github.com/miekg/dns"
	"regexp"
	"strconv"
)

// locReD is the regexp to capture a value in a latitude or longitude.
//
// locReM is for the other values (they can be negative) and can have
// an optional 'm' after. locReOM is an optional version of
//
// locReM. Note that the m character after a number has no meaning at
// all, the values are always in metres.
var locReD = "(\\d+)(?: (\\d+))?(?: (\\d+(?:\\.\\d+)?))?"
var locReM = "(?: (-?\\d+(?:\\.\\d+)?)m?)"
var locReOM = locReM + "?"
var locRe, _ = regexp.Compile(locReD + " (N|S) " + locReD + " (E|W)" + locReM + locReOM + locReOM + locReOM)

// parseSizePrecision parses the siz, hp and vp parts of a LOC string
// and returns them in the weird 8 bit format required. See RFC 1876
// for specification and justification. The p string contains the LOC
// value to parse. It may be empty in which case the default value d
// is returned.  The boolean return is false if the parsing fails.
func parseSizePrecision(p string, d uint8) (uint8, bool) {
	if p == "" {
		return d, true
	}

	f, err := strconv.ParseFloat(p, 64)
	if err != nil || f < 0 || f > 90000000 {
		return 0, false
	}

	// Conversion from m to cm
	f *= 100

	var exponent uint8 = 0
	for f >= 10 {
		exponent += 1
		f /= 10
	}

	// Here both f and exponent will be in the range 0 to 9 and these
	// get packed into a byte in the following manner. The result?
	// Look at the value in hex and you can read it. e.g. 6e3 (i.e. 6000) is 0x63
	return uint8(f) << 4 + exponent, true
}

// parseLatLong parses a latitude/longitude string (see ParseString
// below for format) and returns the value as a single uint32. If the
// bool value is false there was a problem with the format. The limit
// parameter specifies the limit for the number of degrees.
func parseLatLong(d, m, s string, limit uint64) (uint32, bool) {
	n, err := strconv.ParseUint(d, 10, 8)
	if err != nil || n > limit {
		return 0, false
	}
	pos := float64(n) * 60

	if m != "" {
		n, err := strconv.ParseUint(m, 10, 8)
		if err != nil || n > 59 {
			return 0, false
		}
		pos += float64(n)
	}

	pos *= 60

	if s != "" {
		f, err := strconv.ParseFloat(s, 64)
		if err != nil || f > 59.999 {
			return 0, false
		}
		pos += f
	}

	pos *= 1000

	return uint32(pos), pos <= float64(limit * dns.LOC_DEGREES)
}

// parseLOCString parses the string representation of a LOC record and
// fills in the fields in a newly created LOC appropriately. If the
// function returns nil then there was a parsing error, otherwise
// returns a pointer to a new LOC.
//
// Worth reading RFC 1876, Appendix A to understand.
func ParseLOCString(l string, version uint8, header dns.RR_Header) *dns.LOC {
	loc := new(dns.LOC)

	// The string l will be in the following format:
	//
	// d1 [m1 [s1]] {"N"|"S"} d2 [m2 [s2]] {"E"|"W"}
	// alt["m"] [siz["m"] [hp["m"] [vp["m"]]]]
	//
	// d1 is the latitude, d2 is the longitude, alt is the altitude,
	// siz is the size of the planet, hp and vp are the horiz and vert
	// precisions. See RFC 1876 for full detail.
	//
	// Examples:
	// 42 21 54 N 71 06 18 W -24m 30m
	// 42 21 43.952 N 71 5 6.344 W -24m 1m 200m
	// 52 14 05 N 00 08 50 E 10m
	// 2 7 19 S 116 2 25 E 10m
	// 42 21 28.764 N 71 00 51.617 W -44m 2000m
	// 59 N 10 E 15.0 30.0 2000.0 5.0

	parts := locRe.FindStringSubmatch(l)
	if parts == nil {
		return nil
	}

	// Quick reference to the matches
	//
	// parts[1] == latitude degrees
	// parts[2] == latitude minutes (optional)
	// parts[3] == latitude seconds (optional)
	// parts[4] == N or S
	//
	// parts[5] == longitude degrees
	// parts[6] == longitude minutes (optional)
	// parts[7] == longitude seconds (optional)
	// parts[8] == E or W
	//
	// parts[9] == altitude
	//
	// These are completely optional:
	//
	// parts[10] == size
	// parts[11] == horizontal precision
	// parts[12] == vertical precision

	// Convert latitude and longitude to a 32-bit unsigned integer
	latitude, ok := parseLatLong(parts[1], parts[2], parts[3], 90)
	if !ok {
		return nil
	}
	loc.Latitude = dns.LOC_EQUATOR
	if parts[4] == "N" {
		loc.Latitude += latitude
	} else {
		loc.Latitude -= latitude
	}

	longitude, ok := parseLatLong(parts[5], parts[6], parts[7], 180)
	if !ok {
		return nil
	}
	loc.Longitude = dns.LOC_PRIMEMERIDIAN
	if parts[8] == "E" {
		loc.Longitude += longitude
	} else {
		loc.Longitude -= longitude
	}

	// Now parse the altitude. Seriously, read RFC 1876 if you want to
	// understand all the values and conversions here. But altitudes
	// are unsigned 32-bit numbers that start 100,000m below 'sea
	// level' and are expressed in cm.
	//
	// == (2^32-1)/100
	//  - 100,000
	// == 42949672.95
	//  -   100000
	// == 42849672.95
	f, err := strconv.ParseFloat(parts[9], 64)
	if err != nil || f < -dns.LOC_ALTITUDEBASE || f > 42849672.95 {
		return nil
	}
	loc.Altitude = (uint32)((f + dns.LOC_ALTITUDEBASE) * 100)

	// Default values for the optional components, see RFC 1876 for
	// this weird encoding. But top nibble is mantissa, bottom nibble
	// is exponent. Values are in cm. So, for example, 0x12 means 1 *
	// 10^2 or 100cm.
	//
	// 0x12 == 1e2cm == 1m
	if loc.Size, ok = parseSizePrecision(parts[10], 0x12); !ok {
		return nil
	}
	// 0x16 == 1e6cm == 10,000m == 10km
	if loc.HorizPre, ok = parseSizePrecision(parts[11], 0x16); !ok {
		return nil
	}
	// 0x13 == 1e3cm == 10m
	if loc.VertPre, ok = parseSizePrecision(parts[12], 0x13); !ok {
		return nil
	}

	loc.Version = version
	loc.Hdr = header

	return loc
}
