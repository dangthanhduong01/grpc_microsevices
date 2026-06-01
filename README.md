# Demo microservice using grpc

# Techstack
- Go
- Gprc



## Book tutorial
grpc MicroService In go

Gateway HTTP API with Gin and Swagger                                                  
                                                                                         
  New file: gateway/internal/adapters/http/server.go                                     
                                                                                         
###  API Endpoints:                                                                         
  ┌────────┬───────────────────────┬──────┬────────────────────────┐                     
  │ Method │         Path          │ Auth │      Description       │                     
  ├────────┼───────────────────────┼──────┼────────────────────────┤                     
  │ POST   │ /api/v1/auth/register │ No   │ Register new user      │                     
  ├────────┼───────────────────────┼──────┼────────────────────────┤                     
  │ POST   │ /api/v1/auth/login    │ No   │ Login (get JWT tokens) │                     
  ├────────┼───────────────────────┼──────┼────────────────────────┤                     
  │ PUT    │ /api/v1/users/profile │ Yes  │ Update user profile    │                     
  ├────────┼───────────────────────┼──────┼────────────────────────┤                     
  │ POST   │ /api/v1/orders        │ Yes  │ Create new order       │                     
  ├────────┼───────────────────────┼──────┼────────────────────────┤                     
  │ GET    │ /api/v1/health        │ No   │ Health check           │                     
  └────────┴───────────────────────┴──────┴────────────────────────┘                     
###  Swagger Docs:                                                                          
  - URL: http://localhost:8080/swagger/index.html                                        
  - Generated files in gateway/docs/                                                     

### Run order service
ENV=development DATA_SOURCE_URL="postgres://user:pass@localhost:5432/db" \
  APPLICATION_PORT=50051 PAYMENT_SERVICE_URL="localhost:50052" \
  go run ./order/cmd

###  To run the gateway:                                                                    
  ENV=development APPLICATION_PORT=8080 \                                                
  AUTH_SERVICE_URL="localhost:50053" \                                                   
  ORDER_SERVICE_URL="localhost:50052" \                                                  
  JWT_SECRET="your-secret" \                                                             
  go run ./gateway/cmd                                                                   

### Run tests (none exist currently)
go test ./...                          

### Update dependencies
`go mod tidy` 

### Proto Files

Protocol buffers are defined in `/proto/` but generated code is imported from an external repository (`github.com/dangthanhduong01/microservices_proto`). To regenerate:

```bash
protoc --go_out=. --go_opt=paths=source_relative \
  --go-grpc_out=. --go-grpc_opt=paths=source_relative \
  proto/*.proto
```

###  Example requests:                                                                      
  #### Register                                                                             
  curl -X POST http://localhost:8080/api/v1/auth/register \                              
    -H "Content-Type: application/json" \                                                
    -d '{"email":"test@example.com","username":"testuser","password":"password123"}'     
                                                                                         
  #### Login                                                                                
  curl -X POST http://localhost:8080/api/v1/auth/login \                                 
    -H "Content-Type: application/json" \                                                
    -d '{"username":"testuser","password":"password123"}'                                
                                                                                         
  #### Create order (with JWT)                                                              
  curl -X POST http://localhost:8080/api/v1/orders \                                     
    -H "Content-Type: application/json" \                                                
    -H "Authorization: Bearer <token>" \                                                 
    -d '{"user_id":1,"items":[{"product_id":"prod123","quantity":2,"unit_price":10.00}],"
  total_price":20.00}'