package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"html/template"
	"os"
	"regexp"
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

func replaceGlobalJSONVariables(g map[string]string, str string) string {
	m := map[string]interface{}{}
	for v := range g {
		re := regexp.MustCompile(`\$\{` + v + `(\.[\S]+)\}`)
		s := re.ReplaceAllString(str, `{{$1}}`)

		t := template.Must(template.New("").Parse(s))

		if err := json.Unmarshal([]byte(g[v]), &m); err != nil {
			continue
		}

		var tpl bytes.Buffer
		if err := t.Execute(&tpl, m); err != nil {
			panic(err)
		}
		str = tpl.String()
	}

	return str

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
		task.Command = replaceGlobalJSONVariables(g, task.Command)
		task.Command = replaceEnvVariables(task.Command)

		for i := range task.Environment {
			task.Environment[i] = replaceGlobalVariables(g, task.Environment[i])
			task.Environment[i] = replaceGlobalJSONVariables(g, task.Environment[i])
			task.Environment[i] = replaceEnvVariables(task.Environment[i])
		}

		for i := range task.Volumes {
			task.Volumes[i] = replaceGlobalVariables(g, task.Volumes[i])
			task.Volumes[i] = replaceGlobalJSONVariables(g, task.Volumes[i])
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
