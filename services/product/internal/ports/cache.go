package ports

import "services/product/internal/applications/core/domain"

// CachePort defines the interface for caching product data.
// Used for hot products / flash sale — 90% of reads should hit cache.
type CachePort interface {
	GetProduct(id int64) (*domain.Product, error)
	SetProduct(product *domain.Product) error
	DeleteProduct(id int64) error
	GetProductList(key string) ([]*domain.Product, int32, error)
	SetProductList(key string, products []*domain.Product, totalCount int32) error
	DeleteProductList() error
}
