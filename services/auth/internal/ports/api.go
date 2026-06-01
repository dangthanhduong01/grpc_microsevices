package ports

import (
	"context"

	"services/auth/config"

	pb "github.com/dangthanhduong01/microservices_proto/pb/auth"
)

type APIPort interface {
	RegisterUser(ctx context.Context, req *pb.CreateUserRequest) (*pb.CreateUserResponse, error)
	LoginUser(ctx context.Context, req *pb.LoginUserRequest) (*pb.LoginUserResponse, error)
	UpdateUser(ctx context.Context, req *pb.UpdateUserRequest) (*pb.UpdateUserResponse, error)
}

type DBPort interface {
	CreateUser(email, username, password string) (string, error)
	GetUserByEmail(email string) (string, error)
	GetUserByUsername(username string) (string, error)
	GetUserByID(id string) (string, error)
	UpdateUser(id, email, username, fullName string) error
}

func GetJWTSecret() string {
	return config.GetJWTSecret()
}

func GetAccessTokenExpiry() int {
	return config.GetAccessTokenExpiry()
}

func GetRefreshTokenExpiry() int {
	return config.GetRefreshTokenExpiry()
}