package api

import (
	"fmt"
	"log"
	"services/product/internal/applications/core/domain"
	"services/product/internal/ports"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type Application struct {
	db    ports.DBPort
	cache ports.CachePort
}

func NewApplication(db ports.DBPort, cache ports.CachePort) *Application {
	return &Application{db: db, cache: cache}
}

func (a *Application) GetProduct(id int64) (*domain.Product, error) {
	product, err := a.cache.GetProduct(id)
	if err == nil && product != nil {
		return product, nil
	}

	product, err = a.db.GetByID(id)
	if err != nil {
		return nil, status.Errorf(codes.NotFound, "product not found: %v", err)
	}

	if cacheErr := a.cache.SetProduct(product); cacheErr != nil {
		log.Printf("WARN: failed to cache product %d: %v", id, cacheErr)
	}

	return product, nil
}

func (a *Application) ListProducts(page, pageSize int32, category string) ([]*domain.Product, int32, error) {
	if page <= 0 {
		page = 1
	}
	if pageSize <= 0 {
		pageSize = 20
	}

	cacheKey := buildListCacheKey(page, pageSize, category)

	products, totalCount, err := a.cache.GetProductList(cacheKey)
	if err == nil && products != nil {
		return products, totalCount, nil
	}

	offset := int((page - 1) * pageSize)
	products, totalCount, err = a.db.List(offset, int(pageSize), category)
	if err != nil {
		return nil, 0, status.Errorf(codes.Internal, "failed to list products: %v", err)
	}

	if cacheErr := a.cache.SetProductList(cacheKey, products, totalCount); cacheErr != nil {
		log.Printf("WARN: failed to cache product list: %v", cacheErr)
	}

	return products, totalCount, nil
}

func (a *Application) CreateProduct(name, description string, price float32, stock int32, category string) (*domain.Product, error) {
	product := domain.NewProduct(name, description, price, stock, category)
	if err := a.db.Save(product); err != nil {
		return nil, status.Errorf(codes.Internal, "failed to create product: %v", err)
	}

	_ = a.cache.DeleteProductList()

	return product, nil
}

func (a *Application) UpdateProduct(id int64, name, description string, price float32, category string) (*domain.Product, error) {
	product, err := a.db.GetByID(id)
	if err != nil {
		return nil, status.Errorf(codes.NotFound, "product not found: %v", err)
	}

	product.Name = name
	product.Description = description
	product.Price = price
	product.Category = category

	if err := a.db.Update(product); err != nil {
		return nil, status.Errorf(codes.Internal, "failed to update product: %v", err)
	}

	_ = a.cache.DeleteProduct(id)
	_ = a.cache.DeleteProductList()

	return product, nil
}

func (a *Application) UpdateStock(productID int64, quantity int32) (*domain.Product, error) {
	product, err := a.db.GetByID(productID)
	if err != nil {
		return nil, status.Errorf(codes.NotFound, "product not found: %v", err)
	}

	if err := product.UpdateStock(quantity); err != nil {
		return nil, status.Errorf(codes.FailedPrecondition, "insufficient stock")
	}

	if err := a.db.Update(product); err != nil {
		return nil, status.Errorf(codes.Internal, "failed to update stock: %v", err)
	}

	// Invalidate caches
	_ = a.cache.DeleteProduct(productID)
	_ = a.cache.DeleteProductList()

	return product, nil
}

func (a *Application) GetStock(productID int64) (int32, error) {
	product, err := a.GetProduct(productID)
	if err != nil {
		return 0, err
	}
	return product.Stock, nil
}

func buildListCacheKey(page, pageSize int32, category string) string {
	if category == "" {
		category = "all"
	}
	return fmt.Sprintf("products:%s:%d:%d", category, page, pageSize)
}
