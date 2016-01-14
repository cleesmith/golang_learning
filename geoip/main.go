package main

import (
	"fmt"
	"log"
	"net"

	"github.com/oschwald/geoip2-golang"
)

func main() {
	db, err := geoip2.Open("GeoLite2-City.mmdb")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	// if you are using strings that may be invalid, check that ip is not nil
	ip := net.ParseIP("81.2.69.142") // GB London
	record, err := db.City(ip)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("record.City.Names=%#v\n", record.City.Names)
	fmt.Printf("Portuguese (BR) city name: %v\n", record.City.Names["pt-BR"])
	fmt.Printf("English subdivision name: %v\n", record.Subdivisions[0].Names["en"])
	fmt.Printf("Russian country name: %v\n", record.Country.Names["ru"])
	fmt.Printf("ISO country code: %v\n", record.Country.IsoCode)
	fmt.Printf("Time zone: %v\n", record.Location.TimeZone)
	fmt.Printf("Coordinates: %v, %v\n", record.Location.Latitude, record.Location.Longitude)
	// expected output:
	//   Portuguese (BR) city name: Londres
	//   English subdivision name: England
	//   Russian country name: Великобритания
	//   ISO country code: GB
	//   Time zone: Europe/London
	//   Coordinates: 51.5142, -0.0931

	fmt.Println("\ntest ipv6:")
	ip = net.ParseIP("::b110:c400") // BR Brazil
	// ip := net.ParseIP("::b0e4:6000") // IL Israel
	// my local ipv6 addresses, so not found:
	// ip := net.ParseIP("2001:0:5ef5:79fd:89e:3274:3390:b8dd")
	// ip := net.ParseIP("ff02::1")
	if ip.To4() == nil {
		// it's not IPv4, is it IPv6:
		if ip.To16() == nil {
			log.Fatal("invalid IP!")
		}
	}
	fmt.Printf("ip=%T=%#v\n", ip, ip)

	record, err = db.City(ip)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("record=%T=%#v\n", record, record)
	// fmt.Printf("Portuguese (BR) city name: %v\n", record.City.Names["pt-BR"])
	// fmt.Printf("English subdivision name: %v\n", record.Subdivisions[0].Names["en"])
	// fmt.Printf("Russian country name: %v\n", record.Country.Names["ru"])
	fmt.Printf("ISO country code: %v\n", record.Country.IsoCode)
	fmt.Printf("Time zone: %v\n", record.Location.TimeZone)
	fmt.Printf("Coordinates: %v, %v\n", record.Location.Latitude, record.Location.Longitude)
}
