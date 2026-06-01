package domain

import "time"

type Product struct {
	ID          int64   `json:"id"`
	Name        string  `json:"name"`
	Description string  `json:"description"`
	Price       float32 `json:"price"`
	Stock       int32   `json:"stock"`
	Category    string  `json:"category"`
	CreatedAt   int64   `json:"created_at"`
	UpdatedAt   int64   `json:"updated_at"`
}

func NewProduct(name, description string, price float32, stock int32, category string) *Product {
	now := time.Now().Unix()
	return &Product{
		Name:        name,
		Description: description,
		Price:       price,
		Stock:       stock,
		Category:    category,
		CreatedAt:   now,
		UpdatedAt:   now,
	}
}

func (p *Product) UpdateStock(quantity int32) error {
	newStock := p.Stock + quantity
	if newStock < 0 {
		return ErrInsufficientStock
	}
	p.Stock = newStock
	p.UpdatedAt = time.Now().Unix()
	return nil
}
