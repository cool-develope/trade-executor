package server

import (
	"context"
	"fmt"
	"net"
	"os"
	"os/signal"

	"github.com/cool-develope/trade-executor/internal/orderctrl/pb"
	"google.golang.org/grpc"
	"google.golang.org/grpc/health/grpc_health_v1"
)

// RunServer runs gRPC server and HTTP gateway
func RunServer(storage storageInterface, executorCtrl orderCtrlInterface, cfg Config) error {
	ctx := context.Background()

	if len(cfg.GRPCPort) == 0 {
		return fmt.Errorf("invalid TCP port for gRPC server: '%s'", cfg.GRPCPort)
	}

	executorService := NewExecutorService(storage, executorCtrl)

	go func() {
		_ = runGRPCServer(ctx, executorService, cfg.GRPCPort)
	}()

	return nil
}

// HealthChecker will provide an implementation of the HealthCheck interface.
type healthChecker struct{}

// NewHealthChecker returns a health checker according to standard package
// grpc.health.v1.
func newHealthChecker() *healthChecker {
	return &healthChecker{}
}

// HealthCheck interface implementation.

// Check returns the current status of the server for unary gRPC health requests,
// for now if the server is up and able to respond we will always return SERVING.
func (s *healthChecker) Check(ctx context.Context, req *grpc_health_v1.HealthCheckRequest) (*grpc_health_v1.HealthCheckResponse, error) {
	return &grpc_health_v1.HealthCheckResponse{
		Status: grpc_health_v1.HealthCheckResponse_SERVING,
	}, nil
}

// Watch returns the current status of the server for stream gRPC health requests,
// for now if the server is up and able to respond we will always return SERVING.
func (s *healthChecker) Watch(req *grpc_health_v1.HealthCheckRequest, server grpc_health_v1.Health_WatchServer) error {
	return server.Send(&grpc_health_v1.HealthCheckResponse{
		Status: grpc_health_v1.HealthCheckResponse_SERVING,
	})
}

func runGRPCServer(ctx context.Context, executorService pb.ExecutorServiceServer, port string) error {
	listen, err := net.Listen("tcp", ":"+port)
	if err != nil {
		return err
	}

	server := grpc.NewServer()
	pb.RegisterExecutorServiceServer(server, executorService)

	healthService := newHealthChecker()
	grpc_health_v1.RegisterHealthServer(server, healthService)

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	go func() {
		for range c {
			server.GracefulStop()
			<-ctx.Done()
		}
	}()

	fmt.Println("gRPC Server is serving at ", port)
	return server.Serve(listen)
}
