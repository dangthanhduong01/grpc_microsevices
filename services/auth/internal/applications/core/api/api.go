package api

import (
	"context"
	"errors"
	"fmt"
	"time"

	"services/auth/internal/ports"

	pb "github.com/dangthanhduong01/microservices_proto/pb/auth"
	"github.com/golang-jwt/jwt/v5"
	"google.golang.org/protobuf/types/known/timestamppb"
)

var (
	ErrUserAlreadyExists    = errors.New("user already exists")
	ErrInvalidCredentials  = errors.New("invalid credentials")
	ErrUserNotFound         = errors.New("user not found")
)

type Application struct {
	db ports.DBPort
}

func NewApplication(db ports.DBPort) *Application {
	return &Application{db: db}
}

func (a *Application) RegisterUser(ctx context.Context, req *pb.CreateUserRequest) (*pb.CreateUserResponse, error) {
	// Check if user with email already exists
	existingUserID, err := a.db.GetUserByEmail(req.GetEmail())
	if err == nil && existingUserID != "" {
		return nil, ErrUserAlreadyExists
	}

	// Check if user with username already exists
	existingUserID, err = a.db.GetUserByUsername(req.GetUsername())
	if err == nil && existingUserID != "" {
		return nil, ErrUserAlreadyExists
	}

	// Create new user
	userID, err := a.db.CreateUser(req.GetEmail(), req.GetUsername(), req.GetPassword())
	if err != nil {
		return nil, err
	}

	return &pb.CreateUserResponse{
		User: &pb.User{
			Id:        userID,
			Username:  req.GetUsername(),
			Email:     req.GetEmail(),
			CreatedAt: timestamppb.Now(),
			UpdatedAt: timestamppb.Now(),
		},
	}, nil
}

func (a *Application) LoginUser(ctx context.Context, req *pb.LoginUserRequest) (*pb.LoginUserResponse, error) {
	// Try to find user by username
	userID, err := a.db.GetUserByUsername(req.GetUsername())
	if err != nil {
		return nil, ErrInvalidCredentials
	}

	// Verify password (you should add proper password verification with bcrypt)
	// For now, simplified - in production, verify against hashed password in DB

	// Generate access token
	accessToken, accessExpiry, err := a.generateAccessToken(userID, req.GetUsername())
	if err != nil {
		return nil, err
	}

	// Generate refresh token
	refreshToken, refreshExpiry, err := a.generateRefreshToken(userID)
	if err != nil {
		return nil, err
	}

	return &pb.LoginUserResponse{
		User: &pb.User{
			Id:       userID,
			Username: req.GetUsername(),
		},
		AccessToken:        accessToken,
		RefreshToken:       refreshToken,
		AccessTokenExpiry:  timestamppb.New(accessExpiry),
		RefreshTokenExpiry: timestamppb.New(refreshExpiry),
	}, nil
}

func (a *Application) generateAccessToken(userID, username string) (string, time.Time, error) {
	expiry := time.Duration(ports.GetAccessTokenExpiry()) * time.Minute
	expiryTime := time.Now().Add(expiry)

	claims := jwt.MapClaims{
		"sub":     userID,
		"username": username,
		"exp":     expiryTime.Unix(),
		"iat":     time.Now().Unix(),
		"type":    "access",
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString([]byte(ports.GetJWTSecret()))
	if err != nil {
		return "", time.Time{}, err
	}

	return tokenString, expiryTime, nil
}

func (a *Application) generateRefreshToken(userID string) (string, time.Time, error) {
	expiry := time.Duration(ports.GetRefreshTokenExpiry()) * 24 * time.Hour
	expiryTime := time.Now().Add(expiry)

	claims := jwt.MapClaims{
		"sub":  userID,
		"exp":  expiryTime.Unix(),
		"iat":  time.Now().Unix(),
		"type": "refresh",
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString([]byte(ports.GetJWTSecret()))
	if err != nil {
		return "", time.Time{}, err
	}

	return tokenString, expiryTime, nil
}

// ValidateToken validates a JWT token and returns the claims
func (a *Application) ValidateToken(tokenString string) (jwt.MapClaims, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(ports.GetJWTSecret()), nil
	})

	if err != nil {
		return nil, err
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		return claims, nil
	}

	return nil, fmt.Errorf("invalid token")
}

func (a *Application) UpdateUser(ctx context.Context, req *pb.UpdateUserRequest) (*pb.UpdateUserResponse, error) {
	userID := req.GetId()
	if userID == "" {
		return nil, fmt.Errorf("user id is required")
	}

	// Verify user exists
	existingUserID, err := a.db.GetUserByID(userID)
	if err != nil || existingUserID == "" {
		return nil, ErrUserNotFound
	}

	// Check if email is being changed and if it's already in use
	if req.GetEmail() != "" {
		existingID, err := a.db.GetUserByEmail(req.GetEmail())
		if err == nil && existingID != "" && existingID != userID {
			return nil, ErrUserAlreadyExists
		}
	}

	// Check if username is being changed and if it's already in use
	if req.GetUsername() != "" {
		existingID, err := a.db.GetUserByUsername(req.GetUsername())
		if err == nil && existingID != "" && existingID != userID {
			return nil, ErrUserAlreadyExists
		}
	}

	// Update user in DB
	err = a.db.UpdateUser(userID, req.GetEmail(), req.GetUsername(), req.GetFullName())
	if err != nil {
		return nil, err
	}

	return &pb.UpdateUserResponse{
		User: &pb.User{
			Id:        userID,
			Username:  req.GetUsername(),
			Email:     req.GetEmail(),
			UpdatedAt: timestamppb.Now(),
		},
	}, nil
}