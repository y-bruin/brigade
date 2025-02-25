package k8s

import (
	"context"

	common "brigade/internal/gen/common/v1"
)

const (
	PluginName = "docker"
	Version    = "v1"
)

type DockerPluginServer struct{}

func (e *DockerPluginServer) Execute(
	ctx context.Context,
	request *common.ExecuteRequest,
) (*common.ExecuteResponse, error) {
	return &common.ExecuteResponse{}, nil
}

func (e *DockerPluginServer) Logs(
	ctx context.Context,
	request *common.LogsRequest,
) (*common.LogsResponse, error) {
	return &common.LogsResponse{}, nil
}

func (e *DockerPluginServer) Status(
	ctx context.Context,
	request *common.StatusRequest,
) (*common.StatusResponse, error) {
	return &common.StatusResponse{}, nil
}

func (e *DockerPluginServer) Events(
	ctx context.Context,
	request *common.EventsRequest,
) (*common.EventsResponse, error) {
	return &common.EventsResponse{}, nil
}
