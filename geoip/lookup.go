package main

import (
	"fmt"
	"log"
	"net"

	"github.com/oschwald/geoip2-golang"
)

var GeoIp2Lookup *geoip2.Reader

func main() {
	GeoIp2Lookup, err := geoip2.Open("GeoLite2-City.mmdb") // ~40MB
	if err != nil {
		fmt.Printf("\nGeoIp2Lookup: type=%T \t value=%#v\n", GeoIp2Lookup, GeoIp2Lookup)
		GeoIp2Lookup = nil
		if GeoIp2Lookup == nil {
			fmt.Printf("GeoIp2Lookup is nil!\n")
			fmt.Printf("GeoIp2Lookup: type=%T \t value=%#v\n", GeoIp2Lookup, GeoIp2Lookup)
		}
		log.Fatal(err)
	}
	defer GeoIp2Lookup.Close()
	fmt.Printf("\nGeoIp2Lookup=%T=%#v\n", GeoIp2Lookup, GeoIp2Lookup)

	// if you are using strings that may be invalid, check that ip is not nil
	// ip := net.ParseIP("8?1.2.69.142") // GB London
	// ip := net.ParseIP("81.2.69.142") // GB Arnold (near London)
	ip := net.ParseIP("0.0.0.0")
	if ip == nil {
		log.Fatal("net.ParseIP: ip is nil!")
	}
	// ip = nil
	// *************************************************************************
	// note: GeoIp2Lookup.City() always returns a *geoip2.City struct, whereas
	//       GeoLite.GetLocationByIP() may return nil when nothing is found
	// *************************************************************************
	loc, err := GeoIp2Lookup.City(ip)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("\nIP: %v\n", ip)
	fmt.Printf("Lat/Lng coordinates: %v, %v\n", loc.Location.Latitude, loc.Location.Longitude)
	fmt.Printf("ISO country code: %v\n", loc.Country.IsoCode)
	fmt.Printf("City: %#v\n", loc.City.Names["en"])
	fmt.Printf("loc: type=%T \t value=%#v\n", loc, loc)
}
