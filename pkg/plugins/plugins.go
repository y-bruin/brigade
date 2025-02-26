package plugins

import (
	common "brigade/internal/gen/common/v1"
	"brigade/pkg/plugins/docker"
	"brigade/pkg/plugins/k8s"
	"context"
)

type PluginServerInterface interface {
	Execute(ctx context.Context, request *common.ExecuteRequest) (*common.ExecuteResponse, error)
	Logs(ctx context.Context, request *common.LogsRequest) (*common.LogsResponse, error)
	Status(ctx context.Context, request *common.StatusRequest) (*common.StatusResponse, error)
	Events(ctx context.Context, request *common.EventsRequest) (*common.EventsResponse, error)
}

func NewPluginServer(pluginName string) PluginServerInterface {
	switch pluginName {
	case "docker":
		return docker.NewDockerPluginServer()
	case "k8s":
		return k8s.NewK8sPluginServer()
	default:
		return nil
	}
}
