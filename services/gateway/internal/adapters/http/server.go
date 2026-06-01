package http

import (
	"context"
	"fmt"
	"net/http"

	"services/gateway/config"
	"services/gateway/internal/adapters/auth"
	"services/gateway/internal/adapters/order"

	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"

	pb "github.com/dangthanhduong01/microservices_proto/pb/auth"
	pbOrder "github.com/dangthanhduong01/microservices_proto/pb/order"
)

// @title Gateway API
// @version 1.0
// @description API documentation for the Gateway service
// @termsOfService http://swagger.io/terms/

// @contact.name API Support
// @contact.url http://www.example.com/support
// @contact.email support@example.com

// @license.name Apache 2.0
// @license.url http://www.apache.org/licenses/LICENSE-2.0.html

// @host localhost:8080
// @BasePath /api/v1

type Adapter struct {
	router       *gin.Engine
	authClient   *auth.Adapter
	orderClient  *order.Adapter
	server       *http.Server
}

func NewAdapter(
	authServiceURL string,
	orderServiceURL string,
) *Adapter {
	router := gin.Default()

	adapter := &Adapter{
		router: router,
	}

	// Initialize clients
	authClient, err := auth.NewAdapter(authServiceURL)
	if err != nil {
		panic("failed to initialize auth client: " + err.Error())
	}
	adapter.authClient = authClient

	orderClient, err := order.NewAdapter(orderServiceURL)
	if err != nil {
		panic("failed to initialize order client: " + err.Error())
	}
	adapter.orderClient = orderClient

	// Setup routes
	adapter.setupRoutes()

	adapter.server = &http.Server{
		Addr:    fmt.Sprintf(":%d", config.GetApplicationPort()),
		Handler: router,
	}

	return adapter
}

func (a *Adapter) setupRoutes() {
	// Swagger documentation
	a.router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	// API v1 routes
	v1 := a.router.Group("/api/v1")
	{
		// Health check
		v1.GET("/health", func(c *gin.Context) {
			c.JSON(http.StatusOK, gin.H{"status": "ok"})
		})

		// Auth routes (public)
		auth := v1.Group("/auth")
		{
			auth.POST("/register", a.RegisterUser)
			auth.POST("/login", a.LoginUser)
		}

		// Order routes (protected)
		orders := v1.Group("/orders")
		orders.Use(a.AuthMiddleware())
		{
			orders.POST("", a.CreateOrder)
		}

		// User profile routes (protected)
		users := v1.Group("/users")
		users.Use(a.AuthMiddleware())
		{
			users.PUT("/profile", a.UpdateUser)
		}
	}
}

// RegisterUser handles user registration
// @Summary Register a new user
// @Description Create a new user account
// @Tags auth
// @Accept json
// @Produce json
// @Param request body RegisterUserRequest true "User registration request"
// @Success 200 {object} pb.CreateUserResponse
// @Failure 400 {object} ErrorResponse
// @Failure 409 {object} ErrorResponse
// @Router /auth/register [post]
func (a *Adapter) RegisterUser(c *gin.Context) {
	var req RegisterUserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{Error: err.Error()})
		return
	}

	resp, err := a.authClient.RegisterUser(c.Request.Context(), &pb.CreateUserRequest{
		Email:    req.Email,
		Username: req.Username,
		Password: req.Password,
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{Error: err.Error()})
		return
	}

	c.JSON(http.StatusOK, resp)
}

// LoginUser handles user login
// @Summary User login
// @Description Authenticate user and return JWT tokens
// @Tags auth
// @Accept json
// @Produce json
// @Param request body LoginUserRequest true "Login request"
// @Success 200 {object} pb.LoginUserResponse
// @Failure 401 {object} ErrorResponse
// @Router /auth/login [post]
func (a *Adapter) LoginUser(c *gin.Context) {
	var req LoginUserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{Error: err.Error()})
		return
	}

	resp, err := a.authClient.LoginUser(c.Request.Context(), &pb.LoginUserRequest{
		Username: req.Username,
		Password: req.Password,
	})
	if err != nil {
		c.JSON(http.StatusUnauthorized, ErrorResponse{Error: err.Error()})
		return
	}

	c.JSON(http.StatusOK, resp)
}

