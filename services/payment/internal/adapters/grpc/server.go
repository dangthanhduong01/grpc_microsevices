package grpc

import (
	"services/payment/internal/ports"

	"github.com/dangthanhduong01/microservices_proto/pb/payment"
	"google.golang.org/grpc"
)

type Adapter struct {
	api    ports.APIPort
	port   int
	server *grpc.Server
	payment.UnimplementedPaymentServer
}

func NewAdapter(api ports.APIPort, port int) *Adapter {
	return &Adapter{
		api:  api,
		port: port,
	}
}

// func (a *Adapter) Run() {
// 	var err error

// 	listen, err
// }
