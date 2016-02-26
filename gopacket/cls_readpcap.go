package main

import (
	"log"

	"github.com/google/gopacket"
	"github.com/google/gopacket/layers"
	"github.com/google/gopacket/pcap"
)

func main() {
	for _, file := range []struct {
		filename       string
		num            int
		expectedLayers []gopacket.LayerType
	}{
		{"test_loopback.pcap",
			24,
			[]gopacket.LayerType{
				layers.LayerTypeLoopback,
				layers.LayerTypeIPv6,
				layers.LayerTypeTCP,
			},
		},
		{"test_ethernet.pcap",
			16,
			[]gopacket.LayerType{
				layers.LayerTypeEthernet,
				layers.LayerTypeIPv4,
				layers.LayerTypeTCP,
			},
		},
		{"test_dns.pcap",
			10,
			[]gopacket.LayerType{
				layers.LayerTypeEthernet,
				layers.LayerTypeIPv4,
				layers.LayerTypeUDP,
				layers.LayerTypeDNS,
			},
		},
	} {
		log.Printf("Processing file %s\n", file.filename)

		packets := []gopacket.Packet{}
		if handle, err := pcap.OpenOffline(file.filename); err != nil {
			log.Fatal(err)
		} else {
			packetSource := gopacket.NewPacketSource(handle, handle.LinkType())
			for packet := range packetSource.Packets() {
				packets = append(packets, packet)
			}
		}
		log.Printf("packets: got=%v want=%v\n", len(packets), file.num)
		if len(packets) != file.num {
			log.Fatal("Incorrect number of packets, want", file.num, "got", len(packets))
		}
		for i, p := range packets {
			// see issue: https://github.com/google/gopacket/issues/175
			// log.Printf(p.Dump())
			log.Printf("\n%v. packet:\n%#v\n", i+1, p)
			log.Printf("file.expectedLayers=%T=%#v\n", file.expectedLayers, file.expectedLayers)
			for _, layertype := range file.expectedLayers {
				if p.Layer(layertype) == nil {
					log.Fatal("Packet", i, "has no layer type\n%s", layertype, p.Dump())
				}
			}
		}
	}
}
