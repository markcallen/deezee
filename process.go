package main

import (
	"fmt"
	"os"
	"strings"
)

// Process : groups of tasks that need to be run
type Process struct {
	Name  string
	Tasks []Task
}

func replaceGlobalVariables(g map[string]string, str string) string {
	vars := make([]string, 0, len(g)*2)
	for v := range g {
		vars = append(vars, "${"+v+"}", g[v])
	}

	replacer := strings.NewReplacer(vars...)
	return replacer.Replace(str)

}

func replaceEnvVariables(str string) string {
	vars := make([]string, 0, len(os.Environ())*2)

	for _, e := range os.Environ() {
		pair := strings.Split(e, "=")
		vars = append(vars, "$env{"+pair[0]+"}", pair[1])
	}

	replacer := strings.NewReplacer(vars...)
	return replacer.Replace(str)
}

// Run : runs the tasks in the Process
func (p *Process) Run(g map[string]string) {
	fmt.Printf("Process: %s\n", p.Name)
	for _, task := range p.Tasks {
		printGlobals(g)

		task.Command = replaceGlobalVariables(g, task.Command)
		task.Command = replaceEnvVariables(task.Command)

		for i := range task.Environment {
			task.Environment[i] = replaceGlobalVariables(g, task.Environment[i])
			task.Environment[i] = replaceEnvVariables(task.Environment[i])
		}

		for i := range task.Volumes {
			task.Volumes[i] = replaceGlobalVariables(g, task.Volumes[i])
			task.Volumes[i] = replaceEnvVariables(task.Volumes[i])
		}

		if task.Result != "" {
			_, ok := g[task.Result]
			if !ok {
				fmt.Printf("No global variable %s to set\n", task.Result)
			}
			g[task.Result] = task.Run()
		} else {
			task.Run()
		}
		fmt.Printf("\n")
	}
}

func printGlobals(g map[string]string) {
	fmt.Printf("\tGlobals:\n")
	for v := range g {
		fmt.Printf("\t\t%s: %s\n", v, g[v])
	}
}
