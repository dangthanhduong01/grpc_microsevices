package payment

import (
	"context"
	"services/order/internal/applications/core/domain"

	pbpayment "github.com/dangthanhduong01/microservices_proto/pb/payment"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type Adapter struct {
	payment pbpayment.PaymentClient
}

func NewAdapter(paymentServiceUrl string) (*Adapter, error) {
	var opts []grpc.DialOption
	opts = append(opts, grpc.WithTransportCredentials(insecure.NewCredentials()))
	conn, err := grpc.Dial(paymentServiceUrl, opts...)
	if err != nil {
		return nil, err
	}
	defer conn.Close()

	paymentClient := pbpayment.NewPaymentClient(conn)
	return &Adapter{payment: paymentClient}, nil
}

func (a *Adapter) Charge(order *domain.Order) error {
	_, err := a.payment.Create(context.Background(),
		&pbpayment.CreatePaymentRequest{
			UserId:     order.CustomerID,
			OrderId:    order.ID,
			TotalPrice: order.TotalPrice(),
		})
	return err
}
