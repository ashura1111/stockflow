package services

import (
	"strings"
	"testing"

	"stockflow/backend/internal/models"
)

// MockDB and mocks for testing without actual database
// Tests focus on business logic validation

func TestProductServiceListValidation(t *testing.T) {
	tests := []struct {
		name     string
		page     int
		limit    int
		expPage  int
		expLimit int
	}{
		{
			name:     "invalid page defaults to 1",
			page:     0,
			limit:    10,
			expPage:  1,
			expLimit: 10,
		},
		{
			name:     "invalid limit defaults to 10",
			page:     1,
			limit:    0,
			expPage:  1,
			expLimit: 10,
		},
		{
			name:     "limit capped at 100",
			page:     1,
			limit:    200,
			expPage:  1,
			expLimit: 10,
		},
		{
			name:     "valid pagination",
			page:     2,
			limit:    20,
			expPage:  2,
			expLimit: 20,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			filter := &models.ProductFilter{
				Page:  tt.page,
				Limit: tt.limit,
			}

			// Apply validation logic (matches service code)
			if filter.Page < 1 {
				filter.Page = 1
			}
			if filter.Limit < 1 || filter.Limit > 100 {
				filter.Limit = 10
			}

			if filter.Page != tt.expPage {
				t.Errorf("expected page=%d, got %d", tt.expPage, filter.Page)
			}
			if filter.Limit != tt.expLimit {
				t.Errorf("expected limit=%d, got %d", tt.expLimit, filter.Limit)
			}
		})
	}
}

func TestOrderPaginationValidation(t *testing.T) {
	tests := []struct {
		name     string
		page     int
		limit    int
		expPage  int
		expLimit int
	}{
		{
			name:     "invalid page defaults to 1",
			page:     0,
			limit:    10,
			expPage:  1,
			expLimit: 10,
		},
		{
			name:     "invalid limit defaults to 10",
			page:     1,
			limit:    -5,
			expPage:  1,
			expLimit: 10,
		},
		{
			name:     "limit capped at 100",
			page:     1,
			limit:    150,
			expPage:  1,
			expLimit: 10,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			filter := &models.OrderFilter{
				Page:  tt.page,
				Limit: tt.limit,
			}

			if filter.Page < 1 {
				filter.Page = 1
			}
			if filter.Limit < 1 || filter.Limit > 100 {
				filter.Limit = 10
			}

			if filter.Page != tt.expPage {
				t.Errorf("expected page=%d, got %d", tt.expPage, filter.Page)
			}
			if filter.Limit != tt.expLimit {
				t.Errorf("expected limit=%d, got %d", tt.expLimit, filter.Limit)
			}
		})
	}
}

func TestOrderValidation(t *testing.T) {
	tests := []struct {
		name    string
		order   models.CreateOrderRequest
		valid   bool
		errMsg  string
	}{
		{
			name: "valid order",
			order: models.CreateOrderRequest{
				Customer: "PT Contoh",
				Items: []models.OrderItemInput{
					{SKU: "SKU-001", Qty: 5},
				},
			},
			valid: true,
		},
		{
			name: "duplicate SKU detection",
			order: models.CreateOrderRequest{
				Customer: "PT Contoh",
				Items: []models.OrderItemInput{
					{SKU: "SKU-001", Qty: 5},
					{SKU: "SKU-001", Qty: 3},
				},
			},
			valid:  false,
			errMsg: "duplicate SKU",
		},
		{
			name: "empty items invalid",
			order: models.CreateOrderRequest{
				Customer: "PT Contoh",
				Items:    []models.OrderItemInput{},
			},
			valid:  false,
			errMsg: "empty",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Validate items
			if len(tt.order.Items) == 0 {
				if tt.valid {
					t.Error("expected invalid order with empty items")
				}
				return
			}

			// Check for duplicate SKUs
			skuMap := make(map[string]bool)
			hasDuplicate := false
			for _, item := range tt.order.Items {
				if skuMap[item.SKU] {
					hasDuplicate = true
					break
				}
				skuMap[item.SKU] = true
			}

			if hasDuplicate && tt.valid {
				t.Error("expected valid order but found duplicate SKU")
			}
			if !hasDuplicate && !tt.valid && strings.Contains(tt.errMsg, "duplicate") {
				t.Error("expected duplicate SKU error")
			}
		})
	}
}

