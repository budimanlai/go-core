package base

import (
	"context"

	"gorm.io/gorm"
)

type txKey struct{}

// Pagination Result Wrapper (Opsional tapi recommended)
type PaginationResult[T any] struct {
	Data      []T   `json:"data"`
	Total     int64 `json:"total"`
	TotalPage int   `json:"total_page"`
	Page      int   `json:"page"`
	Limit     int   `json:"limit"`
}

type BaseRepository[T any] interface {
	Create(ctx context.Context, entity *T) error
	FindByID(ctx context.Context, id any, scopes ...func(*gorm.DB) *gorm.DB) (*T, error)
	Update(ctx context.Context, entity *T) error
	UpdateFields(ctx context.Context, id any, fields map[string]interface{}) error
	Delete(ctx context.Context, id any) error
	FindAll(ctx context.Context, page, limit int, scopes ...func(*gorm.DB) *gorm.DB) (PaginationResult[T], error)
	FindOne(ctx context.Context, scopes ...func(*gorm.DB) *gorm.DB) (*T, error)

	// Soft delete specific
	Restore(ctx context.Context, id any) error
	ForceDelete(ctx context.Context, id any) error

	// Batch operations
	CreateBatch(ctx context.Context, entities []*T) error
	DeleteBatch(ctx context.Context, ids []any) error

	Count(ctx context.Context, scopes ...func(*gorm.DB) *gorm.DB) (int64, error)
}
