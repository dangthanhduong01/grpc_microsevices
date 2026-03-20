package ports

import "services/order/internal/applications/core/domain"

type PaymentPort interface {
	Charge(*domain.Order) error
}