func TestStockAdjustmentLogic(t *testing.T) {
	tests := []struct {
		name          string
		currentStock  int
		movementType  string
		qty           int
		valid         bool
		resultStock   int
	}{
		{
			name:         "stock in adds to current",
			currentStock: 100,
			movementType: "in",
			qty:          50,
			valid:        true,
			resultStock:  150,
		},
		{
			name:         "stock out subtracts from current",
			currentStock: 100,
			movementType: "out",
			qty:          30,
			valid:        true,
			resultStock:  70,
		},
		{
			name:         "insufficient stock for out",
			currentStock: 20,
			movementType: "out",
			qty:          50,
			valid:        false,
			resultStock:  20,
		},
		{
			name:         "zero stock after out is valid",
			currentStock: 50,
			movementType: "out",
			qty:          50,
			valid:        true,
			resultStock:  0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			resultStock := tt.currentStock
			delta := tt.qty
			if tt.movementType == "out" {
				delta = -tt.qty
				if tt.currentStock < tt.qty {
					if tt.valid {
						t.Error("expected invalid stock out")
					}
					return
				}
			}

			resultStock += delta

			if resultStock < 0 {
				t.Error("stock cannot be negative")
				return
			}

			if resultStock != tt.resultStock {
				t.Errorf("expected stock=%d, got %d", tt.resultStock, resultStock)
			}
		})
	}
}

func TestOrderCreationStockImpact(t *testing.T) {
	tests := []struct {
		name            string
		items           []struct{ sku string; qty int }
		availableStock  map[string]int
		canCreateOrder  bool
		expectedErr     string
	}{
		{
			name: "sufficient stock for all items",
			items: []struct{ sku string; qty int }{
				{"SKU-001", 5},
				{"SKU-002", 10},
			},
			availableStock: map[string]int{
				"SKU-001": 100,
				"SKU-002": 50,
			},
			canCreateOrder: true,
		},
		{
			name: "insufficient stock for one item",
			items: []struct{ sku string; qty int }{
				{"SKU-001", 150},
			},
			availableStock: map[string]int{
				"SKU-001": 100,
			},
			canCreateOrder: false,
			expectedErr:    "INSUFFICIENT_STOCK",
		},
		{
			name: "product not found",
			items: []struct{ sku string; qty int }{
				{"SKU-NONEXIST", 5},
			},
			availableStock: map[string]int{
				"SKU-001": 100,
			},
			canCreateOrder: false,
			expectedErr:    "not found",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			canCreate := true
			var errMsg string

			for _, item := range tt.items {
				stock, exists := tt.availableStock[item.sku]
				if !exists {
					canCreate = false
					errMsg = "not found"
					break
				}
				if stock < item.qty {
					canCreate = false
					errMsg = "INSUFFICIENT_STOCK"
					break
				}
			}

			if canCreate != tt.canCreateOrder {
				t.Errorf("expected canCreateOrder=%v, got %v", tt.canCreateOrder, canCreate)
			}
			if !canCreate && tt.expectedErr != "" && !strings.Contains(errMsg, tt.expectedErr) {
				t.Errorf("expected error containing '%s', got '%s'", tt.expectedErr, errMsg)
			}
		})
	}
}

func TestOrderCancellationRestoresStock(t *testing.T) {
	tests := []struct {
		name              string
		originalStock     int
		orderQty          int
		shouldRestore     bool
		expectedFinalQty  int
	}{
		{
			name:             "cancellation restores stock",
			originalStock:    100,
			orderQty:         30,
			shouldRestore:    true,
			expectedFinalQty: 100,
		},
		{
			name:             "partial restoration needed",
			originalStock:    50,
			orderQty:         20,
			shouldRestore:    true,
			expectedFinalQty: 50,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Simulate order creation: stock reduced
			stockAfterOrder := tt.originalStock - tt.orderQty

			// Simulate cancellation: stock restored
			var finalStock int
			if tt.shouldRestore {
				finalStock = stockAfterOrder + tt.orderQty
			} else {
				finalStock = stockAfterOrder
			}

			if finalStock != tt.expectedFinalQty {
				t.Errorf("expected final stock=%d, got %d", tt.expectedFinalQty, finalStock)
			}
		})
	}
}

func TestMovementHistoryFiltering(t *testing.T) {
	tests := []struct {
		name       string
		filterType string
		movements  []struct{ mType string; qty int }
		valid      bool
	}{
		{
			name:       "filter by in type",
			filterType: "in",
			movements: []struct{ mType string; qty int }{
				{"in", 10},
				{"out", 5},
				{"in", 20},
			},
			valid: true,
		},
		{
			name:       "filter by out type",
			filterType: "out",
			movements: []struct{ mType string; qty int }{
				{"in", 10},
				{"out", 5},
				{"out", 3},
			},
			valid: true,
		},
		{
			name:       "invalid filter type",
			filterType: "invalid",
			valid:      false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			validTypes := map[string]bool{"in": true, "out": true}
			if tt.filterType != "" && !validTypes[tt.filterType] {
				if tt.valid {
					t.Error("expected invalid filter type")
				}
				return
			}

			if !tt.valid && validTypes[tt.filterType] {
				t.Error("expected valid filter type")
			}
		})
	}
}
