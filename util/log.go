package util

import (
	"fmt"
	"github.com/coredns/coredns/request"
	"github.com/miekg/dns"
	"strconv"
	"strings"
	"time"
)

// From CoreDNS logging plugin
func LogResponse(w dns.ResponseWriter, r *dns.Msg, start time.Time) {
	state := request.Request{W: w, Req: r}

	fmt.Printf("%s %s:%s - %s \"%s %s %s %s %s %s %s\" %s %s %s %s\n",
		convertRemote(state.IP()), state.Port(), strconv.Itoa(int(state.Req.Id)),
		state.Type(), state.Class(), state.Name(), state.Proto(),
		strconv.Itoa(state.Req.Len()), boolToString(state.Do()), strconv.Itoa(state.Size()),
		getRCode(r), getRFlags(r), getRSize(r), strconv.FormatFloat(time.Since(start).Seconds(), 'f', -1, 64) + "s",
		time.Now().Format("2006-01-02 15:04:05"))
}

func convertRemote(addr string) string {
	if strings.Contains(addr, ":") {
		return "[" + addr + "]"
	}
	return addr
}

func boolToString(b bool) string {
	if b { return "true" }
	return "false"
}

func getRCode(r *dns.Msg) string {
	if r == nil {
		return "-"
	}
	rcode := dns.RcodeToString[r.Rcode]
	if rcode == "" {
		rcode = strconv.Itoa(r.Rcode)
	}
	return rcode
}

func getRFlags(r *dns.Msg) string {
	if r == nil {
		return "-"
	}
	h := r.MsgHdr

	flags := make([]string, 7)
	i := 0

	if h.Response {
		flags[i] = "qr"
		i++
	}
	if h.Authoritative {
		flags[i] = "aa"
		i++
	}
	if h.Truncated {
		flags[i] = "tc"
		i++
	}
	if h.RecursionDesired {
		flags[i] = "rd"
		i++
	}
	if h.RecursionAvailable {
		flags[i] = "ra"
		i++
	}
	if h.Zero {
		flags[i] = "z"
		i++
	}
	if h.AuthenticatedData {
		flags[i] = "ad"
		i++
	}
	if h.CheckingDisabled {
		flags[i] = "cd"
		i++
	}
	return strings.Join(flags[:i], ",")
}

func getRSize(r *dns.Msg) string {
	if r == nil {
		return "-"
	}
	// Not sure how to get this, so return 0
	return "0"
}
