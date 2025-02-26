package k8s

import (
	"context"

	common "brigade/internal/gen/common/v1"

	"github.com/docker/docker/client"
)

const (
	PluginName = "k8s"
	Version    = "v1"
)

type K8sPluginServer struct {
	client *client.Client
}

func NewK8sPluginServer() *K8sPluginServer {
	cli, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
	if err != nil {
		panic(err)
	}

	return &K8sPluginServer{
		client: cli,
	}
}

func (e *K8sPluginServer) Execute(
	ctx context.Context,
	request *common.ExecuteRequest,
) (*common.ExecuteResponse, error) {
	return &common.ExecuteResponse{}, nil
}

func (e *K8sPluginServer) Logs(
	ctx context.Context,
	request *common.LogsRequest,
) (*common.LogsResponse, error) {
	return &common.LogsResponse{}, nil
}

func (e *K8sPluginServer) Status(
	ctx context.Context,
	request *common.StatusRequest,
) (*common.StatusResponse, error) {
	return &common.StatusResponse{}, nil
}

func (e *K8sPluginServer) Events(
	ctx context.Context,
	request *common.EventsRequest,
) (*common.EventsResponse, error) {
	return &common.EventsResponse{}, nil
}
