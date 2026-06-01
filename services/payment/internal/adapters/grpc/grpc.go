package grpc

import (
	"context"
	"services/payment/internal/applications/core/domain"

	pb "github.com/dangthanhduong01/microservices_proto/pb/payment"
)

func (a *Adapter) Create(ctx context.Context, req *pb.CreatePaymentRequest) (*pb.CreatePaymentResponse, error) {
	newPayment := domain.NewPayment(req.UserId, req.OrderId, req.TotalPrice)
	result, err := a.api.Charge(ctx, newPayment)
	if err != nil {
		return nil, err
	}
	return &pb.CreatePaymentResponse{PaymentId: result.ID}, nil
}
