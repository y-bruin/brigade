module brigade

go 1.22

toolchain go1.22.9

require (
	connectrpc.com/connect v1.14.0
	connectrpc.com/grpchealth v1.2.0
	connectrpc.com/grpcreflect v1.2.0
	github.com/planetscale/vtprotobuf v0.6.1-0.20240319094008-0393e58bdf10
	github.com/rs/cors v1.10.0
	github.com/spf13/pflag v1.0.5
	golang.org/x/net v0.32.0
	google.golang.org/grpc v1.70.0
	google.golang.org/protobuf v1.35.2
)

require (
	github.com/docker/docker v28.0.0+incompatible // indirect
	golang.org/x/sys v0.28.0 // indirect
	golang.org/x/text v0.21.0 // indirect
	google.golang.org/genproto/googleapis/rpc v0.0.0-20241202173237-19429a94021a // indirect
)
