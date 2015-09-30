package main

import (
	"flag"
	"fmt"
	"os"
)

var conf = flag.String("conf", "seo.toml", "a config toml file is required containing your SharedCount API key, etc.")

func main() {
	flag.Usage = usage
	flag.Parse()
	tomlFile := *conf
	if len(tomlFile) < 1 {
		fmt.Println("A toml config file is required!")
		usage()
	}
	fmt.Printf("main: conf flag: tomlFile=%T=%+v\n", tomlFile, tomlFile)

	// used for trailing positional arguments (unnamed):
	// args := flag.Args()
	// fmt.Printf("args=%T=%v\n", args, args)
	// if len(args) < 1 {
	// 	fmt.Println("A toml config file is required!")
	// 	usage()
	// 	os.Exit(1)
	// }
}

func usage() {
	fmt.Fprintf(os.Stdout, "Usage: ./enseomble -conf=whatever.toml\n")
	flag.PrintDefaults()
	os.Exit(2)
}
