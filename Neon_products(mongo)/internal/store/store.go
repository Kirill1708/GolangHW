package store

import (
	"context"
	"neon_products/internal/models"
)

type Store interface {
	Create(ctx context.Context, neon *models.Neon) error
	All(ctx context.Context) ([]*models.Neon, error)
	ByID(ctx context.Context, id string) (*models.Neon, error)
	Update(ctx context.Context, neon *models.Neon) error
	Delete(ctx context.Context, id string) error
	// Laptops() LaptopsRepository
	// Phones() PhonesRepository
}

// electronics
//   laptops
//   phones

// TODO дома почитать, вернемся в будущих лекциях
// type LaptopsRepository interface {
// 	Create(ctx context.Context, laptop *models.Laptop) error
// 	All(ctx context.Context) ([]*models.Laptop, error)
// 	ByID(ctx context.Context, id int) (*models.Laptop, error)
// 	Update(ctx context.Context, laptop *models.Laptop) error
// 	Delete(ctx context.Context, id int) error
// }
