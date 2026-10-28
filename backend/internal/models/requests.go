package models

type CreateProductRequest struct {
	SKU      string `json:"sku" validate:"required,min=1"`
	Name     string `json:"name" validate:"required,min=1"`
	Price    int    `json:"price" validate:"required,min=0"`
	Stock    int    `json:"stock" validate:"required,min=0"`
	IsActive bool   `json:"is_active"`
}

type UpdateProductRequest struct {
	SKU      string `json:"sku" validate:"required,min=1"`
	Name     string `json:"name" validate:"required,min=1"`
	Price    int    `json:"price" validate:"required,min=0"`
	IsActive bool   `json:"is_active"`
}

type StockMovementList struct {
	Data  []StockMovement `json:"data"`
	Page  int             `json:"page"`
	Limit int             `json:"limit"`
	Total int             `json:"total"`
}

type OrderList struct {
	Data  []SalesOrder `json:"data"`
	Page  int          `json:"page"`
	Limit int          `json:"limit"`
	Total int          `json:"total"`
}

type ProductFilter struct {
	Search   string
	Status   string
	LowStock *bool
	Page     int
	Limit    int
}

type StockMovementFilter struct {
	ProductID int
	Type      string
	Search    string
	Page      int
	Limit     int
}

type OrderFilter struct {
	Search string
	Status string
	Page   int
	Limit  int
}

type DashboardStats struct {
	TotalProducts    int             `json:"total_products"`
	ActiveProducts   int             `json:"active_products"`
	TotalStockValue  int             `json:"total_stock_value"`
	LowStockProducts int             `json:"low_stock_products"`
	RecentMovements  []StockMovement `json:"recent_movements"`
	RecentOrders     []SalesOrder    `json:"recent_orders"`
	LowStockList     []Product       `json:"low_stock_list"`
}
