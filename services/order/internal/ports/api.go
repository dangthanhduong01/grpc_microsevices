package ports

import "services/order/internal/applications/core/domain"

type APIPort interface {
	PlaceOrder(order domain.Order) (domain.Order, error)
}
