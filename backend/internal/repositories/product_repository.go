package repositories

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"stockflow/backend/internal/models"
)

type ProductRepository struct {
	db *pgxpool.Pool
}

func NewProductRepository(db *pgxpool.Pool) *ProductRepository {
	return &ProductRepository{db: db}
}

func (r *ProductRepository) GetByID(ctx context.Context, id int) (*models.Product, error) {
	var product models.Product
	var lastMovement *time.Time

	query := `
		SELECT p.id, p.sku, p.name, p.price, p.stock, p.is_active, p.created_at,
		       MAX(sm.created_at) as last_movement_at
		FROM products p
		LEFT JOIN stock_movements sm ON p.id = sm.product_id
		WHERE p.id = $1
		GROUP BY p.id
	`

	err := r.db.QueryRow(ctx, query, id).Scan(
		&product.ID, &product.SKU, &product.Name, &product.Price, &product.Stock,
		&product.IsActive, &product.CreatedAt, &lastMovement,
	)

	if err == pgx.ErrNoRows {
		return nil, errors.New("product not found")
	}
	if err != nil {
		return nil, err
	}

	product.LastMovementAt = lastMovement
	return &product, nil
}

func (r *ProductRepository) GetBySKU(ctx context.Context, sku string) (*models.Product, error) {
	var product models.Product
	var lastMovement *time.Time

	query := `
		SELECT p.id, p.sku, p.name, p.price, p.stock, p.is_active, p.created_at,
		       MAX(sm.created_at) as last_movement_at
		FROM products p
		LEFT JOIN stock_movements sm ON p.id = sm.product_id
		WHERE p.sku = $1
		GROUP BY p.id
	`

	err := r.db.QueryRow(ctx, query, sku).Scan(
		&product.ID, &product.SKU, &product.Name, &product.Price, &product.Stock,
		&product.IsActive, &product.CreatedAt, &lastMovement,
	)

	if err == pgx.ErrNoRows {
		return nil, errors.New("product not found")
	}
	if err != nil {
		return nil, err
	}

	product.LastMovementAt = lastMovement
	return &product, nil
}

