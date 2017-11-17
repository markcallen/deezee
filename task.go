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
	"bufio"
	"bytes"
	"context"
	"fmt"
	"strings"

	// https://github.com/moby/moby/tree/master/client
	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/client"
)

// Task : docker image and command to return
type Task struct {
	Name        string
	Image       string
	Command     string
	Environment []string `yaml:",flow"`
	Volumes     []string `yaml:",flow"`
	Result      string
}

// Run : runs this task
func (t *Task) Run() string {
	fmt.Printf("\tTask: %s\n", t.Name)
	t.getImage()

	res := t.runImage()

	return res
}

func (t *Task) getImage() {
	if !strings.Contains(t.Image, ":") {
		t.Image = t.Image + ":latest"
	}

	cli, err := client.NewEnvClient()
	if err != nil {
		panic(err)
	}

	ctx := context.Background()

	_, err = cli.ImagePull(ctx, t.Image, types.ImagePullOptions{})
	if err != nil {
		panic(err)
	}

	images, err := cli.ImageList(ctx, types.ImageListOptions{})
	if err != nil {
		panic(err)
	}

	for _, image := range images {
		for _, tag := range image.RepoTags {
			if strings.Compare(tag, t.Image) == 0 {
				fmt.Printf("\t\tUsing image: %s\n", tag)
			}
		}
	}
}

func (t *Task) runImage() string {
	//TODO: put this into a global variable
	cli, err := client.NewEnvClient()
	if err != nil {
		panic(err)
	}

	ctx := context.Background()

	containerConfig := &container.Config{
		Image:        t.Image,
		Cmd:          strings.Fields(t.Command),
		Env:          t.Environment,
		Tty:          true,
		AttachStdout: true,
		AttachStderr: true,
	}

	hostConfig := &container.HostConfig{
		AutoRemove: false,
		Binds:      t.Volumes,
	}

	fmt.Printf("\t\tCreating Container\n")
	fmt.Printf("\t\t\tImage: %s\n", containerConfig.Image)
	fmt.Printf("\t\t\tCommand: %s\n", containerConfig.Cmd)
	fmt.Printf("\t\t\tEnvironment: %s\n", containerConfig.Env)
	fmt.Printf("\t\t\tBinds: %s\n", hostConfig.Binds)

	resp, err := cli.ContainerCreate(ctx, containerConfig, hostConfig, nil, t.Name)
	if err != nil {
		panic(err)
	}

	if err2 := cli.ContainerStart(ctx, resp.ID, types.ContainerStartOptions{}); err != nil {
		panic(err2)
	}
	fmt.Printf("\t\tStarting Container\n")

	if len(resp.Warnings) > 0 {
		fmt.Println("\t\t\tWarnings:", resp.Warnings)
	}

	fmt.Println("\t\tOutput")

	var buffer bytes.Buffer

	go func() {
		reader, err := cli.ContainerLogs(context.Background(), resp.ID, types.ContainerLogsOptions{
			ShowStdout: true,
			ShowStderr: true,
			Follow:     true,
			Timestamps: false,
		})
		if err != nil {
			panic(err)
		}
		defer reader.Close()

		scanner := bufio.NewScanner(reader)
		for scanner.Scan() {
			fmt.Printf("\t\t\t%s\n", scanner.Text())
			buffer.WriteString(scanner.Text())
		}
	}()

	statusCh, errCh := cli.ContainerWait(ctx, resp.ID, container.WaitConditionNotRunning)
	select {
	case err3 := <-errCh:
		if err3 != nil {
			panic(err3)
		}
	case <-statusCh:
	}

	fmt.Printf("\n")

	fmt.Printf("\t\tContainer Stopped\n")

	if err := cli.ContainerRemove(ctx, resp.ID, types.ContainerRemoveOptions{}); err != nil {
		panic(err)
	}

	res := stripCtlAndExtFromBytes(buffer.String())

	return res
}

func stripCtlAndExtFromBytes(str string) string {
	b := make([]byte, len(str))
	var bl int
	for i := 0; i < len(str); i++ {
		c := str[i]
		if c >= 32 && c < 127 {
			b[bl] = c
			bl++
		}
	}
	return string(b[:bl])
}
