package services

import (
	"context"
	"errors"
	"fmt"

	"github.com/jackc/pgx/v5/pgxpool"
	"stockflow/backend/internal/models"
	"stockflow/backend/internal/repositories"
	"stockflow/backend/internal/validation"
)

type SalesOrderService struct {
	db          *pgxpool.Pool
	productRepo *repositories.ProductRepository
	orderRepo   *repositories.SalesOrderRepository
	movementRepo *repositories.StockMovementRepository
}

func NewSalesOrderService(db *pgxpool.Pool) *SalesOrderService {
	return &SalesOrderService{
		db:          db,
		productRepo: repositories.NewProductRepository(db),
		orderRepo:   repositories.NewSalesOrderRepository(db),
		movementRepo: repositories.NewStockMovementRepository(db),
	}
}

func (s *SalesOrderService) GetOrder(ctx context.Context, id int) (*models.SalesOrder, error) {
	return s.orderRepo.GetByID(ctx, id)
}

func (s *SalesOrderService) ListOrders(ctx context.Context, filter *models.OrderFilter) (*models.OrderList, error) {
	if filter.Page < 1 {
		filter.Page = 1
	}
	if filter.Limit < 1 || filter.Limit > 100 {
		filter.Limit = 10
	}
	return s.orderRepo.List(ctx, filter)
}

func (s *SalesOrderService) CreateOrder(ctx context.Context, req *models.CreateOrderRequest) (*models.SalesOrder, error) {
	if err := validation.ValidateStruct(req); err != nil {
		return nil, errors.New(validation.FormatValidationError(err))
	}

	if len(req.Items) == 0 {
		return nil, errors.New("items cannot be empty")
	}

	skuMap := make(map[string]bool)
	for _, item := range req.Items {
		if skuMap[item.SKU] {
			return nil, fmt.Errorf("duplicate SKU in order: %s", item.SKU)
		}
		skuMap[item.SKU] = true
	}

	tx, err := s.db.Begin(ctx)
	if err != nil {
		return nil, err
	}
	defer tx.Rollback(ctx)

	var total int
	orderItems := make([]models.OrderItem, 0)

	for _, item := range req.Items {
		product, err := s.productRepo.GetBySKU(ctx, item.SKU)
		if err != nil {
			return nil, fmt.Errorf("product not found: %s", item.SKU)
		}

		if !product.IsActive {
			return nil, fmt.Errorf("inactive product cannot be ordered: %s", item.SKU)
		}

		if product.Stock < item.Qty {
			return nil, fmt.Errorf("INSUFFICIENT_STOCK:%s:%d", item.SKU, product.Stock)
		}

		lineTotal := item.Qty * product.Price
		total += lineTotal

		orderItems = append(orderItems, models.OrderItem{
			SKU:       item.SKU,
			Qty:       item.Qty,
			UnitPrice: product.Price,
			Subtotal:  lineTotal,
			ProductName: product.Name,
		})
	}

	orderID, err := s.orderRepo.Create(ctx, tx, req.Customer, total)
	if err != nil {
		return nil, err
	}

	for _, item := range orderItems {
		product, _ := s.productRepo.GetBySKU(ctx, item.SKU)
		if err := s.orderRepo.AddItem(ctx, tx, orderID, item.SKU, item.Qty, item.UnitPrice); err != nil {
			return nil, err
		}

		if err := s.productRepo.UpdateStock(ctx, tx, product.ID, -item.Qty); err != nil {
			return nil, fmt.Errorf("INSUFFICIENT_STOCK:%s:%d", item.SKU, product.Stock)
		}

		if err := s.movementRepo.Create(ctx, tx, product.ID, "out", item.Qty, fmt.Sprintf("Sales order %d", orderID)); err != nil {
			return nil, err
		}
	}

	if err := tx.Commit(ctx); err != nil {
		return nil, err
	}

	return s.orderRepo.GetByID(ctx, orderID)
}

func (s *SalesOrderService) CancelOrder(ctx context.Context, id int) (*models.OrderCancelResponse, error) {
	order, err := s.orderRepo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}

	if order.Status == "CANCELLED" {
		return nil, errors.New("order already cancelled")
	}

	tx, err := s.db.Begin(ctx)
	if err != nil {
		return nil, err
	}
	defer tx.Rollback(ctx)

	if err := s.orderRepo.Cancel(ctx, tx, id); err != nil {
		return nil, err
	}

	items, err := s.orderRepo.GetOrderItems(ctx, tx, id)
	if err != nil {
		return nil, err
	}

	for _, item := range items {
		product, err := s.productRepo.GetBySKU(ctx, item.SKU)
		if err != nil {
			return nil, err
		}

		if err := s.productRepo.UpdateStock(ctx, tx, product.ID, item.Qty); err != nil {
			return nil, err
		}

		if err := s.movementRepo.Create(ctx, tx, product.ID, "in", item.Qty, fmt.Sprintf("Cancelled sales order %d", id)); err != nil {
			return nil, err
		}
	}

	if err := tx.Commit(ctx); err != nil {
		return nil, err
	}

	updated, _ := s.orderRepo.GetByID(ctx, id)

	return &models.OrderCancelResponse{
		ID:          updated.ID,
		Customer:    updated.Customer,
		Status:      updated.Status,
		Total:       updated.Total,
		CreatedAt:   updated.CreatedAt,
		CancelledAt: updated.CancelledAt,
	}, nil
}

func (s *SalesOrderService) GetDashboardStats(ctx context.Context) (*models.DashboardStats, error) {
	products, err := s.productRepo.List(ctx, &models.ProductFilter{
		Page:  1,
		Limit: 10000,
	})
	if err != nil {
		return nil, err
	}

	stats := &models.DashboardStats{
		TotalProducts:    products.Total,
		RecentMovements:  make([]models.StockMovement, 0),
		RecentOrders:     make([]models.SalesOrder, 0),
		LowStockList:     make([]models.Product, 0),
	}

	activeCount := 0
	totalValue := 0
	lowStockCount := 0

	for _, p := range products.Data {
		if p.IsActive {
			activeCount++
		}
		if p.Stock <= 10 {
			lowStockCount++
			stats.LowStockList = append(stats.LowStockList, p)
		}
		totalValue += p.Price * p.Stock
	}

	stats.ActiveProducts = activeCount
	stats.TotalStockValue = totalValue
	stats.LowStockProducts = lowStockCount

	recentMovements, err := s.movementRepo.GetRecent(ctx, 5)
	if err == nil {
		stats.RecentMovements = recentMovements
	}

	recentOrders, err := s.orderRepo.GetRecent(ctx, 5)
	if err == nil {
		stats.RecentOrders = recentOrders
	}

	return stats, nil
}
