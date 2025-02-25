package cloud

import (
	"context"
	"errors"
	"fmt"
	"io"
	"time"

	"connectrpc.com/connect"

	cloudv1 "brigade/internal/gen/cloud/v1"

	cloudv1connect "brigade/internal/gen/cloud/v1/v1connect"
)

type BrigadeCloudServer struct {
	streamDelay time.Duration // sleep between streaming response messages
}

// NewBrigadeCloudServer returns a new cloud implementation which sleeps for the
// provided duration between streaming responses.
func NewBrigadeCloudServer(streamDelay time.Duration) cloudv1connect.BrigadeCloudServiceHandler {
	return &BrigadeCloudServer{streamDelay: streamDelay}
}

func (e *BrigadeCloudServer) Heartbeat(
	ctx context.Context,
	request *connect.Request[cloudv1.HeartbeatRequest],
) (*connect.Response[cloudv1.HeartbeatResponse], error) {
	return connect.NewResponse(&cloudv1.HeartbeatResponse{}), nil
}

func (e *BrigadeCloudServer) Events(
	ctx context.Context,
	request *connect.Request[cloudv1.EventsRequest],
) (*connect.Response[cloudv1.EventsResponse], error) {
	return connect.NewResponse(&cloudv1.EventsResponse{}), nil
}

func (e *BrigadeCloudServer) Logs(
	ctx context.Context,
	request *connect.Request[cloudv1.LogsRequest],
) (*connect.Response[cloudv1.LogsResponse], error) {
	return connect.NewResponse(&cloudv1.LogsResponse{}), nil
}

func (e *BrigadeCloudServer) Subscribe(
	ctx context.Context,
	stream *connect.BidiStream[cloudv1.SubscribeRequest, cloudv1.SubscribeResponse],
) error {
	for {
		if err := ctx.Err(); err != nil {
			return err
		}
		request, err := stream.Receive()
		if err != nil && errors.Is(err, io.EOF) {
			return nil
		} else if err != nil {
			return fmt.Errorf("receive request: %w", err)
		}
		fmt.Println(request.RequestType.Type)
		if err := stream.Send(&cloudv1.SubscribeResponse{Sentence: request.RequestType.String()}); err != nil {
			return fmt.Errorf("send response: %w", err)
		}
	}
}
