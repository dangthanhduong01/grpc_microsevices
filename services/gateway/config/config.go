package config

import (
	"log"
	"os"
	"strconv"
)

func GetEnv() string {
	return getEnvironmentVariable("ENV")
}

func GetDataSourceURL() string {
	return getEnvironmentVariable("DATA_SOURCE_URL")
}

func GetApplicationPort() int {
	portStr := getEnvironmentVariable("APPLICATION_PORT")
	port, err := strconv.Atoi(portStr)
	if err != nil {
		log.Fatalf("environment variable APPLICATION_PORT is not a valid int: %v", err)
	}
	return port
}

func GetPaymentServiceURL() string {
	return getEnvironmentVariable("PAYMENT_SERVICE_URL")
}

func GetOrderServiceURL() string {
	return getEnvironmentVariable("ORDER_SERVICE_URL")
}

func GetAuthServiceURL() string {
	return getEnvironmentVariable("AUTH_SERVICE_URL")
}

func GetRateLimitRequests() int {
	requestsStr := os.Getenv("RATE_LIMIT_REQUESTS")
	if requestsStr == "" {
		return 100 // default
	}
	requests, err := strconv.Atoi(requestsStr)
	if err != nil {
		return 100
	}
	return requests
}

func GetRateLimitWindow() int {
	windowStr := os.Getenv("RATE_LIMIT_WINDOW_SECONDS")
	if windowStr == "" {
		return 60 // default 60 seconds
	}
	window, err := strconv.Atoi(windowStr)
	if err != nil {
		return 60
	}
	return window
}

func getEnvironmentVariable(key string) string {
	if os.Getenv(key) == "" {
		log.Fatalf("environment variable %s is not set", key)
	}
	return os.Getenv(key)
}