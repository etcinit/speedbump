package ginbump

import (
	"net"
	"net/http"
	"strings"
)

// Originally from: https://github.com/sebest/xff/blob/master/xff.go

var privateMasks = func() []net.IPNet {
	masks := []net.IPNet{}
	for _, cidr := range []string{"10.0.0.0/8", "172.16.0.0/12", "192.168.0.0/16", "fc00::/7"} {
		if _, network, err := net.ParseCIDR(cidr); err != nil {
			panic(err)
		} else {
			masks = append(masks, *network)
		}
	}
	return masks
}()

// IsPublicIP returns true if the given IP can be routed on the Internet
func IsPublicIP(ip net.IP) bool {
	if !ip.IsGlobalUnicast() {
		return false
	}

	for _, mask := range privateMasks {
		if mask.Contains(ip) {
			return false
		}
	}

	return true
}

// ParseForwarded parses the value of the X-Forwarded-For Header and returns the
// IP address.
func ParseForwarded(ipList string) string {
	for _, ip := range strings.Split(ipList, ",") {
		ip = strings.TrimSpace(ip)

		if parsed := net.ParseIP(ip); parsed != nil && IsPublicIP(parsed) {
			return ip
		}
	}

	return ""
}

// GetRequesterAddress does a best effort lookup for the real IP address of the
// requester. Many load balancers (such as AWS's ELB) set a X-Forwarded-For
// header which can be used to determine the IP address of the client when the
// server is behind a load balancer.
//
// It is possible however for the client to spoof this header if the load
// balancer is not configured to remove it from the request or if the server is
// accessed directly.
//
// For uses such as rate limitting, only use this function if you can trust that
// the load balancer will strip the header from the client and that the server
// will not be directly accessible by the public (only though the load
// balancer).
func GetRequesterAddress(r *http.Request) string {
	if xff := r.Header.Get("X-Forwarded-For"); xff != "" {
		return ParseForwarded(xff)
	}

	if ip, _, err := net.SplitHostPort(r.RemoteAddr); err == nil {
		return ip
	}

	return ""
}
