/*
Copyright 2017 Mark C Allen <mark@markcallen.com>

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at
     http://www.apache.org/licenses/LICENSE-2.0
Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

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
