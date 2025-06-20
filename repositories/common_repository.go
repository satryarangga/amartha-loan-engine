package repositories

import (
	"context"

	"github.com/satryarangga/amartha-loan-engine/models"
	"gorm.io/gorm"
)

type TransactionFunc func(tx *gorm.DB) error

// CommonRepository interface with generic type parameter
type CommonRepository[T any] interface {

	// Insert inserts a new record into the repository
	Insert(ctx context.Context, tx *gorm.DB, model *T) (string, error)

	// Update updates an existing record in the repository
	Update(ctx context.Context, tx *gorm.DB, model *T) error

	// FindByID finds a record by its ID
	FindByID(ctx context.Context, id string, relations []string) (*T, error)

	// FindAll finds all records matching the provided parameters
	FindAll(ctx context.Context, param models.FindAllParam) ([]T, error)

	// Wrapper transaction
	WithTransaction(ctx context.Context, fn TransactionFunc) error
}
