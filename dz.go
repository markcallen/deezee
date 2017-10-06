package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"

	//https://github.com/go-yaml/yaml
	"gopkg.in/yaml.v2"

	// https://github.com/alecthomas/kingpin
	"gopkg.in/alecthomas/kingpin.v2"
)

type conf struct {
	Version   string
	Globals   map[string]string
	Processes []Process
}

var (
	app   = kingpin.New("dz", "Run tasks using Docker")
	debug = app.Flag("debug", "Enable debug mode.").Bool()

	run     = app.Command("run", "Run processes and tasks.")
	runFile = run.Flag("file", "yaml file.").Required().String()
)

func main() {
	switch kingpin.MustParse(app.Parse(os.Args[1:])) {

	// run
	case run.FullCommand():
		var dz conf
		dz.getConf(*runFile)
		dz.printConf()

		for _, process := range dz.Processes {
			process.Run(dz.Globals)
		}
	}

}

func (c *conf) getConf(filename string) *conf {
	yamlFile, err := ioutil.ReadFile(filename)
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
	fmt.Printf("yaml:\n%s\n\n", string(d))

}
