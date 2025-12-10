package base

import (
	"context"
	"errors"
	"math"

	"gorm.io/gorm"
)

// Implementasi Struct
type BaseRepositoryImpl[T any] struct {
	db *gorm.DB
}

// InjectTx: Memasukkan object Transaction ke dalam Context
func InjectTx(ctx context.Context, tx *gorm.DB) context.Context {
	return context.WithValue(ctx, txKey{}, tx)
}

// ExtractTx: Mengambil object Transaction dari Context (jika ada)
func ExtractTx(ctx context.Context) *gorm.DB {
	if tx, ok := ctx.Value(txKey{}).(*gorm.DB); ok {
		return tx
	}
	return nil
}

// Constructor Public
func NewGormRepository[T any](db *gorm.DB) BaseRepository[T] {
	return &BaseRepositoryImpl[T]{
		db: db,
	}
}

// Helper internal untuk memilih DB mana yang dipakai
func (r *BaseRepositoryImpl[T]) getDB(ctx context.Context) *gorm.DB {
	// 1. Cek apakah ada Transaksi "titipan" di context?
	tx := ExtractTx(ctx)
	if tx != nil {
		// Jika ada, gunakan TX tersebut, dan jangan lupa WithContext
		return tx.WithContext(ctx)
	}

	// 2. Jika tidak ada, gunakan DB default milik repo
	return r.db.WithContext(ctx)
}

func (r *BaseRepositoryImpl[T]) Create(ctx context.Context, entity *T) error {
	return r.getDB(ctx).Create(entity).Error
}

func (r *BaseRepositoryImpl[T]) FindByID(ctx context.Context, id any, scopes ...func(*gorm.DB) *gorm.DB) (*T, error) {
	var entity T

	// 1. Ambil DB dasar
	db := r.getDB(ctx)

	// 2. Apply Scopes (misal: Preload("Profile"))
	for _, scope := range scopes {
		db = scope(db)
	}

	// 3. Eksekusi
	err := db.First(&entity, id).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &entity, nil
}

func (r *BaseRepositoryImpl[T]) UpdateFields(ctx context.Context, id any, fields map[string]interface{}) error {
	return r.getDB(ctx).Model(new(T)).Where("id = ?", id).Updates(fields).Error
}

func (r *BaseRepositoryImpl[T]) Update(ctx context.Context, entity *T) error {
	return r.getDB(ctx).Updates(entity).Error
}

func (r *BaseRepositoryImpl[T]) Delete(ctx context.Context, id any) error {
	var entity T
	return r.getDB(ctx).Delete(&entity, id).Error
}

// 5. LIST with Pagination
func (r *BaseRepositoryImpl[T]) FindAll(ctx context.Context, page, limit int, scopes ...func(*gorm.DB) *gorm.DB) (PaginationResult[T], error) {
	var entities []T
	var total int64

	// Default value safety
	if page <= 0 {
		page = 1
	}
	if limit <= 0 {
		limit = 10
	}
	if limit > 100 {
		limit = 100
	}

	// Mulai build query
	db := r.getDB(ctx).Model(new(T))

	// Apply Filter Dinamis (Scopes) SEBELUM count
	for _, scope := range scopes {
		db = scope(db)
	}

	if err := db.Session(&gorm.Session{}).
		Limit(-1).
		Offset(-1).
		Count(&total).Error; err != nil {
		return PaginationResult[T]{}, err
	}

	// Pagination Offset
	offset := (page - 1) * limit

	// Ambil Data
	if err := db.Offset(offset).Limit(limit).Find(&entities).Error; err != nil {
		return PaginationResult[T]{}, err
	}

	// Hitung Total Page
	totalPage := int(math.Ceil(float64(total) / float64(limit)))

	return PaginationResult[T]{
		Data:      entities,
		Total:     total,
		TotalPage: totalPage,
		Page:      page,
		Limit:     limit,
	}, nil
}

func (r *BaseRepositoryImpl[T]) Restore(ctx context.Context, id any) error {
	var entity T
	return r.getDB(ctx).Unscoped().Model(&entity).
		Where("id = ?", id).Update("deleted_at", nil).Error
}

func (r *BaseRepositoryImpl[T]) ForceDelete(ctx context.Context, id any) error {
	var entity T
	return r.getDB(ctx).Unscoped().Delete(&entity, id).Error
}

// FindOne: Find single entity by any condition using scopes
func (r *BaseRepositoryImpl[T]) FindOne(ctx context.Context, scopes ...func(*gorm.DB) *gorm.DB) (*T, error) {
	var entity T

	// Build query
	db := r.getDB(ctx)

	// Apply scopes (e.g., Where conditions, Preload, etc.)
	for _, scope := range scopes {
		db = scope(db)
	}

	// Execute
	err := db.First(&entity).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &entity, nil
}

// CreateBatch: Bulk insert entities (efficient batch processing)
func (r *BaseRepositoryImpl[T]) CreateBatch(ctx context.Context, entities []*T) error {
	if len(entities) == 0 {
		return nil
	}

	// GORM's CreateInBatches automatically handles chunking
	// Default batch size: 100 records per INSERT
	return r.getDB(ctx).CreateInBatches(entities, 100).Error
}

// DeleteBatch: Delete multiple entities by IDs
func (r *BaseRepositoryImpl[T]) DeleteBatch(ctx context.Context, ids []any) error {
	if len(ids) == 0 {
		return nil
	}

	var entity T
	// DELETE FROM table WHERE id IN (?, ?, ?)
	return r.getDB(ctx).Delete(&entity, ids).Error
}

// Count: Count entities with optional filters
func (r *BaseRepositoryImpl[T]) Count(ctx context.Context, scopes ...func(*gorm.DB) *gorm.DB) (int64, error) {
	var count int64

	// Build query
	db := r.getDB(ctx).Model(new(T))

	// Apply scopes (filters)
	for _, scope := range scopes {
		db = scope(db)
	}

	// Execute count
	err := db.Session(&gorm.Session{}).Limit(-1).Offset(-1).Count(&count).Error
	return count, err
}
