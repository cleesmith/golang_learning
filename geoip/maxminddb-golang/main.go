package main

import (
	"fmt"
	"log"
	_ "net"

	"github.com/oschwald/maxminddb-golang"
)

func main() {
	// db, err := maxminddb.Open("GeoIP2-Connection-Type-Test.mmdb")
	db, err := maxminddb.Open("../GeoLite2-City.mmdb")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	record := struct {
		Domain string `maxminddb:"connection_type"`
	}{}

	networks := db.Networks()
	for networks.Next() {
		subnet, err := networks.Network(&record)
		if err != nil {
			log.Fatal(err)
		}
		fmt.Printf("%s: %s\n", subnet.String(), record.Domain)
	}
	if networks.Err() != nil {
		log.Fatal(networks.Err())
	}
}
