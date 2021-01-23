package dockertools

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"os"
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
	ctx := context.Background()
	host = fmt.Sprintf("tcp://%s", host)
	opt := client.WithHost(host)
	cli, err := client.NewClientWithOpts(opt, client.WithAPIVersionNegotiation())
	if err != nil {
		panic(err)
	}

	dock.ctx = ctx
	dock.cli = *cli
	return dock
}

func (dock Dock) Create(image string, memory int64, size string) string {
	// hostConfig docs https://docs.docker.com/engine/api/v1.24/
	hostConfig := new(container.HostConfig)
	hostConfig.Resources.Memory = memory << 20 // 限制内存
	hostConfig.Resources.CPUShares = 256
	hostConfig.StorageOpt = map[string]string{
		"size": size, // 限制磁盘 单位 M、G
	}
	// hostConfig.Resources.CPUPeriod = 100000
	// hostConfig.Resources.CPUQuota = 50000
	resp, err := dock.cli.ContainerCreate(
		dock.ctx,
		&container.Config{
			Image:        image,
			Tty:          true,
			User:         "root",
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

func (dock Dock) CreateAndStart(image string, memory int64, size string) string {
	containerID := dock.Create(image, memory, size)
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

func (dock Dock) ContainerRemove(containerID string) error {
	dock.Stop(containerID)
	return dock.cli.ContainerRemove(dock.ctx, containerID, types.ContainerRemoveOptions{})
}

func (dock Dock) ImageRemove(imageID string) ([]types.ImageDeleteResponseItem, error) {
	return dock.cli.ImageRemove(dock.ctx, imageID, types.ImageRemoveOptions{})
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

func (dock Dock) Commit(containerID string, imageName string) (types.IDResponse, error) {
	commitResp, err := dock.cli.ContainerCommit(dock.ctx, containerID, types.ContainerCommitOptions{Reference: imageName})
	return commitResp, err
}

func (dock Dock) Push(imageName string) error {
	auth := types.AuthConfig{
		Username: "foxsir",
		Password: "123456",
	}
	authBytes, _ := json.Marshal(auth)
	authBase64 := base64.URLEncoding.EncodeToString(authBytes)
	fmt.Println(authBase64)

	fmt.Println(fmt.Sprintf("%s:latest", imageName))

	pusher, err := dock.cli.ImagePush(dock.ctx, fmt.Sprintf("%s:latest", imageName), types.ImagePushOptions{
		All:           false,
		RegistryAuth:  authBase64,
		PrivilegeFunc: nil,
	})
	if err != nil {
		panic(err)
	}
	defer pusher.Close()
	io.Copy(os.Stdout, pusher)
	return err
}

func (dock Dock) Close() bool {
	dock.cli.Close()
	return true
}