// UpdateUser handles user profile update
// @Summary Update user profile
// @Description Update user profile information
// @Tags users
// @Accept json
// @Produce json
// @Param request body UpdateUserRequest true "Update profile request"
// @Success 200 {object} pb.UpdateUserResponse
// @Failure 401 {object} ErrorResponse
// @Security BearerAuth
// @Router /users/profile [put]
func (a *Adapter) UpdateUser(c *gin.Context) {
	var req UpdateUserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{Error: err.Error()})
		return
	}

	resp, err := a.authClient.UpdateUser(c.Request.Context(), &pb.UpdateUserRequest{
		Id:       req.Id,
		Email:    req.Email,
		Username: req.Username,
		FullName: req.FullName,
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{Error: err.Error()})
		return
	}

	c.JSON(http.StatusOK, resp)
}

// CreateOrder handles order creation
// @Summary Create a new order
// @Description Create a new order
// @Tags orders
// @Accept json
// @Produce json
// @Param request body CreateOrderRequest true "Order request"
// @Success 200 {object} pbOrder.CreateOrderResponse
// @Failure 400 {object} ErrorResponse
// @Security BearerAuth
// @Router /orders [post]
func (a *Adapter) CreateOrder(c *gin.Context) {
	var req CreateOrderRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{Error: err.Error()})
		return
	}

	items := make([]*pbOrder.Item, len(req.Items))
	for i, item := range req.Items {
		items[i] = &pbOrder.Item{
			ProductId: item.ProductId,
			Quantity:  item.Quantity,
			UnitPrice:  float32(item.UnitPrice),
		}
	}

	resp, err := a.orderClient.CreateOrder(c.Request.Context(), &pbOrder.CreateOrderRequest{
		UserId:     req.UserId,
		Items:      items,
		TotalPrice: float32(req.TotalPrice),
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{Error: err.Error()})
		return
	}

	c.JSON(http.StatusOK, resp)
}

// AuthMiddleware validates JWT token
func (a *Adapter) AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.JSON(http.StatusUnauthorized, ErrorResponse{Error: "authorization required"})
			c.Abort()
			return
		}

		token := authHeader
		if len(token) > 7 && token[:7] == "Bearer " {
			token = token[7:]
		}

		_, err := a.authClient.ValidateToken(c.Request.Context(), token)
		if err != nil {
			c.JSON(http.StatusUnauthorized, ErrorResponse{Error: "invalid token"})
			c.Abort()
			return
		}

		c.Next()
	}
}

func (a *Adapter) Run() error {
	return a.server.ListenAndServe()
}

func (a *Adapter) RunTLS(certFile, keyFile string) error {
	return a.server.ListenAndServeTLS(certFile, keyFile)
}

func (a *Adapter) Shutdown(ctx context.Context) error {
	return a.server.Shutdown(ctx)
}

type ErrorResponse struct {
	Error string `json:"error"`
}

// RegisterUserRequest represents the request body for registration
type RegisterUserRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

// LoginUserRequest represents the request body for login
type LoginUserRequest struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

// UpdateUserRequest represents the request body for profile update
type UpdateUserRequest struct {
	Id       string `json:"id" binding:"required"`
	Email    string `json:"email"`
	Username string `json:"username"`
	FullName string `json:"full_name"`
}

// CreateOrderRequest represents the request body for order creation
type CreateOrderRequest struct {
	UserId     int64   `json:"user_id" binding:"required"`
	Items      []Item  `json:"items" binding:"required,min=1"`
	TotalPrice float64 `json:"total_price" binding:"required"`
}

// Item represents an order item
type Item struct {
	ProductId string  `json:"product_id" binding:"required"`
	Quantity  int32   `json:"quantity" binding:"required,min=1"`
	UnitPrice float64 `json:"unit_price" binding:"required"`
}