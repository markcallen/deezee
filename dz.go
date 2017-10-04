package main

import (
	"fmt"
	"io/ioutil"
	"log"

	//https://github.com/go-yaml/yaml
	"gopkg.in/yaml.v2"
)

type conf struct {
	Version   string
	Globals   map[string]string
	Processes []Process
}

func main() {
	var dz conf
	dz.getConf()
	dz.printConf()

	fmt.Printf("Globals: %v\n", dz.Globals)

	for _, process := range dz.Processes {
		process.Run(dz.Globals)
	}

	fmt.Printf("Globals: %v\n", dz.Globals)
}

func (c *conf) getConf() *conf {
	yamlFile, err := ioutil.ReadFile("sample.yml")
	if err != nil {
		log.Printf("yamlFile.Get err   #%v ", err)
	}
	err = yaml.Unmarshal(yamlFile, c)
	if err != nil {
		log.Fatalf("Unmarshal: %v", err)
	}

	return c
}

func (c *conf) printConf() {
	d, err := yaml.Marshal(&c)
	if err != nil {
		log.Fatalf("error: %v", err)
	}
	fmt.Printf("--- t dump:\n%s\n\n", string(d))

}
