package models

import (
	"time"
)

type Product struct {
	ID             int       `json:"id"`
	SKU            string    `json:"sku"`
	Name           string    `json:"name"`
	Price          int       `json:"price"`
	Stock          int       `json:"stock"`
	IsActive       bool      `json:"is_active"`
	CreatedAt      time.Time `json:"created_at"`
	LastMovementAt *time.Time `json:"last_movement_at"`
}

type ProductList struct {
	Data  []Product `json:"data"`
	Page  int       `json:"page"`
	Limit int       `json:"limit"`
	Total int       `json:"total"`
}

type StockMovement struct {
	ID         int       `json:"id"`
	ProductID  int       `json:"product_id"`
	SKU        string    `json:"sku"`
	ProductName string    `json:"product_name"`
	Type       string    `json:"type"`
	Qty        int       `json:"qty"`
	Note       string    `json:"note"`
	CreatedAt  time.Time `json:"created_at"`
}

type StockMovementRequest struct {
	Type string `json:"type" validate:"required,oneof=in out"`
	Qty  int    `json:"qty" validate:"required,min=1"`
	Note string `json:"note"`
}

type SalesOrder struct {
	ID           int            `json:"id"`
	Customer     string         `json:"customer"`
	Status       string         `json:"status"`
	Total        int            `json:"total"`
	CreatedAt    time.Time      `json:"created_at"`
	CancelledAt  *time.Time     `json:"cancelled_at"`
	Items        []OrderItem    `json:"items"`
}

type OrderItem struct {
	ID        int    `json:"id"`
	OrderID   int    `json:"order_id"`
	SKU       string `json:"sku"`
	ProductName string `json:"product_name"`
	Qty       int    `json:"qty"`
	UnitPrice int    `json:"unit_price"`
	Subtotal  int    `json:"subtotal"`
}

type CreateOrderRequest struct {
	Customer string       `json:"customer" validate:"required,min=1"`
	Items    []OrderItemInput `json:"items" validate:"required,min=1,dive"`
}

type OrderItemInput struct {
	SKU string `json:"sku" validate:"required,min=1"`
	Qty int    `json:"qty" validate:"required,min=1"`
}

type OrderCancelResponse struct {
	ID          int        `json:"id"`
	Customer    string     `json:"customer"`
	Status      string     `json:"status"`
	Total       int        `json:"total"`
	CreatedAt   time.Time  `json:"created_at"`
	CancelledAt *time.Time `json:"cancelled_at"`
}

type ProductSummary struct {
	TotalProducts    int `json:"total_products"`
	ActiveProducts   int `json:"active_products"`
	TotalStockValue  int `json:"total_stock_value"`
	LowStockProducts int `json:"low_stock_products"`
}

type ErrorResponse struct {
	Error struct {
		Code    string `json:"code"`
		Message string `json:"message"`
		SKU     string `json:"sku,omitempty"`
		Available int    `json:"available,omitempty"`
	} `json:"error"`
}
