package main

import (
	"fmt"
	"strings"
)

// Process : groups of tasks that need to be run
type Process struct {
	Name  string
	Tasks []Task
}

// Run : runs the tasks in the Process
func (p *Process) Run(g map[string]string) {
	fmt.Printf("Process: %s\n", p.Name)
	for _, task := range p.Tasks {

		vars := make([]string, 0, len(g)*2)
		for v := range g {
			vars = append(vars, "${"+v+"}", g[v])
		}

		replacer := strings.NewReplacer(vars...)
		task.Command = replacer.Replace(task.Command)

		if task.Result != "" {
			_, ok := g[task.Result]
			if !ok {
				fmt.Printf("No global variable %s to set\n", task.Result)
			}
			g[task.Result] = task.Run()
		} else {
			task.Run()
		}
	}
}
