package grpc

import (
	"payment/services/order/internal/pb"
	"payment/services/order/internal/ports"
)

type Adapter struct {
	api  ports.APIPort
	port int
	pb.UnimplementedOrderServiceServer
}
