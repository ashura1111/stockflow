package repositories

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"stockflow/backend/internal/models"
)

type StockMovementRepository struct {
	db *pgxpool.Pool
}

func NewStockMovementRepository(db *pgxpool.Pool) *StockMovementRepository {
	return &StockMovementRepository{db: db}
}

func (r *StockMovementRepository) Create(ctx context.Context, tx pgx.Tx, productID int, movementType string, qty int, note string) error {
	query := `
        INSERT INTO stock_movements (product_id, type, qty, note)
        VALUES ($1, $2, $3, $4)
    `
	_, err := tx.Exec(ctx, query, productID, movementType, qty, note)
	return err
}

func (r *StockMovementRepository) GetByProductID(ctx context.Context, filter *models.StockMovementFilter) (*models.StockMovementList, error) {
	where := " WHERE 1=1"
	args := []interface{}{}
	argNum := 1

	if filter.ProductID > 0 {
		where += fmt.Sprintf(" AND sm.product_id = $%d", argNum)
		args = append(args, filter.ProductID)
		argNum++
	}

	if filter.Type != "" {
		where += fmt.Sprintf(" AND sm.type = $%d", argNum)
		args = append(args, filter.Type)
		argNum++
	}

	if filter.Search != "" {
		where += fmt.Sprintf(" AND (p.sku ILIKE $%d OR p.name ILIKE $%d)", argNum, argNum)
		args = append(args, "%"+filter.Search+"%")
		argNum++
	}

	countQuery := `
        SELECT COUNT(*)
        FROM stock_movements sm
        JOIN products p ON sm.product_id = p.id
    ` + where

	var totalCount int
	if err := r.db.QueryRow(ctx, countQuery, args...).Scan(&totalCount); err != nil {
		return nil, err
	}

	query := `
        SELECT sm.id, sm.product_id, p.sku, p.name,
               sm.type, sm.qty, sm.note, sm.created_at
        FROM stock_movements sm
        JOIN products p ON sm.product_id = p.id
    ` + where + " ORDER BY sm.created_at DESC, sm.id DESC"

	offset := (filter.Page - 1) * filter.Limit
	query += fmt.Sprintf(" LIMIT $%d OFFSET $%d", argNum, argNum+1)
	args = append(args, filter.Limit, offset)

	rows, err := r.db.Query(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	movements := []models.StockMovement{}
	for rows.Next() {
		var m models.StockMovement
		if err := rows.Scan(
			&m.ID, &m.ProductID, &m.SKU, &m.ProductName,
			&m.Type, &m.Qty, &m.Note, &m.CreatedAt,
		); err != nil {
			return nil, err
		}
		movements = append(movements, m)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return &models.StockMovementList{
		Data:  movements,
		Page:  filter.Page,
		Limit: filter.Limit,
		Total: totalCount,
	}, nil
}

func (r *StockMovementRepository) GetRecent(ctx context.Context, limit int) ([]models.StockMovement, error) {
	query := `
        SELECT sm.id, sm.product_id, p.sku, p.name,
               sm.type, sm.qty, sm.note, sm.created_at
        FROM stock_movements sm
        JOIN products p ON sm.product_id = p.id
        ORDER BY sm.created_at DESC, sm.id DESC
        LIMIT $1
    `

	rows, err := r.db.Query(ctx, query, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	movements := []models.StockMovement{}
	for rows.Next() {
		var m models.StockMovement
		if err := rows.Scan(
			&m.ID, &m.ProductID, &m.SKU, &m.ProductName,
			&m.Type, &m.Qty, &m.Note, &m.CreatedAt,
		); err != nil {
			return nil, err
		}
		movements = append(movements, m)
	}

	return movements, rows.Err()
}
