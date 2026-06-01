package cache

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"services/product/internal/applications/core/domain"
	"time"

	"github.com/redis/go-redis/v9"
)

const (
	productKeyPrefix  = "product:"
	productListPrefix = "product_list:"
	defaultTTL        = 10 * time.Minute // Hot product cache: 10 min
)

type Adapter struct {
	client *redis.Client
	ctx    context.Context
}

func NewAdapter(redisURL string) (*Adapter, error) {
	opts, err := redis.ParseURL(redisURL)
	if err != nil {
		return nil, fmt.Errorf("redis url parse error: %v", err)
	}
	client := redis.NewClient(opts)

	ctx := context.Background()
	if err := client.Ping(ctx).Err(); err != nil {
		return nil, fmt.Errorf("redis connection error: %v", err)
	}
	log.Println("Connected to Redis")
	return &Adapter{client: client, ctx: ctx}, nil
}

// ── Single product cache ──

func (a *Adapter) GetProduct(id int64) (*domain.Product, error) {
	key := fmt.Sprintf("%s%d", productKeyPrefix, id)
	data, err := a.client.Get(a.ctx, key).Bytes()
	if err != nil {
		return nil, err
	}
	var product domain.Product
	if err := json.Unmarshal(data, &product); err != nil {
		return nil, err
	}
	return &product, nil
}

func (a *Adapter) SetProduct(product *domain.Product) error {
	key := fmt.Sprintf("%s%d", productKeyPrefix, product.ID)
	data, err := json.Marshal(product)
	if err != nil {
		return err
	}
	return a.client.Set(a.ctx, key, data, defaultTTL).Err()
}

func (a *Adapter) DeleteProduct(id int64) error {
	key := fmt.Sprintf("%s%d", productKeyPrefix, id)
	return a.client.Del(a.ctx, key).Err()
}

// ── Product list cache ──

type cachedList struct {
	Products   []*domain.Product `json:"products"`
	TotalCount int32             `json:"total_count"`
}

func (a *Adapter) GetProductList(key string) ([]*domain.Product, int32, error) {
	fullKey := productListPrefix + key
	data, err := a.client.Get(a.ctx, fullKey).Bytes()
	if err != nil {
		return nil, 0, err
	}
	var cached cachedList
	if err := json.Unmarshal(data, &cached); err != nil {
		return nil, 0, err
	}
	return cached.Products, cached.TotalCount, nil
}

func (a *Adapter) SetProductList(key string, products []*domain.Product, totalCount int32) error {
	fullKey := productListPrefix + key
	data, err := json.Marshal(cachedList{Products: products, TotalCount: totalCount})
	if err != nil {
		return err
	}
	return a.client.Set(a.ctx, fullKey, data, defaultTTL).Err()
}

func (a *Adapter) DeleteProductList() error {
	// Delete all list cache keys using pattern scan
	iter := a.client.Scan(a.ctx, 0, productListPrefix+"*", 100).Iterator()
	for iter.Next(a.ctx) {
		a.client.Del(a.ctx, iter.Val())
	}
	return iter.Err()
}
