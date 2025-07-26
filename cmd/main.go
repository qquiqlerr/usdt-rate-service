package main

import (
	"context"
	"os"
	"os/signal"
	"syscall"
	"usdt-rate-service/config"
	"usdt-rate-service/internal/adapter"
	handler "usdt-rate-service/internal/handler/grpc"
	"usdt-rate-service/internal/infra/grinex"
	"usdt-rate-service/internal/repository"
	server "usdt-rate-service/internal/server/grpc"
	"usdt-rate-service/internal/service"
	"usdt-rate-service/pkg/database"
	"usdt-rate-service/pkg/logger"

	"go.uber.org/zap"
)

func main() {
	config := config.MustLoad()
	logger := logger.MustLoadLogger(config.LogLevel)

	logger.Info("config loaded", zap.String("logLevel", config.LogLevel))
	logger.Debug("config", zap.Any("config", config))

	ctx, cancel := context.WithCancel(context.Background())

	pgPool, err := database.NewPostgresPool(ctx, config.DatabaseAddress)
	if err != nil {
		logger.Fatal("failed to create postgres pool", zap.Error(err))
	}
	defer pgPool.Close()
	logger.Info("Postgres pool created", zap.String("address", config.DatabaseAddress))

	ratesRepo := repository.NewRates(pgPool)

	grinex, err := grinex.NewClient(config.GrinexAddress)
	if err != nil {
		logger.Fatal("failed to create grinex client", zap.Error(err))
	}

	depthProvider := adapter.NewGrinexDepthProvider(grinex)

	service := service.NewRatesService(logger, depthProvider, ratesRepo)

	ratesHandler := handler.NewRatesHandler(service)

	grpcServer := server.NewServer(ratesHandler)

	go func() {
		if err = grpcServer.Start(ctx, config.GRPCAddress); err != nil {
			logger.Fatal("failed to start gRPC server", zap.Error(err))
		}
	}()
	// Wait for interrupt signal
	logger.Info("gRPC server started", zap.String("address", config.GRPCAddress))

	// Graceful shutdown
	ch := make(chan os.Signal, 1)
	signal.Notify(ch, os.Interrupt, syscall.SIGTERM)
	<-ch
	logger.Info("Shutting down")
	cancel()
	grpcServer.Stop()
	logger.Info("Shutdown complete")
}
