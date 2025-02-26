package main

import (
	"context"
	"log"
	"net/http"

	"github.com/spf13/pflag"

	v1 "brigade/internal/gen/cloud/v1"
	cloudv1connect "brigade/internal/gen/cloud/v1/v1connect"
	common "brigade/internal/gen/common/v1"
	plugins "brigade/pkg/plugins"
)

func main() {
	helpArg := pflag.BoolP("help", "h", false, "")
	streamDelayArg := pflag.DurationP(
		"server-stream-delay",
		"d",
		0,
		"The duration to delay sending responses on the server stream.",
	)
	pflag.Parse()
	inter := "docker"

	if *helpArg {
		pflag.PrintDefaults()
		return
	}

	if *streamDelayArg < 0 {
		log.Printf("Server stream delay cannot be negative.")
		return
	}
	client := cloudv1connect.NewBrigadeCloudServiceClient(
		http.DefaultClient,
		"http://localhost:8080",
	)

	stream := client.Subscribe(context.Background())

	for {
		msg, err := stream.Receive()
		if err != nil {
			log.Printf("Error receiving message: %v", err)
			return
		}

		if msg.RequestType.Type == v1.RequestType_TASK {

			plugin := plugins.NewPluginServer(inter)
			resp, err := plugin.Execute(context.Background(), &common.ExecuteRequest{})
			if err != nil {
				log.Printf("Error executing task: %v", err)
				return
			}
			log.Printf("Received message: %v", resp.Output)
		}
	}
}