func (r *ProductRepository) List(ctx context.Context, filter *models.ProductFilter) (*models.ProductList, error) {
	query := `
		SELECT p.id, p.sku, p.name, p.price, p.stock, p.is_active, p.created_at,
		       MAX(sm.created_at) as last_movement_at
		FROM products p
		LEFT JOIN stock_movements sm ON p.id = sm.product_id
		WHERE 1=1
	`

	args := []interface{}{}
	argNum := 1

	if filter.Search != "" {
		query += fmt.Sprintf(" AND (p.sku ILIKE $%d OR p.name ILIKE $%d)", argNum, argNum+1)
		searchTerm := "%" + filter.Search + "%"
		args = append(args, searchTerm, searchTerm)
		argNum += 2
	}

	if filter.Status == "active" {
		query += fmt.Sprintf(" AND p.is_active = true")
	} else if filter.Status == "inactive" {
		query += fmt.Sprintf(" AND p.is_active = false")
	}

	if filter.LowStock != nil && *filter.LowStock {
		query += fmt.Sprintf(" AND p.stock <= 10")
	}

	query += " GROUP BY p.id ORDER BY p.created_at DESC"

	var totalCount int
	countQuery := "SELECT COUNT(*) FROM products p WHERE 1=1"
	if filter.Search != "" {
		countQuery += fmt.Sprintf(" AND (p.sku ILIKE $1 OR p.name ILIKE $2)")
	}
	if filter.Status == "active" {
		countQuery += " AND p.is_active = true"
	} else if filter.Status == "inactive" {
		countQuery += " AND p.is_active = false"
	}
	if filter.LowStock != nil && *filter.LowStock {
		countQuery += " AND p.stock <= 10"
	}

	if filter.Search != "" {
		r.db.QueryRow(ctx, countQuery, "%" + filter.Search + "%", "%" + filter.Search + "%").Scan(&totalCount)
	} else {
		r.db.QueryRow(ctx, countQuery).Scan(&totalCount)
	}

	offset := (filter.Page - 1) * filter.Limit
	query += fmt.Sprintf(" LIMIT $%d OFFSET $%d", argNum, argNum+1)
	args = append(args, filter.Limit, offset)

	rows, err := r.db.Query(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var products []models.Product
	for rows.Next() {
		var p models.Product
		var lastMovement *time.Time
		if err := rows.Scan(
			&p.ID, &p.SKU, &p.Name, &p.Price, &p.Stock, &p.IsActive, &p.CreatedAt, &lastMovement,
		); err != nil {
			return nil, err
		}
		p.LastMovementAt = lastMovement
		products = append(products, p)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return &models.ProductList{
		Data:  products,
		Page:  filter.Page,
		Limit: filter.Limit,
		Total: totalCount,
	}, nil
}

func (r *ProductRepository) Create(ctx context.Context, req *models.CreateProductRequest) (*models.Product, error) {
	var product models.Product

	query := `
		INSERT INTO products (sku, name, price, stock, is_active)
		VALUES ($1, $2, $3, $4, $5)
		RETURNING id, sku, name, price, stock, is_active, created_at
	`

	err := r.db.QueryRow(ctx, query, req.SKU, req.Name, req.Price, req.Stock, req.IsActive).Scan(
		&product.ID, &product.SKU, &product.Name, &product.Price, &product.Stock,
		&product.IsActive, &product.CreatedAt,
	)

	if err != nil {
		if err.Error() == "ERROR: duplicate key value violates unique constraint \"products_sku_key\" (SQLSTATE 23505)" {
			return nil, errors.New("SKU_ALREADY_EXISTS")
		}
		return nil, err
	}

	return &product, nil
}

func (r *ProductRepository) Update(ctx context.Context, id int, req *models.UpdateProductRequest) (*models.Product, error) {
	var product models.Product
	var lastMovement *time.Time

	query := `
		UPDATE products
		SET sku = $1, name = $2, price = $3, is_active = $4
		WHERE id = $5
		RETURNING id, sku, name, price, stock, is_active, created_at
	`

	err := r.db.QueryRow(ctx, query, req.SKU, req.Name, req.Price, req.IsActive, id).Scan(
		&product.ID, &product.SKU, &product.Name, &product.Price, &product.Stock,
		&product.IsActive, &product.CreatedAt,
	)

	if err == pgx.ErrNoRows {
		return nil, errors.New("product not found")
	}
	if err != nil {
		if err.Error() == "ERROR: duplicate key value violates unique constraint \"products_sku_key\" (SQLSTATE 23505)" {
			return nil, errors.New("SKU_ALREADY_EXISTS")
		}
		return nil, err
	}

	movementQuery := `
		SELECT MAX(created_at)
		FROM stock_movements
		WHERE product_id = $1
	`
	r.db.QueryRow(ctx, movementQuery, id).Scan(&lastMovement)
	product.LastMovementAt = lastMovement

	return &product, nil
}

func (r *ProductRepository) Delete(ctx context.Context, id int) error {
	var movementCount, orderCount int

	movementQuery := "SELECT COUNT(*) FROM stock_movements WHERE product_id = $1"
	err := r.db.QueryRow(ctx, movementQuery, id).Scan(&movementCount)
	if err != nil {
		return err
	}

	orderQuery := "SELECT COUNT(*) FROM sales_order_items soi JOIN sales_orders so ON soi.order_id = so.id WHERE soi.sku = (SELECT sku FROM products WHERE id = $1)"
	err = r.db.QueryRow(ctx, orderQuery, id).Scan(&orderCount)
	if err != nil {
		return err
	}

	if movementCount > 0 || orderCount > 0 {
		return errors.New("PRODUCT_HAS_HISTORY")
	}

	query := "DELETE FROM products WHERE id = $1"
	_, err = r.db.Exec(ctx, query, id)
	return err
}

func (r *ProductRepository) UpdateStock(ctx context.Context, tx pgx.Tx, id int, delta int) error {
	query := `
		UPDATE products
		SET stock = stock + $1
		WHERE id = $2 AND stock + $1 >= 0
	`

	commandTag, err := tx.Exec(ctx, query, delta, id)
	if err != nil {
		return err
	}

	if commandTag.RowsAffected() == 0 {
		return errors.New("INSUFFICIENT_STOCK")
	}

	return nil
}

func (r *ProductRepository) GetStockForUpdate(ctx context.Context, tx pgx.Tx, id int) (int, error) {
	var stock int
	query := "SELECT stock FROM products WHERE id = $1 FOR UPDATE"
	err := tx.QueryRow(ctx, query, id).Scan(&stock)
	return stock, err
}
