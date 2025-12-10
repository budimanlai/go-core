package base

import (
	"context"
	"fmt"

	"gorm.io/gorm"
)

// Example entity
type Product struct {
	ID       uint   `gorm:"primaryKey"`
	Name     string `gorm:"not null"`
	Category string `gorm:"index"`
	Price    float64
	Stock    int
	Status   string `gorm:"default:'active'"`
}

// ExampleFindOne demonstrates how to use FindOne method
func ExampleFindOne() {
	// Assume db is initialized
	var db *gorm.DB
	repo := NewGormRepository[Product](db)

	// Example 1: Find by unique field (e.g., email, code)
	product, err := repo.FindOne(context.Background(),
		func(db *gorm.DB) *gorm.DB {
			return db.Where("name = ?", "Laptop Dell XPS")
		},
	)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}
	if product == nil {
		fmt.Println("Product not found")
		return
	}
	fmt.Printf("Found: %+v\n", product)

	// Example 2: Find with multiple conditions
	activeProduct, err := repo.FindOne(context.Background(),
		func(db *gorm.DB) *gorm.DB {
			return db.Where("category = ? AND status = ?", "Electronics", "active")
		},
		func(db *gorm.DB) *gorm.DB {
			return db.Order("price DESC") // Get most expensive
		},
	)
	fmt.Printf("Most expensive active electronics: %+v\n", activeProduct)

	// Example 3: Find with Preload
	productWithRelations, err := repo.FindOne(context.Background(),
		func(db *gorm.DB) *gorm.DB {
			return db.Where("id = ?", 1)
		},
		func(db *gorm.DB) *gorm.DB {
			return db.Preload("Reviews").Preload("Category")
		},
	)
	fmt.Printf("Product with relations: %+v\n", productWithRelations)
}

// ExampleCreateBatch demonstrates bulk insert
func ExampleCreateBatch() {
	var db *gorm.DB
	repo := NewGormRepository[Product](db)

	// Create multiple products at once
	products := []*Product{
		{Name: "Product 1", Category: "Electronics", Price: 1000, Stock: 10},
		{Name: "Product 2", Category: "Electronics", Price: 2000, Stock: 5},
		{Name: "Product 3", Category: "Books", Price: 50, Stock: 100},
		{Name: "Product 4", Category: "Books", Price: 75, Stock: 50},
		{Name: "Product 5", Category: "Clothing", Price: 150, Stock: 20},
	}

	// Bulk insert - efficient single query
	// SQL: INSERT INTO products (name, category, price, stock) VALUES (...), (...), (...)
	err := repo.CreateBatch(context.Background(), products)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}

	fmt.Println("✅ Inserted", len(products), "products")
	// After CreateBatch, each product will have its ID populated
	for _, p := range products {
		fmt.Printf("Product ID: %d, Name: %s\n", p.ID, p.Name)
	}
}

// ExampleCreateBatchLarge demonstrates handling large batches
func ExampleCreateBatchLarge() {
	var db *gorm.DB
	repo := NewGormRepository[Product](db)

	// Create 1000 products
	products := make([]*Product, 1000)
	for i := 0; i < 1000; i++ {
		products[i] = &Product{
			Name:     fmt.Sprintf("Product %d", i+1),
			Category: "Bulk",
			Price:    float64(i * 10),
			Stock:    i,
		}
	}

	// CreateBatch automatically chunks into batches of 100
	// This will execute 10 INSERT queries (1000 / 100)
	err := repo.CreateBatch(context.Background(), products)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}

	fmt.Println("✅ Inserted 1000 products in batches of 100")
}

// ExampleDeleteBatch demonstrates bulk delete
func ExampleDeleteBatch() {
	var db *gorm.DB
	repo := NewGormRepository[Product](db)

	// Delete multiple products by IDs
	idsToDelete := []any{1, 2, 3, 4, 5}

	// Single query: DELETE FROM products WHERE id IN (1, 2, 3, 4, 5)
	err := repo.DeleteBatch(context.Background(), idsToDelete)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}

	fmt.Println("✅ Deleted", len(idsToDelete), "products")
}

// ExampleDeleteBatchWithSoftDelete demonstrates soft delete behavior
func ExampleDeleteBatchWithSoftDelete() {
	var db *gorm.DB
	repo := NewGormRepository[Product](db)

	// If Product has DeletedAt field, this will be soft delete
	// SQL: UPDATE products SET deleted_at = NOW() WHERE id IN (...)
	idsToDelete := []any{10, 11, 12}

	err := repo.DeleteBatch(context.Background(), idsToDelete)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}

	fmt.Println("✅ Soft deleted", len(idsToDelete), "products")

	// Verify - FindByID won't return soft-deleted records
	product, _ := repo.FindByID(context.Background(), 10)
	if product == nil {
		fmt.Println("Product 10 is now soft-deleted (not visible)")
	}
}

