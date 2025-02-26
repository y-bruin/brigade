package plugins

import (
	common "brigade/internal/gen/common/v1"
	"context"
	"fmt"
	"io"

	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
)

const (
	PluginName = "k8s"
	Version    = "v1"
)

type K8sPluginServer struct {
	client *kubernetes.Clientset
}

func NewK8sPluginServer() PluginServerInterface {
	// Create in-cluster config
	config, err := rest.InClusterConfig()
	if err != nil {
		panic(err)
	}

	// Create clientset
	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		panic(err)
	}

	return &K8sPluginServer{
		client: clientset,
	}
}

func (e *K8sPluginServer) Execute(
	ctx context.Context,
	request *common.ExecuteRequest,
) (*common.ExecuteResponse, error) {
	// Create pod specification
	pod := &corev1.Pod{
		ObjectMeta: metav1.ObjectMeta{
			GenerateName: "brigade-task-",
		},
		Spec: corev1.PodSpec{
			RestartPolicy: corev1.RestartPolicyNever,
			Containers: []corev1.Container{
				{
					Name: "task",
				},
			},
		},
	}

	// Create pod
	createdPod, err := e.client.CoreV1().Pods("default").Create(ctx, pod, metav1.CreateOptions{})
	if err != nil {
		return nil, fmt.Errorf("failed to create pod: %v", err)
	}

	return &common.ExecuteResponse{
		Output: createdPod.Name,
	}, nil
}

func (e *K8sPluginServer) Logs(
	ctx context.Context,
	request *common.LogsRequest,
) (*common.LogsResponse, error) {
	// Get pod logs
	podLogs, err := e.client.CoreV1().Pods("default").GetLogs(request.JobId, &corev1.PodLogOptions{}).Stream(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get pod logs: %v", err)
	}
	defer podLogs.Close()

	// Read logs
	logs, err := io.ReadAll(podLogs)
	if err != nil {
		return nil, fmt.Errorf("failed to read pod logs: %v", err)
	}

	return &common.LogsResponse{
		Log: string(logs),
	}, nil
}

func (e *K8sPluginServer) Status(
	ctx context.Context,
	request *common.StatusRequest,
) (*common.StatusResponse, error) {
	// Get pod status
	pod, err := e.client.CoreV1().Pods("default").Get(ctx, request.JobId, metav1.GetOptions{})
	if err != nil {
		return nil, fmt.Errorf("failed to get pod status: %v", err)
	}

	return &common.StatusResponse{
		Status: string(pod.Status.Phase),
	}, nil
}

func (e *K8sPluginServer) Events(
	ctx context.Context,
	request *common.EventsRequest,
) (*common.EventsResponse, error) {
	// Get pod events
	events, err := e.client.CoreV1().Events("default").List(ctx, metav1.ListOptions{
		FieldSelector: fmt.Sprintf("involvedObject.name=%s", request.JobId),
	})
	if err != nil {
		return nil, fmt.Errorf("failed to get pod events: %v", err)
	}

	// Format events
	var eventMessages string
	for _, event := range events.Items {
		eventMessages += fmt.Sprintf("[%s] %s: %s\n",
			event.LastTimestamp.String(),
			event.Reason,
			event.Message)
	}

	return &common.EventsResponse{
		Event: eventMessages,
	}, nil
}
