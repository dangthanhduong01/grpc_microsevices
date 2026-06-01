package ports

import "services/product/internal/applications/core/domain"

type DBPort interface {
	GetByID(id int64) (*domain.Product, error)
	List(offset, limit int, category string) ([]*domain.Product, int32, error)
	Save(product *domain.Product) error
	Update(product *domain.Product) error
}
