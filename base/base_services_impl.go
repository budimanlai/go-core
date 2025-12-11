package base

import (
	"context"

	"gorm.io/gorm"
)

type baseServiceImpl[E any] struct {
	repo DomainRepository[E]
	db   *gorm.DB
}

// Constructor Service
// Parameter 'repo' bisa menerima 'BaseRepositoryImpl[E, M]' apapun M-nya.
func NewBaseService[E any](repo DomainRepository[E], db *gorm.DB) BaseService[E] {
	return &baseServiceImpl[E]{
		repo: repo,
		db:   db,
	}
}

// Helper untuk akses DB (buat transaction manual di custom service)
func (s *baseServiceImpl[E]) GetDB() *gorm.DB {
	return s.db
}

// --- Transaction Helper ---
func (s *baseServiceImpl[E]) WithTransaction(ctx context.Context, fn func(context.Context) error) error {
	return s.db.Transaction(func(tx *gorm.DB) error {
		txCtx := InjectTx(ctx, tx)
		return fn(txCtx)
	})
}

// --- Standard CRUD (Simple Delegation) ---
func (s *baseServiceImpl[E]) Create(ctx context.Context, entity *E) error {
	return s.repo.Create(ctx, entity)
}

func (s *baseServiceImpl[E]) Update(ctx context.Context, entity *E) error {
	return s.repo.Update(ctx, entity)
}

func (s *baseServiceImpl[E]) UpdateFields(ctx context.Context, id any, fields map[string]interface{}) error {
	return s.repo.UpdateFields(ctx, id, fields)
}

func (s *baseServiceImpl[E]) Delete(ctx context.Context, id any) error {
	return s.repo.Delete(ctx, id)
}

// --- Query Methods ---
func (s *baseServiceImpl[E]) FindByID(ctx context.Context, id any) (*E, error) {
	return s.repo.FindByID(ctx, id)
}

func (s *baseServiceImpl[E]) FindAll(ctx context.Context, page, limit int, scopes ...func(*gorm.DB) *gorm.DB) (PaginationResult[E], error) {
	return s.repo.FindAll(ctx, page, limit, scopes...)
}

func (s *baseServiceImpl[E]) FindOne(ctx context.Context, scopes ...func(*gorm.DB) *gorm.DB) (*E, error) {
	return s.repo.FindOne(ctx, scopes...)
}

func (s *baseServiceImpl[E]) Count(ctx context.Context, scopes ...func(*gorm.DB) *gorm.DB) (int64, error) {
	return s.repo.Count(ctx, scopes...)
}

// --- Batch Operations ---
func (s *baseServiceImpl[E]) CreateBatch(ctx context.Context, entities []*E) error {
	return s.repo.CreateBatch(ctx, entities)
}

func (s *baseServiceImpl[E]) DeleteBatch(ctx context.Context, ids []any) error {
	return s.repo.DeleteBatch(ctx, ids)
}

// --- Soft Delete Management ---
func (s *baseServiceImpl[E]) Restore(ctx context.Context, id any) error {
	return s.repo.Restore(ctx, id)
}

func (s *baseServiceImpl[E]) ForceDelete(ctx context.Context, id any) error {
	return s.repo.ForceDelete(ctx, id)
}
