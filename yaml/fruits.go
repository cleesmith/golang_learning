package main

import (
	"fmt"
	"io/ioutil"
	"path/filepath"

	"gopkg.in/yaml.v2"
)

type Config struct {
  Description string // `yaml:"description"`
  Actions map[string][]string // `yaml:"act_ions"`
}

func main() {
	// quick example:
	// type T struct {
	// 	F int `yaml:"fff,omitempty"`
	// 	B int
	// }
	// var t T
	// yaml.Unmarshal([]byte("fff: 1\nb: 2"), &t)
	// fmt.Printf("t=%T=%#v\n", t, t)

	var config Config

	filename, _ := filepath.Abs("./fruits.yml")
	yamlFile, err := ioutil.ReadFile(filename)
	if err != nil {
		panic(err)
	}

	// err = yaml.Unmarshal(yamlFile, &config)
	// if err != nil {
	// 	panic(err)
	// }

	if err := config.Parse(yamlFile); err != nil {
  	panic(err)
  }
  fmt.Printf("\nconfig=%+v\n\n", config)
	fmt.Printf("config.Description=%T=%#v\n", config.Description, config.Description)
	fmt.Printf("config.Actions=%T=%#v\n", config.Actions, config.Actions)
  for k, v := range config.Actions {
    fmt.Printf("\nkey=%T=%v\n", k, k)
    fmt.Printf("value=%T=%v\n", v, v)
    for k2, v2 := range v {
	    fmt.Printf("\tk2=%T=%v\n", k2, k2)
	    fmt.Printf("\tv2=%T=%v\n", v2, v2)
    }
  }
}

func (c *Config) Parse(data []byte) error {
	return yaml.Unmarshal(data, c)
}
