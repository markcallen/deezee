package main

import (
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
		Image: t.Image,
		Cmd:   strings.Fields(t.Command),
	}

	hostConfig := &container.HostConfig{
		AutoRemove: false,
	}

	resp, err := cli.ContainerCreate(ctx, containerConfig, hostConfig, nil, t.Name)
	if err != nil {
		panic(err)
	}

	if err2 := cli.ContainerStart(ctx, resp.ID, types.ContainerStartOptions{}); err != nil {
		panic(err2)
	}

	statusCh, errCh := cli.ContainerWait(ctx, resp.ID, container.WaitConditionNotRunning)
	select {
	case err3 := <-errCh:
		if err3 != nil {
			panic(err3)
		}
	case <-statusCh:
	}

	out, err := cli.ContainerLogs(ctx, resp.ID, types.ContainerLogsOptions{ShowStdout: true})
	if err != nil {
		panic(err)
	}

	buf := new(bytes.Buffer)
	buf.ReadFrom(out)
	res := buf.String()

	if err := cli.ContainerRemove(ctx, resp.ID, types.ContainerRemoveOptions{}); err != nil {
		panic(err)
	}

	fmt.Printf("\t\tResult: %s", res)

	return res
}
