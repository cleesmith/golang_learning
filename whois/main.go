package main

import (
	"fmt"
	"net"
)

func whois(domainName, server string) string {
	conn, err := net.Dial("tcp", server+":43")

	if err != nil {
		fmt.Println("Error")
	}

	defer conn.Close()

	conn.Write([]byte(domainName + "\r\n"))

	buf := make([]byte, 1024)

	result := []byte{}

	for {
		numBytes, err := conn.Read(buf)
		sbuf := buf[0:numBytes]
		result = append(result, sbuf...)
		if err != nil {
			break
		}
	}

	return string(result)
}

func main() {
	// does NOT do IP addresses:
	// result := whois("82.165.177.154", "com.whois-servers.net")
	// but this does:
	// http://www.tcpiputils.com/browse/ip-address/82.165.177.154

	result := whois("socketloop.com", "com.whois-servers.net")
	fmt.Println(result)
}
