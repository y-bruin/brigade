package server

import (
	v1 "brigade/internal/gen/cloud/v1"
	"brigade/internal/gen/cloud/v1/v1connect"
	"context"

	"connectrpc.com/connect"
)

type Server struct {
}

func NewServer() v1connect.BrigadeCloudServiceHandler {
	return &Server{}
}

func (s *Server) Heartbeat(ctx context.Context, req *connect.Request[v1.HeartbeatRequest]) (*connect.Response[v1.HeartbeatResponse], error) {
	return &connect.Response[v1.HeartbeatResponse]{
		Msg: &v1.HeartbeatResponse{
			Payload: &v1.Payload{
				Sentence: "Hello, world!",
			},
		},
	}, nil
}

func (s *Server) Events(ctx context.Context, req *connect.Request[v1.EventsRequest]) (*connect.Response[v1.EventsResponse], error) {
	return &connect.Response[v1.EventsResponse]{
		Msg: &v1.EventsResponse{},
	}, nil
}

func (s *Server) Logs(ctx context.Context, req *connect.Request[v1.LogsRequest]) (*connect.Response[v1.LogsResponse], error) {
	return &connect.Response[v1.LogsResponse]{
		Msg: &v1.LogsResponse{},
	}, nil
}
