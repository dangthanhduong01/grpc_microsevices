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

func GetJWTSecret() string {
	return getEnvironmentVariable("JWT_SECRET")
}

func GetAccessTokenExpiry() int {
	expiryStr := getEnvironmentVariable("ACCESS_TOKEN_EXPIRY")
	expiry, err := strconv.Atoi(expiryStr)
	if err != nil {
		log.Fatalf("environment variable ACCESS_TOKEN_EXPIRY is not a valid int: %v", err)
	}
	return expiry
}

func GetRefreshTokenExpiry() int {
	expiryStr := getEnvironmentVariable("REFRESH_TOKEN_EXPIRY")
	expiry, err := strconv.Atoi(expiryStr)
	if err != nil {
		log.Fatalf("environment variable REFRESH_TOKEN_EXPIRY is not a valid int: %v", err)
	}
	return expiry
}

func getEnvironmentVariable(key string) string {
	if os.Getenv(key) == "" {
		log.Fatalf("environment variable %s is not set", key)
	}
	return os.Getenv(key)
}