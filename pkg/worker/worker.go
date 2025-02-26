package worker

import (
	common "brigade/internal/gen/common/v1"
	"brigade/pkg/channel"
	"brigade/pkg/config"
	"brigade/pkg/plugins"
	"brigade/pkg/server"
	"context"
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	v1connect "brigade/internal/gen/cloud/v1/v1connect"

	"connectrpc.com/connect"
	"connectrpc.com/grpchealth"
	"connectrpc.com/grpcreflect"
	"github.com/rs/cors"
	"golang.org/x/net/http2"
	"golang.org/x/net/http2/h2c"
)

type Worker struct {
	pluginName string
	workerName string
}

func NewWorker(cfg *config.Config) *Worker {
	return &Worker{
		pluginName: cfg.PluginName,
		workerName: cfg.ChannelName,
	}
}

func (w *Worker) Start(ctx context.Context) error {
	channel := channel.NewChannel(w.workerName)
	stream, err := channel.Receive(context.Background())
	if err != nil {
		return err
	}
	for msg := range stream {
		plugin := plugins.NewPluginServer(w.pluginName)
		resp, err := plugin.Execute(ctx, &common.ExecuteRequest{
			Message: msg,
		})
		if err != nil {
			return err
		}
		fmt.Println(msg.JobId)
		fmt.Println(msg.PluginName)
		fmt.Println(msg.Message)
		fmt.Println(resp.Output)
	}
	ctx.Done()
	return nil
}

func (w *Worker) Server(ctx context.Context) error {
	mux := http.NewServeMux()
	mux.Handle(
		"/",
		http.RedirectHandler("https://connectrpc.com/demo", http.StatusFound),
	)
	compress1KB := connect.WithCompressMinBytes(1024)
	mux.Handle(v1connect.NewBrigadeCloudServiceHandler(
		server.NewServer(),
		compress1KB,
	))
	mux.Handle(grpchealth.NewHandler(
		grpchealth.NewStaticChecker(v1connect.BrigadeCloudServiceName),
		compress1KB,
	))
	mux.Handle(grpcreflect.NewHandlerV1(
		grpcreflect.NewStaticReflector(v1connect.BrigadeCloudServiceName),
		compress1KB,
	))
	mux.Handle(grpcreflect.NewHandlerV1Alpha(
		grpcreflect.NewStaticReflector(v1connect.BrigadeCloudServiceName),
		compress1KB,
	))
	addr := "localhost:8080"
	if port := os.Getenv("PORT"); port != "" {
		addr = ":" + port
	}
	srv := &http.Server{
		Addr: addr,
		Handler: h2c.NewHandler(
			newCORS().Handler(mux),
			&http2.Server{},
		),
		ReadHeaderTimeout: time.Second,
		ReadTimeout:       5 * time.Minute,
		WriteTimeout:      5 * time.Minute,
		MaxHeaderBytes:    8 * 1024, // 8KiB
	}
	signals := make(chan os.Signal, 1)
	signal.Notify(signals, os.Interrupt, syscall.SIGTERM)
	go func() {
		if err := srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Fatalf("HTTP listen and serve: %v", err)
		}
	}()

	<-signals
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		log.Fatalf("HTTP shutdown: %v", err) //nolint:gocritic
	}
	return nil
}

func newCORS() *cors.Cors {
	// To let web developers play with the demo service from browsers, we need a
	// very permissive CORS setup.
	return cors.New(cors.Options{
		AllowedMethods: []string{
			http.MethodHead,
			http.MethodGet,
			http.MethodPost,
			http.MethodPut,
			http.MethodPatch,
			http.MethodDelete,
		},
		AllowOriginFunc: func(_ /* origin */ string) bool {
			// Allow all origins, which effectively disables CORS.
			return true
		},
		AllowedHeaders: []string{"*"},
		ExposedHeaders: []string{
			// Content-Type is in the default safelist.
			"Accept",
			"Accept-Encoding",
			"Accept-Post",
			"Connect-Accept-Encoding",
			"Connect-Content-Encoding",
			"Content-Encoding",
			"Grpc-Accept-Encoding",
			"Grpc-Encoding",
			"Grpc-Message",
			"Grpc-Status",
			"Grpc-Status-Details-Bin",
		},
		// Let browsers cache CORS information for longer, which reduces the number
		// of preflight requests. Any changes to ExposedHeaders won't take effect
		// until the cached data expires. FF caps this value at 24h, and modern
		// Chrome caps it at 2h.
		MaxAge: int(2 * time.Hour / time.Second),
	})
}
