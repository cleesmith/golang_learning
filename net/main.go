package main

import (
	"fmt"
	"net"
)

func main() {
	// ip4 := net.IP([]byte{1, 2, 3, 4})
	// ip := net.IP("ff02::1").String()
	// ip = net.IP("2001:0:9d38:90d7:1ce6:39af:3390:b8dd").String()

	var ip string = "::1"
	fmt.Println(isIP(ip))
	ip4, ip6, ips := isIP(ip)
	if ip4 {
		fmt.Printf("ip=%v is valid IPv4: '%v'\n", ip, ips)
	} else if ip6 {
		fmt.Printf("ip=%v is valid IPv6: '%v'\n", ip, ips)
	} else {
		fmt.Printf("ip=%v is NOT valid: '%v'\n", ip, ips)
	}
}

func isIP(s string) (ip4 bool, ip6 bool, ips string) {
	ip := net.ParseIP(s)
	if ip.To4() == nil {
		// it's not IPv4, is it IPv6:
		if ip.To16() == nil {
			return false, false, ""
		} else {
			return false, true, ip.String()
		}
	} else {
		return true, false, ip.String()
	}
}
