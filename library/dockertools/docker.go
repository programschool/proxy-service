package dockertools

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/strslice"
	"github.com/docker/docker/client"
)

type Dock struct {
	ctx context.Context
	cli client.Client
}

func (dock Dock) New(host string) Dock {
	var httpClient *http.Client
	host = fmt.Sprintf("tcp://%s", host)
	var clientVersion = "v1.40"
	defaultHeaders := map[string]string{"User-Agent": "engine-api-cli-1.0"}

	cli, err := client.NewClient(host, clientVersion, httpClient, defaultHeaders) // engine 1.13.1
	ctx := context.Background()
	if err != nil {
		panic(err)
	}

	dock.ctx = ctx
	dock.cli = *cli
	return dock
}

func (dock Dock) Create(image string) string {
	hostConfig := new(container.HostConfig)
	hostConfig.Resources.Memory = 150 << 20
	hostConfig.Resources.CPUShares = 256
	// hostConfig.Resources.CPUPeriod = 100000
	// hostConfig.Resources.CPUQuota = 50000
	resp, err := dock.cli.ContainerCreate(
		dock.ctx,
		&container.Config{
			Image:        image,
			Tty:          true,
			User:         "root",
			Cmd:          strslice.StrSlice([]string{"/bin/bash"}),
			AttachStdin:  true,
			AttachStdout: true,
			AttachStderr: true,
			OpenStdin:    true,
		}, hostConfig, nil, nil, "")

	if err != nil {
		panic(err)
	}

	return resp.ID
}

func (dock Dock) CreateAndStart(image string) string {
	containerID := dock.Create(image)
	dock.Start(containerID)

	return containerID

}

func (dock Dock) Start(containerID string) bool {
	err := dock.cli.ContainerStart(dock.ctx, containerID, types.ContainerStartOptions{})
	if err != nil {
		panic(err)
	}
	return true
}

func (dock Dock) ReStart(containerID string) bool {
	timeout := 0 * time.Second
	err := dock.cli.ContainerRestart(dock.ctx, containerID, &timeout)
	if err != nil {
		panic(err)
	}
	return true
}

func (dock Dock) Stop(containerID string) bool {
	timeout := 0 * time.Second
	err := dock.cli.ContainerStop(dock.ctx, containerID, &timeout)
	if err != nil {
		panic(err)
	}
	return true
}

func (dock Dock) Remove(containerID string) bool {
	dock.Stop(containerID)
	dock.cli.ContainerRemove(dock.ctx, containerID, types.ContainerRemoveOptions{})
	return true
}

func (dock Dock) Inspect(containerID string) types.ContainerJSON {
	inspect, err := dock.cli.ContainerInspect(dock.ctx, containerID)
	if err != nil {
		panic(err)
	}
	return inspect
}

func (dock Dock) ExecCommand(containerID string, command []string) (types.ContainerExecInspect, error) {
	exec, _ := dock.cli.ContainerExecCreate(dock.ctx, containerID, types.ExecConfig{
		User:         "root",
		Privileged:   false,
		Tty:          true,
		AttachStdin:  true,
		AttachStderr: true,
		AttachStdout: true,
		Cmd:          strslice.StrSlice(command),
	})

	err := dock.cli.ContainerExecStart(dock.ctx, exec.ID, types.ExecStartCheck{
		Detach: true,
		Tty:    false,
	})

	if err != nil {
		panic(err)
	}

	/*
		res, _ := dock.cli.ContainerExecAttach(dock.ctx, exec.ID, types.ExecStartCheck{})

		bs := bufio.NewScanner(res.Reader)

		for k := 0; bs.Scan(); k++ {
			fmt.Printf("%s %v\n", bs.Bytes(), bs.Text())
		}
	*/

	return dock.cli.ContainerExecInspect(dock.ctx, exec.ID)

}

// get all image
func (dock Dock) ImageList() []types.ImageSummary {
	images, err := dock.cli.ImageList(dock.ctx, types.ImageListOptions{})
	if err != nil {
		panic(err)
	}
	return images
}

func (dock Dock) IpAddress(containerID string) string {
	inspect := dock.Inspect(containerID)
	return inspect.NetworkSettings.DefaultNetworkSettings.IPAddress
}

func (dock Dock) IsRunning(containerID string) bool {
	inspect := dock.Inspect(containerID)
	if inspect.State.Status == "running" {
		return true
	}
	return false
}

func (dock Dock) Close() bool {
	dock.cli.Close()
	return true
}
