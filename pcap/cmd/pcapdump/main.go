// Copyright 2015 David du Colombier. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"encoding/binary"
	"flag"
	"fmt"
	"io"
	"log"
	"net"
	"os"

	// "io/ioutil"

	"github.com/0intro/pcap"
)

type Ethernet struct {
	DstAddr []byte
	SrcAddr []byte
	Type    uint16
}

var verbose = flag.Bool("v", false, "verbose")

func usage() {
	fmt.Fprintf(os.Stderr, "usage: pcapdump [ -v ] file.pcap\n")
	os.Exit(1)
}

func main() {
	flag.Parse()
	args := flag.Args()
	if flag.NArg() != 1 {
		usage()
	}

	f, err := os.Open(args[0])
	if err != nil {
		log.Fatal(err)
	}
	defer func() {
		if err := f.Close(); err != nil {
			log.Fatal(err)
		}
	}()

	var totRecs int

	pr, err := pcap.NewReader(f)
	if err != nil {
		log.Fatal(err)
	}

	if *verbose {
		fmt.Println("Header")
		fmt.Printf("Magic 0x%.8x\n", pr.Header.Magic)
		fmt.Println("VersionMajor", pr.Header.VersionMajor)
		fmt.Println("VersionMinor", pr.Header.VersionMinor)
		fmt.Println("ThisZone", pr.Header.ThisZone)
		fmt.Println("SigFigs", pr.Header.SigFigs)
		fmt.Println("SnapLen", pr.Header.SnapLen)
		fmt.Println("LinkType", pr.Header.LinkType)
	}

	for {

		record, err := pr.Next()
		if err != nil {
			if err == io.EOF {
				break
			}
			log.Fatal(err)
		}
		totRecs++

		if *verbose {
			// https://wiki.wireshark.org/Development/LibpcapFileFormat
			fmt.Printf("PacketRecordHeader(%T):\tTsSecond=%v\tTsMicrosecond=%v\tCapLen=%v\tLen=%v\n", record, record.TsSec, record.TsUsec, record.CapLen, record.Len)
		}

		var ip []byte
		ip = make([]byte, 4)
		if err = binary.Read(pr, binary.BigEndian, &ip); err != nil {
			log.Fatalln(err)
		}
		// fmt.Printf("ip=%T=%#v=%#v\n", ip, ip, string(ip[:]))
		nip := net.IP(ip).String()
		fmt.Printf("nip=%T=%#v\n", nip, nip)
		ip4, ip6, ips := isIP(nip)
		fmt.Printf("ip4=%v ip6=%v ips=%v\n", ip4, ip6, ips)

		// eth := &Ethernet{}
		// if err = binary.Read(pr, binary.BigEndian, eth); err != nil {
		// 	log.Fatalln(err)
		// }
		// _, _, ips := isIP(net.IP(eth.SrcAddr).String())
		// fmt.Printf("%d.%.6d %d %v -> %x %x\n", record.TsSec, record.TsUsec, record.Len, ips, eth.DstAddr, eth.Type)

		// buf, err := ioutil.ReadAll(pr)
		// if err != nil {
		// 	log.Fatal(err)
		// }
		// if *verbose {
		// 	fmt.Println()
		// 	fmt.Println("Payload")
		// 	fmt.Println(buf)
		// }

	} // for
	fmt.Printf("Records=%v\n", totRecs)

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
