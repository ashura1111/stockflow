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

type ProductService struct {
	db         *pgxpool.Pool
	productRepo *repositories.ProductRepository
	movementRepo *repositories.StockMovementRepository
}

func NewProductService(db *pgxpool.Pool) *ProductService {
	return &ProductService{
		db:         db,
		productRepo: repositories.NewProductRepository(db),
		movementRepo: repositories.NewStockMovementRepository(db),
	}
}

func (s *ProductService) GetProduct(ctx context.Context, id int) (*models.Product, error) {
	return s.productRepo.GetByID(ctx, id)
}

func (s *ProductService) ListProducts(ctx context.Context, filter *models.ProductFilter) (*models.ProductList, error) {
	if filter.Page < 1 {
		filter.Page = 1
	}
	if filter.Limit < 1 || filter.Limit > 100 {
		filter.Limit = 10
	}
	return s.productRepo.List(ctx, filter)
}

func (s *ProductService) CreateProduct(ctx context.Context, req *models.CreateProductRequest) (*models.Product, error) {
	if err := validation.ValidateStruct(req); err != nil {
		return nil, errors.New(validation.FormatValidationError(err))
	}

	product, err := s.productRepo.Create(ctx, req)
	if err != nil {
		if err.Error() == "SKU_ALREADY_EXISTS" {
			return nil, errors.New("SKU_ALREADY_EXISTS")
		}
		return nil, err
	}

	if req.Stock > 0 {
		tx, err := s.db.Begin(ctx)
		if err != nil {
			return nil, err
		}
		defer tx.Rollback(ctx)

		if err := s.movementRepo.Create(ctx, tx, product.ID, "in", req.Stock, "Initial stock"); err != nil {
			return nil, err
		}

		if err := tx.Commit(ctx); err != nil {
			return nil, err
		}
	}

	return product, nil
}

func (s *ProductService) UpdateProduct(ctx context.Context, id int, req *models.UpdateProductRequest) (*models.Product, error) {
	if err := validation.ValidateStruct(req); err != nil {
		return nil, errors.New(validation.FormatValidationError(err))
	}

	product, err := s.productRepo.Update(ctx, id, req)
	if err != nil {
		if err.Error() == "SKU_ALREADY_EXISTS" {
			return nil, errors.New("SKU_ALREADY_EXISTS")
		}
		if err.Error() == "product not found" {
			return nil, err
		}
		return nil, err
	}

	return product, nil
}

func (s *ProductService) DeleteProduct(ctx context.Context, id int) error {
	return s.productRepo.Delete(ctx, id)
}

func (s *ProductService) AdjustStock(ctx context.Context, id int, req *models.StockMovementRequest) (*models.Product, error) {
	if err := validation.ValidateStruct(req); err != nil {
		return nil, errors.New(validation.FormatValidationError(err))
	}

	tx, err := s.db.Begin(ctx)
	if err != nil {
		return nil, err
	}
	defer tx.Rollback(ctx)

	product, err := s.productRepo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}

	delta := req.Qty
	if req.Type == "out" {
		delta = -req.Qty
		if product.Stock < req.Qty {
			return nil, fmt.Errorf("INSUFFICIENT_STOCK:%d", product.Stock)
		}
	}

	if err := s.productRepo.UpdateStock(ctx, tx, id, delta); err != nil {
		if err.Error() == "INSUFFICIENT_STOCK" {
			return nil, fmt.Errorf("INSUFFICIENT_STOCK:%d", product.Stock)
		}
		return nil, err
	}

	if err := s.movementRepo.Create(ctx, tx, id, req.Type, req.Qty, req.Note); err != nil {
		return nil, err
	}

	if err := tx.Commit(ctx); err != nil {
		return nil, err
	}

	updated, err := s.productRepo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}

	return updated, nil
}

func (s *ProductService) GetMovementHistory(ctx context.Context, productID int, filter *models.StockMovementFilter) (*models.StockMovementList, error) {
	if filter.Page < 1 {
		filter.Page = 1
	}
	if filter.Limit < 1 || filter.Limit > 100 {
		filter.Limit = 10
	}
	filter.ProductID = productID
	return s.movementRepo.GetByProductID(ctx, filter)
}
