package grpc

import (
	"context"

	pb "github.com/dangthanhduong01/microservices_proto/pb/product"
)

func (a Adapter) GetProduct(ctx context.Context, req *pb.GetProductRequest) (*pb.GetProductResponse, error) {
	product, err := a.api.GetProduct(req.Id)
	if err != nil {
		return nil, err
	}
	return &pb.GetProductResponse{
		Product: &pb.Product{
			Id:          product.ID,
			Name:        product.Name,
			Description: product.Description,
			Price:       product.Price,
			Stock:       product.Stock,
			Category:    product.Category,
			CreatedAt:   product.CreatedAt,
			UpdatedAt:   product.UpdatedAt,
		},
	}, nil
}

func (a Adapter) ListProducts(ctx context.Context, req *pb.ListProductsRequest) (*pb.ListProductsResponse, error) {
	products, totalCount, err := a.api.ListProducts(req.Page, req.PageSize, req.Category)
	if err != nil {
		return nil, err
	}

	var pbProducts []*pb.Product
	for _, p := range products {
		pbProducts = append(pbProducts, &pb.Product{
			Id:          p.ID,
			Name:        p.Name,
			Description: p.Description,
			Price:       p.Price,
			Stock:       p.Stock,
			Category:    p.Category,
			CreatedAt:   p.CreatedAt,
			UpdatedAt:   p.UpdatedAt,
		})
	}

	return &pb.ListProductsResponse{
		Products:   pbProducts,
		TotalCount: totalCount,
	}, nil
}

func (a Adapter) CreateProduct(ctx context.Context, req *pb.CreateProductRequest) (*pb.CreateProductResponse, error) {
	product, err := a.api.CreateProduct(req.Name, req.Description, req.Price, req.Stock, req.Category)
	if err != nil {
		return nil, err
	}
	return &pb.CreateProductResponse{
		Product: &pb.Product{
			Id:          product.ID,
			Name:        product.Name,
			Description: product.Description,
			Price:       product.Price,
			Stock:       product.Stock,
			Category:    product.Category,
			CreatedAt:   product.CreatedAt,
			UpdatedAt:   product.UpdatedAt,
		},
	}, nil
}

func (a Adapter) UpdateProduct(ctx context.Context, req *pb.UpdateProductRequest) (*pb.UpdateProductResponse, error) {
	product, err := a.api.UpdateProduct(req.Id, req.Name, req.Description, req.Price, req.Category)
	if err != nil {
		return nil, err
	}
	return &pb.UpdateProductResponse{
		Product: &pb.Product{
			Id:          product.ID,
			Name:        product.Name,
			Description: product.Description,
			Price:       product.Price,
			Stock:       product.Stock,
			Category:    product.Category,
			CreatedAt:   product.CreatedAt,
			UpdatedAt:   product.UpdatedAt,
		},
	}, nil
}

func (a Adapter) UpdateStock(ctx context.Context, req *pb.UpdateStockRequest) (*pb.UpdateStockResponse, error) {
	product, err := a.api.UpdateStock(req.ProductId, req.Quantity)
	if err != nil {
		return nil, err
	}
	return &pb.UpdateStockResponse{
		ProductId: product.ID,
		NewStock:  product.Stock,
	}, nil
}

func (a Adapter) GetStock(ctx context.Context, req *pb.GetStockRequest) (*pb.GetStockResponse, error) {
	stock, err := a.api.GetStock(req.ProductId)
	if err != nil {
		return nil, err
	}
	return &pb.GetStockResponse{
		ProductId: req.ProductId,
		Stock:     stock,
	}, nil
}
