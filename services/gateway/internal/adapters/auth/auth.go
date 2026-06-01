package auth

import (
	"context"
	"errors"

	pb "github.com/dangthanhduong01/microservices_proto/pb/auth"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

var ErrAuthServiceNotAvailable = errors.New("auth service not available")

type Adapter struct {
	client pb.AuthServiceClient
	conn   *grpc.ClientConn
}

func NewAdapter(authServiceURL string) (*Adapter, error) {
	conn, err := grpc.NewClient(authServiceURL, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, err
	}

	client := pb.NewAuthServiceClient(conn)
	return &Adapter{client: client, conn: conn}, nil
}

func (a *Adapter) RegisterUser(ctx context.Context, req *pb.CreateUserRequest) (*pb.CreateUserResponse, error) {
	if a.client == nil {
		return nil, ErrAuthServiceNotAvailable
	}
	return a.client.RegisterUser(ctx, req)
}

func (a *Adapter) LoginUser(ctx context.Context, req *pb.LoginUserRequest) (*pb.LoginUserResponse, error) {
	if a.client == nil {
		return nil, ErrAuthServiceNotAvailable
	}
	return a.client.LoginUser(ctx, req)
}

func (a *Adapter) UpdateUser(ctx context.Context, req *pb.UpdateUserRequest) (*pb.UpdateUserResponse, error) {
	if a.client == nil {
		return nil, ErrAuthServiceNotAvailable
	}
	return a.client.UpdateUser(ctx, req)
}

func (a *Adapter) ValidateToken(ctx context.Context, token string) (string, error) {
	if a.client == nil {
		return "", ErrAuthServiceNotAvailable
	}
	// Use LoginUser to validate token - the token can be used as the password
	resp, err := a.client.LoginUser(ctx, &pb.LoginUserRequest{
		Username: token,
		Password: token,
	})
	if err != nil {
		return "", err
	}
	if resp.User == nil {
		return "", errors.New("user not found")
	}
	return resp.User.GetId(), nil
}

func (a *Adapter) Close() error {
	if a.conn != nil {
		return a.conn.Close()
	}
	return nil
}
