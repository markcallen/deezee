package main

import "fmt"

// Process : groups of tasks that need to be run
type Process struct {
	Name  string
	Tasks []Task
}

// Run : runs the tasks in the Process
func (p *Process) Run(g []string) {
	fmt.Printf("Process: %s\n", p.Name)
	for _, task := range p.Tasks {
		task.Run()
	}
}
