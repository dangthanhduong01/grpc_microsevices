package main

import (
	"log"
	"services/product/config"
	"services/product/internal/adapters/cache"
	"services/product/internal/adapters/db"
	"services/product/internal/adapters/grpc"
	"services/product/internal/applications/core/api"

	"github.com/joho/godotenv"
)

func main() {
	if err := godotenv.Load(".env"); err != nil {
		log.Println("No .env file found, using system environment variables")
	}
	dbAdapter, err := db.NewAdapter(config.GetDataSourceURL())
	if err != nil {
		log.Fatalf("failed to initialize db adapter: %v", err)
	}

	cacheAdapter, err := cache.NewAdapter(config.GetRedisURL())
	if err != nil {
		log.Fatalf("failed to initialize cache adapter: %v", err)
	}

	application := api.NewApplication(dbAdapter, cacheAdapter)

	grpcAdapter := grpc.NewAdapter(application, config.GetApplicationPort())
	grpcAdapter.Run()
}