// ExampleCount demonstrates counting with filters
func ExampleCount() {
	var db *gorm.DB
	repo := NewGormRepository[Product](db)

	// Example 1: Count all products
	total, err := repo.Count(context.Background())
	if err != nil {
		fmt.Println("Error:", err)
		return
	}
	fmt.Printf("Total products: %d\n", total)

	// Example 2: Count with filter
	activeCount, err := repo.Count(context.Background(),
		func(db *gorm.DB) *gorm.DB {
			return db.Where("status = ?", "active")
		},
	)
	fmt.Printf("Active products: %d\n", activeCount)

	// Example 3: Count with multiple conditions
	expensiveElectronics, err := repo.Count(context.Background(),
		func(db *gorm.DB) *gorm.DB {
			return db.Where("category = ?", "Electronics")
		},
		func(db *gorm.DB) *gorm.DB {
			return db.Where("price > ?", 1000)
		},
	)
	fmt.Printf("Expensive electronics: %d\n", expensiveElectronics)

	// Example 4: Count low stock items
	lowStock, err := repo.Count(context.Background(),
		func(db *gorm.DB) *gorm.DB {
			return db.Where("stock < ? AND status = ?", 10, "active")
		},
	)
	fmt.Printf("Low stock items: %d\n", lowStock)
}

// ExampleCountVsFindAll demonstrates efficiency difference
func ExampleCountVsFindAll() {
	var db *gorm.DB
	repo := NewGormRepository[Product](db)

	// ❌ INEFFICIENT: Using FindAll just to get count
	result, _ := repo.FindAll(context.Background(), 1, 1,
		func(db *gorm.DB) *gorm.DB {
			return db.Where("status = ?", "active")
		},
	)
	total := result.Total // This fetches 1 record unnecessarily
	fmt.Printf("Total (inefficient): %d\n", total)

	// ✅ EFFICIENT: Using Count directly
	count, _ := repo.Count(context.Background(),
		func(db *gorm.DB) *gorm.DB {
			return db.Where("status = ?", "active")
		},
	)
	fmt.Printf("Total (efficient): %d\n", count)
	// SQL: SELECT COUNT(*) FROM products WHERE status = 'active'
	// No data fetching, much faster!
}

// ExampleCombinedUsage demonstrates real-world scenario
func ExampleCombinedUsage() {
	var db *gorm.DB
	repo := NewGormRepository[Product](db)
	ctx := context.Background()

	// Scenario: Bulk import products from CSV

	// 1. Check if category has products
	existingCount, _ := repo.Count(ctx,
		func(db *gorm.DB) *gorm.DB {
			return db.Where("category = ?", "Electronics")
		},
	)
	fmt.Printf("Existing electronics: %d\n", existingCount)

	// 2. Bulk insert new products
	newProducts := []*Product{
		{Name: "iPhone 15", Category: "Electronics", Price: 12000000, Stock: 50},
		{Name: "MacBook Pro", Category: "Electronics", Price: 25000000, Stock: 20},
		{Name: "iPad Air", Category: "Electronics", Price: 8000000, Stock: 30},
	}
	repo.CreateBatch(ctx, newProducts)

	// 3. Verify
	newCount, _ := repo.Count(ctx,
		func(db *gorm.DB) *gorm.DB {
			return db.Where("category = ?", "Electronics")
		},
	)
	fmt.Printf("After import: %d (added %d products)\n", newCount, newCount-existingCount)

	// 4. Find most expensive
	mostExpensive, _ := repo.FindOne(ctx,
		func(db *gorm.DB) *gorm.DB {
			return db.Where("category = ?", "Electronics")
		},
		func(db *gorm.DB) *gorm.DB {
			return db.Order("price DESC")
		},
	)
	fmt.Printf("Most expensive: %s - Rp %.0f\n", mostExpensive.Name, mostExpensive.Price)

	// 5. Clean up test data (batch delete)
	idsToDelete := []any{}
	for _, p := range newProducts {
		idsToDelete = append(idsToDelete, p.ID)
	}
	repo.DeleteBatch(ctx, idsToDelete)
	fmt.Println("✅ Cleanup complete")
}
