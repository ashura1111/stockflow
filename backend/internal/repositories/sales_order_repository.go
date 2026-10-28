package repositories

import (
	"context"
	"errors"
	"fmt"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"stockflow/backend/internal/models"
)

type SalesOrderRepository struct {
	db *pgxpool.Pool
}

func NewSalesOrderRepository(db *pgxpool.Pool) *SalesOrderRepository {
	return &SalesOrderRepository{db: db}
}

func (r *SalesOrderRepository) GetByID(ctx context.Context, id int) (*models.SalesOrder, error) {
	query := `
		SELECT id, customer, status, total, created_at, cancelled_at
		FROM sales_orders
		WHERE id = $1
	`

	var order models.SalesOrder
	err := r.db.QueryRow(ctx, query, id).Scan(
		&order.ID, &order.Customer, &order.Status, &order.Total, &order.CreatedAt, &order.CancelledAt,
	)

	if err == pgx.ErrNoRows {
		return nil, errors.New("order not found")
	}
	if err != nil {
		return nil, err
	}

	itemsQuery := `
		SELECT id, order_id, sku, qty, unit_price
		FROM sales_order_items
		WHERE order_id = $1
	`

	itemRows, err := r.db.Query(ctx, itemsQuery, id)
	if err != nil {
		return nil, err
	}
	defer itemRows.Close()

	var items []models.OrderItem
	for itemRows.Next() {
		var item models.OrderItem
		if err := itemRows.Scan(&item.ID, &item.OrderID, &item.SKU, &item.Qty, &item.UnitPrice); err != nil {
			return nil, err
		}
		item.Subtotal = item.Qty * item.UnitPrice

		productQuery := "SELECT name FROM products WHERE sku = $1"
		_ = r.db.QueryRow(ctx, productQuery, item.SKU).Scan(&item.ProductName)

		items = append(items, item)
	}

	order.Items = items
	return &order, nil
}

func (r *SalesOrderRepository) List(ctx context.Context, filter *models.OrderFilter) (*models.OrderList, error) {
	query := `
		SELECT id, customer, status, total, created_at, cancelled_at
		FROM sales_orders
		WHERE 1=1
	`

	args := []interface{}{}
	argNum := 1

	if filter.Search != "" {
		query += fmt.Sprintf(" AND customer ILIKE $%d", argNum)
		args = append(args, "%"+filter.Search+"%")
		argNum++
	}

	if filter.Status != "" {
		query += fmt.Sprintf(" AND status = $%d", argNum)
		args = append(args, filter.Status)
		argNum++
	}

	query += " ORDER BY created_at DESC"

	var totalCount int
	countQuery := "SELECT COUNT(*) FROM sales_orders WHERE 1=1"
	if filter.Search != "" {
		countQuery += " AND customer ILIKE $1"
		r.db.QueryRow(ctx, countQuery, "%"+filter.Search+"%").Scan(&totalCount)
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

	var orders []models.SalesOrder
	for rows.Next() {
		var order models.SalesOrder
		if err := rows.Scan(
			&order.ID, &order.Customer, &order.Status, &order.Total, &order.CreatedAt, &order.CancelledAt,
		); err != nil {
			return nil, err
		}
		orders = append(orders, order)
	}

	return &models.OrderList{
		Data:  orders,
		Page:  filter.Page,
		Limit: filter.Limit,
		Total: totalCount,
	}, nil
}

func (r *SalesOrderRepository) Create(ctx context.Context, tx pgx.Tx, customer string, total int) (int, error) {
	query := `
		INSERT INTO sales_orders (customer, status, total)
		VALUES ($1, 'CREATED', $2)
		RETURNING id
	`

	var id int
	err := tx.QueryRow(ctx, query, customer, total).Scan(&id)
	return id, err
}

func (r *SalesOrderRepository) AddItem(ctx context.Context, tx pgx.Tx, orderID int, sku string, qty, unitPrice int) error {
	query := `
		INSERT INTO sales_order_items (order_id, sku, qty, unit_price)
		VALUES ($1, $2, $3, $4)
	`
	_, err := tx.Exec(ctx, query, orderID, sku, qty, unitPrice)
	return err
}

func (r *SalesOrderRepository) Cancel(ctx context.Context, tx pgx.Tx, id int) error {
	query := `
		UPDATE sales_orders
		SET status = 'CANCELLED', cancelled_at = NOW()
		WHERE id = $1 AND status = 'CREATED'
	`

	commandTag, err := tx.Exec(ctx, query, id)
	if err != nil {
		return err
	}

	if commandTag.RowsAffected() == 0 {
		return errors.New("order not found or already cancelled")
	}

	return nil
}

func (r *SalesOrderRepository) GetRecent(ctx context.Context, limit int) ([]models.SalesOrder, error) {
	query := `
		SELECT id, customer, status, total, created_at, cancelled_at
		FROM sales_orders
		ORDER BY created_at DESC
		LIMIT $1
	`

	rows, err := r.db.Query(ctx, query, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var orders []models.SalesOrder
	for rows.Next() {
		var order models.SalesOrder
		if err := rows.Scan(
			&order.ID, &order.Customer, &order.Status, &order.Total, &order.CreatedAt, &order.CancelledAt,
		); err != nil {
			return nil, err
		}
		orders = append(orders, order)
	}

	return orders, rows.Err()
}

func (r *SalesOrderRepository) GetOrderItems(ctx context.Context, tx pgx.Tx, orderID int) ([]models.OrderItem, error) {
	query := `
		SELECT id, order_id, sku, qty, unit_price
		FROM sales_order_items
		WHERE order_id = $1
	`

	rows, err := tx.Query(ctx, query, orderID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var items []models.OrderItem
	for rows.Next() {
		var item models.OrderItem
		if err := rows.Scan(&item.ID, &item.OrderID, &item.SKU, &item.Qty, &item.UnitPrice); err != nil {
			return nil, err
		}
		items = append(items, item)
	}

	return items, rows.Err()
}
