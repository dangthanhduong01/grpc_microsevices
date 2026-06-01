package order

import (
	"context"

	pb "github.com/dangthanhduong01/microservices_proto/pb/order"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type Adapter struct {
	client pb.OrderServiceClient
}

func NewAdapter(orderServiceURL string) (*Adapter, error) {
	conn, err := grpc.NewClient(orderServiceURL, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, err
	}

	client := pb.NewOrderServiceClient(conn)
	return &Adapter{client: client}, nil
}

func (a *Adapter) CreateOrder(ctx context.Context, req *pb.CreateOrderRequest) (*pb.CreateOrderResponse, error) {
	return a.client.CreateOrder(ctx, req)
}