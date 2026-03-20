package main

import (
	"log"
	"services/order/config"
	"services/order/internal/adapters/db"
	"services/order/internal/adapters/grpc"
	"services/order/internal/adapters/payment"
	"services/order/internal/applications/core/api"
)

func main() {
	dbAdapter, err := db.NewAdapter(config.GetDataSourceURL())
	if err != nil {
		log.Fatalf("failed to initialize db adapter: %v", err)
	}

	paymentAdapter, err := payment.NewAdapter(config.GetPaymentServiceURL())
	if err != nil {
		log.Fatalf("failed to initialize payment adapter: %v", err)
	}

	application := api.NewApplication(dbAdapter, paymentAdapter)

	grpcAdapter := grpc.NewAdapter(application, config.GetApplicationPort())
	grpcAdapter.Run()
}
