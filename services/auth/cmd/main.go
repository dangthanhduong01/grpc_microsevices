package main

import (
	"log"
	"services/auth/config"
	"services/auth/internal/adapters/db"
	"services/auth/internal/adapters/grpc"
	"services/auth/internal/applications/core/api"
)

func main() {
	dbAdapter, err := db.NewAdapter(config.GetDataSourceURL())
	if err != nil {
		log.Fatalf("failed to initialize db adapter: %v", err)
	}

	application := api.NewApplication(dbAdapter)

	grpcAdapter := grpc.NewAdapter(application, config.GetApplicationPort())
	grpcAdapter.Run()
}