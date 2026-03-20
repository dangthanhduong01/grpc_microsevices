package ports

import (
	"context"
	"services/payment/internal/applications/core/domain"
)

type APIPort interface {
	Charge(ctx context.Context, payment domain.Payment) (domain.Payment, error)
}
