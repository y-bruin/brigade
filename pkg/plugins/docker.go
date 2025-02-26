package plugins

import (
	"context"
	"fmt"
	"os"

	common "brigade/internal/gen/common/v1"

	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/client"
	"github.com/docker/docker/pkg/stdcopy"
)

type DockerPluginServer struct {
	client *client.Client
}

func NewDockerPluginServer() PluginServerInterface {
	cli, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
	if err != nil {
		panic(err)
	}
	return &DockerPluginServer{
		client: cli,
	}
}

func (e *DockerPluginServer) Execute(
	ctx context.Context,
	request *common.ExecuteRequest,
) (*common.ExecuteResponse, error) {
	container, err := e.client.ContainerCreate(ctx, &container.Config{
		Image: "ubuntu:latest",
	}, nil, nil, nil, "test-container")
	if err != nil {
		return nil, err
	}
	fmt.Println(container.ID)
	return &common.ExecuteResponse{}, nil
}

func (e *DockerPluginServer) Logs(
	ctx context.Context,
	request *common.LogsRequest,
) (*common.LogsResponse, error) {
	reader, err := e.client.ContainerLogs(ctx, "test-container", container.LogsOptions{
		ShowStdout: true,
		ShowStderr: true,
	})
	if err != nil {
		return nil, err
	}
	defer reader.Close()

	stdcopy.StdCopy(os.Stdout, os.Stderr, reader)

	return &common.LogsResponse{}, nil
}

func (e *DockerPluginServer) Status(
	ctx context.Context,
	request *common.StatusRequest,
) (*common.StatusResponse, error) {
	container, err := e.client.ContainerInspect(ctx, "test-container")
	if err != nil {
		return nil, err
	}

	return &common.StatusResponse{
		Status: container.State.Status,
	}, nil
}

func (e *DockerPluginServer) Events(
	ctx context.Context,
	request *common.EventsRequest,
) (*common.EventsResponse, error) {

	return &common.EventsResponse{}, nil
}
