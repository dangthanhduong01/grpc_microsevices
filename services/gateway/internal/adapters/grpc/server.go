package grpc

import (
	"context"
	"fmt"
	"log"
	"net"
	"strings"

	"services/gateway/internal/adapters/auth"
	"services/gateway/internal/adapters/order"
	"services/gateway/internal/adapters/ratelimiter"

	pb "github.com/dangthanhduong01/microservices_proto/pb/auth"
	pbOrder "github.com/dangthanhduong01/microservices_proto/pb/order"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/reflection"
	"google.golang.org/grpc/status"
)

type Adapter struct {
	authClient   *auth.Adapter
	orderClient  *order.Adapter
	rateLimiter  *ratelimiter.RateLimiter
	port         int
	pb.UnimplementedAuthServiceServer
	pbOrder.UnimplementedOrderServiceServer
}

func NewAdapter(
	authServiceURL string,
	orderServiceURL string,
	port int,
	rateLimitRequests int,
	rateLimitWindowSeconds int,
) *Adapter {
	authClient, err := auth.NewAdapter(authServiceURL)
	if err != nil {
		log.Fatalf("failed to initialize auth client: %v", err)
	}

	orderClient, err := order.NewAdapter(orderServiceURL)
	if err != nil {
		log.Fatalf("failed to initialize order client: %v", err)
	}

	rateLimiter := ratelimiter.NewRateLimiter(rateLimitRequests, rateLimitWindowSeconds)

	return &Adapter{
		authClient:  authClient,
		orderClient: orderClient,
		rateLimiter: rateLimiter,
		port:        port,
	}
}

func (a *Adapter) Run() {
	listen, err := net.Listen("tcp", fmt.Sprintf(":%d", a.port))
	if err != nil {
		log.Fatalf("failed to listen on port %d: %v", a.port, err)
	}

	grpcServer := grpc.NewServer(
		grpc.UnaryInterceptor(a.unaryInterceptor),
	)
	pb.RegisterAuthServiceServer(grpcServer, a)
	pbOrder.RegisterOrderServiceServer(grpcServer, a)
	reflection.Register(grpcServer)

	log.Printf("Gateway server starting on port %d", a.port)
	if err := grpcServer.Serve(listen); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}

func (a *Adapter) unaryInterceptor(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
	// Check rate limit
	allowed, err := a.rateLimiter.Allow(ctx, "global")
	if err != nil {
		return nil, status.Errorf(codes.ResourceExhausted, "rate limit exceeded")
	}
	if !allowed {
		return nil, status.Errorf(codes.ResourceExhausted, "rate limit exceeded")
	}

	// For now, just pass through to the handler
	// Auth validation will be added when auth service is ready
	return handler(ctx, req)
}

func (a *Adapter) CreateOrder(ctx context.Context, req *pbOrder.CreateOrderRequest) (*pbOrder.CreateOrderResponse, error) {
	// Extract token from metadata
	token := a.extractToken(ctx)

	// Validate token
	if token != "" {
		userID, err := a.authClient.ValidateToken(ctx, token)
		if err != nil {
			return nil, status.Errorf(codes.Unauthenticated, "invalid token: %v", err)
		}
		log.Printf("Authenticated user: %s", userID)
	}

	// Forward to order service
	resp, err := a.orderClient.CreateOrder(ctx, req)
	if err != nil {
		return nil, err
	}

	return resp, nil
}

func (a *Adapter) extractToken(ctx context.Context) string {
	// Extract token from gRPC metadata
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return ""
	}

	// Check for authorization header
	authHeader := md.Get("authorization")
	if len(authHeader) > 0 {
		// Remove "Bearer " prefix if present
		token := authHeader[0]
		if strings.HasPrefix(token, "Bearer ") {
			return strings.TrimPrefix(token, "Bearer ")
		}
		return token
	}

	return ""
}

// RegisterUser handles user registration - forwards to auth service
func (a *Adapter) RegisterUser(ctx context.Context, req *pb.CreateUserRequest) (*pb.CreateUserResponse, error) {
	return a.authClient.RegisterUser(ctx, req)
}

// LoginUser handles user login - forwards to auth service
func (a *Adapter) LoginUser(ctx context.Context, req *pb.LoginUserRequest) (*pb.LoginUserResponse, error) {
	return a.authClient.LoginUser(ctx, req)
}

// UpdateUser handles user profile update - requires authentication
func (a *Adapter) UpdateUser(ctx context.Context, req *pb.UpdateUserRequest) (*pb.UpdateUserResponse, error) {
	// Extract and validate token
	token := a.extractToken(ctx)
	if token == "" {
		return nil, status.Errorf(codes.Unauthenticated, "authorization token required")
	}

	_, err := a.authClient.ValidateToken(ctx, token)
	if err != nil {
		return nil, status.Errorf(codes.Unauthenticated, "invalid token: %v", err)
	}

	// Forward to auth service
	return a.authClient.UpdateUser(ctx, req)
}
