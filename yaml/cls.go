package main

import (
	"fmt"
	"io/ioutil"
	"log"

	"gopkg.in/yaml.v2"
)

type ConfigType struct {
	keys map[string]string "yaml: keys"
}

func main() {
	// config := make(map[interface{}]interface{})
	var config ConfigType
	filename := "xls_config.yml"
	yamlFile, err := ioutil.ReadFile(filename)
	if err != nil {
  	panic(err)
  }
	err = yaml.Unmarshal(yamlFile, &config)
	// err = yaml.UnmarshalYAML(yamlFile, &config)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("config type is %T\n", config)
	fmt.Printf("--- config:\n%v\n\n", config.keys)
	// fmt.Printf("\nconfig.urls=%v\n", config.urls)
}

// func (v *ConfigType) UnmarshalYAML(unmarshal func(interface{}) error) error {
// 	return unmarshal(&v.keys)
// }
