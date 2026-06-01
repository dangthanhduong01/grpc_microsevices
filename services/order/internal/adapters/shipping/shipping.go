package shipping

import (
	pbshipping "github.com/dangthanhduong01/microservices_proto/pb/shipping"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type Adapter struct {
	pbshipping.ShippingClient
}

func NewAdapter(shippingServiceUrl string) (*Adapter, error) {
	var opts []grpc.DialOption
	opts = append(opts, grpc.WithTransportCredentials(insecure.NewCredentials()))
	conn, err := grpc.Dial(shippingServiceUrl, opts...)
	if err != nil {
		return nil, err
	}
	defer conn.Close()
	client := pbshipping.NewShippingClient(conn)
	return &Adapter{client}, nil
}

func (a *Adapter) Charge() error {
	return nil
}
