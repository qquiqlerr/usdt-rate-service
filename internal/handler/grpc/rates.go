package grpc

import (
	"context"
	"usdt-rate-service/internal/pb"
	"usdt-rate-service/internal/service"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// RatesHandler is a gRPC handler for managing rates.
type RatesHandler struct {
	pb.UnimplementedRatesServiceServer

	ratesService *service.RatesService
}

// NewRatesHandler creates a new RatesHandler with the provided logger and RatesService.
func NewRatesHandler(ratesService *service.RatesService) *RatesHandler {
	return &RatesHandler{
		ratesService: ratesService,
	}
}

// GetRates handles the gRPC request to get rates for a specific market.
func (h *RatesHandler) GetRates(ctx context.Context, req *pb.GetRatesRequest) (*pb.GetRatesResponse, error) {
	market := req.GetMarket()
	if market == "" {
		return nil, status.Error(codes.InvalidArgument, "market must be specified")
	}
	rate, err := h.ratesService.GetRates(ctx, market)
	if err != nil {
		return nil, status.Error(codes.Internal, "failed to get rates")
	}

	return &pb.GetRatesResponse{
		Rate: &pb.Rate{
			AskPrice:  rate.AskPrice,
			BidPrice:  rate.BidPrice,
			Timestamp: rate.Timestamp,
		},
	}, nil
}
