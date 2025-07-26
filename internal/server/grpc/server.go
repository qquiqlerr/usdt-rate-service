package grpc

import (
	"context"
	"net"
	grpcHandler "usdt-rate-service/internal/handler/grpc"
	"usdt-rate-service/internal/pb"

	"google.golang.org/grpc"
)

// Server represents the gRPC server for the USDT rate service.
type Server struct {
	grpcServer   *grpc.Server
	ratesHandler *grpcHandler.RatesHandler
}

// NewServer creates a new gRPC server with the provided RatesHandler.
// It registers the RatesService with the gRPC server.
func NewServer(ratesHandler *grpcHandler.RatesHandler) *Server {
	grpcServer := grpc.NewServer()

	pb.RegisterRatesServiceServer(grpcServer, ratesHandler)

	return &Server{
		grpcServer:   grpcServer,
		ratesHandler: ratesHandler,
	}
}

// Start starts the gRPC server on the specified address.
func (s *Server) Start(ctx context.Context, addr string) error {
	config := &net.ListenConfig{}
	listener, err := config.Listen(ctx, "tcp", addr)
	if err != nil {
		return err
	}
	return s.grpcServer.Serve(listener)
}

// Stop gracefully stops the gRPC server.
func (s *Server) Stop() {
	s.grpcServer.GracefulStop()
}
