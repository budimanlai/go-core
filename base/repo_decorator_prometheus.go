package base

import (
	"context"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"gorm.io/gorm"
)

// -----------------------------------------------------------
// 1. Definisi Metric Global
// -----------------------------------------------------------
var (
	// Histogram sangat cocok untuk mengukur durasi (latency)
	// dan otomatis menyediakan count (jumlah request)
	dbDurationHistogram = prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Namespace: "gocore", // Prefix metric
			Subsystem: "db",
			Name:      "query_duration_seconds",
			Help:      "Duration of database queries in seconds.",
			// Buckets untuk mengelompokkan durasi (dari 5ms sampai 10 detik)
			Buckets: []float64{0.005, 0.01, 0.025, 0.05, 0.1, 0.25, 0.5, 1, 2.5, 5, 10},
		},
		// Label yang akan kita isi dinamis
		[]string{"entity", "operation", "status"},
	)
)

// Init function berjalan otomatis saat aplikasi start
func init() {
	// Mendaftarkan metric ke Default Registerer Prometheus
	prometheus.MustRegister(dbDurationHistogram)
}

// -----------------------------------------------------------
// 2. Struct Decorator
// -----------------------------------------------------------

type prometheusRepository[T any] struct {
	next BaseRepository[T]
	name string // Nama Entity, misal "User" atau "Charger"
}

// -----------------------------------------------------------
// 3. Implementasi Interface BaseRepository
// -----------------------------------------------------------

func (r *prometheusRepository[T]) Create(ctx context.Context, entity *T) error {
	start := time.Now()
	err := r.next.Create(ctx, entity)
	r.record("Create", time.Since(start), err)
	return err
}

// Perhatikan: Signature FindByID sudah sesuai update terakhir (menerima scopes)
func (r *prometheusRepository[T]) FindByID(ctx context.Context, id any, scopes ...func(*gorm.DB) *gorm.DB) (*T, error) {
	start := time.Now()
	res, err := r.next.FindByID(ctx, id, scopes...)
	r.record("FindByID", time.Since(start), err)
	return res, err
}

func (r *prometheusRepository[T]) UpdateFields(ctx context.Context, id any, fields map[string]interface{}) error {
	start := time.Now()
	err := r.next.UpdateFields(ctx, id, fields)
	r.record("UpdateFields", time.Since(start), err)
	return err
}

func (r *prometheusRepository[T]) Update(ctx context.Context, entity *T) error {
	start := time.Now()
	err := r.next.Update(ctx, entity)
	r.record("Update", time.Since(start), err)
	return err
}

func (r *prometheusRepository[T]) Delete(ctx context.Context, id any) error {
	start := time.Now()
	err := r.next.Delete(ctx, id)
	r.record("Delete", time.Since(start), err)
	return err
}

func (r *prometheusRepository[T]) FindAll(ctx context.Context, page, limit int, scopes ...func(*gorm.DB) *gorm.DB) (PaginationResult[T], error) {
	start := time.Now()
	res, err := r.next.FindAll(ctx, page, limit, scopes...)
	r.record("FindAll", time.Since(start), err)
	return res, err
}

func (r *prometheusRepository[T]) Restore(ctx context.Context, id any) error {
	start := time.Now()
	err := r.next.Restore(ctx, id)
	r.record("Restore", time.Since(start), err)
	return err
}

func (r *prometheusRepository[T]) ForceDelete(ctx context.Context, id any) error {
	start := time.Now()
	err := r.next.ForceDelete(ctx, id)
	r.record("ForceDelete", time.Since(start), err)
	return err
}

func (r *prometheusRepository[T]) FindOne(ctx context.Context, scopes ...func(*gorm.DB) *gorm.DB) (*T, error) {
	start := time.Now()
	res, err := r.next.FindOne(ctx, scopes...)
	r.record("FindOne", time.Since(start), err)
	return res, err
}

func (r *prometheusRepository[T]) CreateBatch(ctx context.Context, entities []*T) error {
	start := time.Now()
	err := r.next.CreateBatch(ctx, entities)
	r.record("CreateBatch", time.Since(start), err)
	return err
}

func (r *prometheusRepository[T]) DeleteBatch(ctx context.Context, ids []any) error {
	start := time.Now()
	err := r.next.DeleteBatch(ctx, ids)
	r.record("DeleteBatch", time.Since(start), err)
	return err
}

func (r *prometheusRepository[T]) Count(ctx context.Context, scopes ...func(*gorm.DB) *gorm.DB) (int64, error) {
	start := time.Now()
	count, err := r.next.Count(ctx, scopes...)
	r.record("Count", time.Since(start), err)
	return count, err
}

// -----------------------------------------------------------
// 4. Helper Recorder
// -----------------------------------------------------------

func (r *prometheusRepository[T]) record(operation string, duration time.Duration, err error) {
	status := "success"
	if err != nil {
		status = "error"
	}

	// Mengisi label dan mencatat durasi dalam detik
	dbDurationHistogram.WithLabelValues(
		r.name,    // entity: "User"
		operation, // operation: "FindByID"
		status,    // status: "success"/"error"
	).Observe(duration.Seconds())
}
