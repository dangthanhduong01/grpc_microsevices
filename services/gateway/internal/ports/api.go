package ports

import (
	"context"

	pb "github.com/dangthanhduong01/microservices_proto/pb/auth"
)

type AuthPort interface {
	RegisterUser(ctx context.Context, req *pb.CreateUserRequest) (*pb.CreateUserResponse, error)
	LoginUser(ctx context.Context, req *pb.LoginUserRequest) (*pb.LoginUserResponse, error)
	UpdateUser(ctx context.Context, req *pb.UpdateUserRequest) (*pb.UpdateUserResponse, error)
	ValidateToken(ctx context.Context, token string) (string, error)
}

type RateLimiterPort interface {
	Allow(ctx context.Context, key string) (bool, error)
}

type OrderPort interface {
	PlaceOrder(ctx context.Context, request interface{}) (interface{}, error)
}