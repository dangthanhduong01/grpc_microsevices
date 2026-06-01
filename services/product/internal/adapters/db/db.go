package db

import (
	"fmt"
	"services/product/internal/applications/core/domain"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type ProductModel struct {
	gorm.Model
	Name        string  `gorm:"not null"`
	Description string  `gorm:"type:text"`
	Price       float32 `gorm:"not null"`
	Stock       int32   `gorm:"not null;default:0"`
	Category    string  `gorm:"index;not null"`
}

func (ProductModel) TableName() string {
	return "products"
}

type Adapter struct {
	db *gorm.DB
}

func NewAdapter(dataSourceURL string) (*Adapter, error) {
	db, err := gorm.Open(postgres.Open(dataSourceURL), &gorm.Config{})
	if err != nil {
		return nil, fmt.Errorf("db connection error: %v", err)
	}
	if err := db.AutoMigrate(&ProductModel{}); err != nil {
		return nil, fmt.Errorf("db migration error: %v", err)
	}
	return &Adapter{db: db}, nil
}

func (a *Adapter) GetByID(id int64) (*domain.Product, error) {
	var model ProductModel
	if err := a.db.First(&model, id).Error; err != nil {
		return nil, err
	}
	return toDomain(&model), nil
}

func (a *Adapter) List(offset, limit int, category string) ([]*domain.Product, int32, error) {
	var models []ProductModel
	var totalCount int64

	query := a.db.Model(&ProductModel{})
	if category != "" {
		query = query.Where("category = ?", category)
	}

	if err := query.Count(&totalCount).Error; err != nil {
		return nil, 0, err
	}

	if err := query.Offset(offset).Limit(limit).Order("created_at DESC").Find(&models).Error; err != nil {
		return nil, 0, err
	}

	products := make([]*domain.Product, len(models))
	for i, m := range models {
		products[i] = toDomain(&m)
	}
	return products, int32(totalCount), nil
}

func (a *Adapter) Save(product *domain.Product) error {
	model := toModel(product)
	if err := a.db.Create(model).Error; err != nil {
		return err
	}
	product.ID = int64(model.ID)
	product.CreatedAt = model.CreatedAt.Unix()
	product.UpdatedAt = model.UpdatedAt.Unix()
	return nil
}

func (a *Adapter) Update(product *domain.Product) error {
	return a.db.Model(&ProductModel{}).Where("id = ?", product.ID).Updates(map[string]interface{}{
		"name":        product.Name,
		"description": product.Description,
		"price":       product.Price,
		"stock":       product.Stock,
		"category":    product.Category,
	}).Error
}

// ── Converters ──

func toDomain(m *ProductModel) *domain.Product {
	return &domain.Product{
		ID:          int64(m.ID),
		Name:        m.Name,
		Description: m.Description,
		Price:       m.Price,
		Stock:       m.Stock,
		Category:    m.Category,
		CreatedAt:   m.CreatedAt.Unix(),
		UpdatedAt:   m.UpdatedAt.Unix(),
	}
}

func toModel(p *domain.Product) *ProductModel {
	return &ProductModel{
		Name:        p.Name,
		Description: p.Description,
		Price:       p.Price,
		Stock:       p.Stock,
		Category:    p.Category,
	}
}
