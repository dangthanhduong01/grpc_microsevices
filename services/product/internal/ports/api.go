package ports

import "services/product/internal/applications/core/domain"

type APIPort interface {
	GetProduct(id int64) (*domain.Product, error)
	ListProducts(page, pageSize int32, category string) ([]*domain.Product, int32, error)
	CreateProduct(name, description string, price float32, stock int32, category string) (*domain.Product, error)
	UpdateProduct(id int64, name, description string, price float32, category string) (*domain.Product, error)
	UpdateStock(productID int64, quantity int32) (*domain.Product, error)
	GetStock(productID int64) (int32, error)
}
