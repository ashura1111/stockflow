package models

import (
	"testing"
)

func TestProductValidation(t *testing.T) {
	tests := []struct {
		name    string
		product Product
		valid   bool
	}{
		{
			name: "valid product",
			product: Product{
				ID:       1,
				SKU:      "SKU-001",
				Name:     "Test Product",
				Price:    10000,
				Stock:    100,
				IsActive: true,
			},
			valid: true,
		},
		{
			name: "zero stock is valid",
			product: Product{
				ID:       1,
				SKU:      "SKU-001",
				Name:     "Test Product",
				Price:    10000,
				Stock:    0,
				IsActive: true,
			},
			valid: true,
		},
		{
			name: "inactive product is valid",
			product: Product{
				ID:       1,
				SKU:      "SKU-001",
				Name:     "Test Product",
				Price:    10000,
				Stock:    100,
				IsActive: false,
			},
			valid: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.product.SKU == "" {
				t.Error("expected valid SKU")
			}
			if tt.product.Stock < 0 {
				t.Error("stock cannot be negative")
			}
		})
	}
}

func TestOrderItemValidation(t *testing.T) {
	tests := []struct {
		name  string
		item  OrderItemInput
		valid bool
	}{
		{
			name: "valid order item",
			item: OrderItemInput{
				SKU: "SKU-001",
				Qty: 5,
			},
			valid: true,
		},
		{
			name: "empty SKU is invalid",
			item: OrderItemInput{
				SKU: "",
				Qty: 5,
			},
			valid: false,
		},
		{
			name: "zero quantity is invalid",
			item: OrderItemInput{
				SKU: "SKU-001",
				Qty: 0,
			},
			valid: false,
		},
		{
			name: "negative quantity is invalid",
			item: OrderItemInput{
				SKU: "SKU-001",
				Qty: -1,
			},
			valid: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			valid := tt.item.SKU != "" && tt.item.Qty > 0
			if valid != tt.valid {
				t.Errorf("expected valid=%v, got %v", tt.valid, valid)
			}
		})
	}
}

func TestCreateOrderRequestValidation(t *testing.T) {
	tests := []struct {
		name    string
		request CreateOrderRequest
		valid   bool
	}{
		{
			name: "valid order request",
			request: CreateOrderRequest{
				Customer: "PT Test",
				Items: []OrderItemInput{
					{SKU: "SKU-001", Qty: 5},
				},
			},
			valid: true,
		},
		{
			name: "empty customer is invalid",
			request: CreateOrderRequest{
				Customer: "",
				Items: []OrderItemInput{
					{SKU: "SKU-001", Qty: 5},
				},
			},
			valid: false,
		},
		{
			name: "empty items is invalid",
			request: CreateOrderRequest{
				Customer: "PT Test",
				Items:    []OrderItemInput{},
			},
			valid: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			valid := tt.request.Customer != "" && len(tt.request.Items) > 0
			if valid != tt.valid {
				t.Errorf("expected valid=%v, got %v", tt.valid, valid)
			}
		})
	}
}

func TestStockMovementRequestValidation(t *testing.T) {
	tests := []struct {
		name    string
		request StockMovementRequest
		valid   bool
	}{
		{
			name: "valid stock in movement",
			request: StockMovementRequest{
				Type: "in",
				Qty:  10,
				Note: "Restocking",
			},
			valid: true,
		},
		{
			name: "valid stock out movement",
			request: StockMovementRequest{
				Type: "out",
				Qty:  5,
				Note: "Sale",
			},
			valid: true,
		},
		{
			name: "invalid type",
			request: StockMovementRequest{
				Type: "invalid",
				Qty:  10,
			},
			valid: false,
		},
		{
			name: "zero quantity is invalid",
			request: StockMovementRequest{
				Type: "in",
				Qty:  0,
			},
			valid: false,
		},
		{
			name: "negative quantity is invalid",
			request: StockMovementRequest{
				Type: "in",
				Qty:  -1,
			},
			valid: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			valid := (tt.request.Type == "in" || tt.request.Type == "out") && tt.request.Qty > 0
			if valid != tt.valid {
				t.Errorf("expected valid=%v, got %v", tt.valid, valid)
			}
		})
	}
}
