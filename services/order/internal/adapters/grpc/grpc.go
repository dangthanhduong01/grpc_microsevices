package grpc

import (
	"context"
	"services/order/internal/applications/core/domain"

	pb "github.com/dangthanhduong01/microservices_proto/pb/order"
)

func (a Adapter) Create(ctx context.Context, req *pb.CreateOrderRequest) (*pb.CreateOrderResponse, error) {
	var orderItems []domain.OrderItem
	for _, orderItem := range req.Items {
		orderItems = append(orderItems, domain.OrderItem{
			ProductCode: orderItem.ProductId,
			UnitPrice:   orderItem.UnitPrice,
			Quantity:    orderItem.Quantity,
		})
	}
	newOrder := domain.NewOrder(req.UserId, orderItems)
	rs, err := a.api.PlaceOrder(*newOrder)
	if err != nil {
		return nil, err
	}
	return &pb.CreateOrderResponse{
		OrderId: rs.ID,
		Message: "Order created successfully",
	}, nil
}
