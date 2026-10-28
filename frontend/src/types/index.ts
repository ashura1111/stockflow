export interface Product {
  id: number;
  sku: string;
  name: string;
  price: number;
  stock: number;
  is_active: boolean;
  created_at: string;
  last_movement_at: string | null;
}

export interface ProductList {
  data: Product[];
  page: number;
  limit: number;
  total: number;
}

export interface StockMovement {
  id: number;
  product_id: number;
  sku: string;
  product_name: string;
  type: 'in' | 'out';
  qty: number;
  note: string;
  created_at: string;
}

export interface StockMovementList {
  data: StockMovement[];
  page: number;
  limit: number;
  total: number;
}

export interface OrderItem {
  id: number;
  order_id: number;
  sku: string;
  product_name: string;
  qty: number;
  unit_price: number;
  subtotal: number;
}

export interface SalesOrder {
  id: number;
  customer: string;
  status: 'CREATED' | 'CANCELLED';
  total: number;
  created_at: string;
  cancelled_at: string | null;
  items: OrderItem[];
}

export interface OrderList {
  data: SalesOrder[];
  page: number;
  limit: number;
  total: number;
}

export interface DashboardStats {
  total_products: number;
  active_products: number;
  total_stock_value: number;
  low_stock_products: number;
  recent_movements: StockMovement[];
  recent_orders: SalesOrder[];
  low_stock_list: Product[];
}

export interface ErrorResponse {
  error: {
    code: string;
    message: string;
    sku?: string;
    available?: number;
  };
}

export interface CreateProductRequest {
  sku: string;
  name: string;
  price: number;
  stock: number;
  is_active: boolean;
}

export interface UpdateProductRequest {
  sku: string;
  name: string;
  price: number;
  is_active: boolean;
}

export interface StockMovementRequest {
  type: 'in' | 'out';
  qty: number;
  note: string;
}

export interface CreateOrderRequest {
  customer: string;
  items: OrderItemInput[];
}

export interface OrderItemInput {
  sku: string;
  qty: number;
}
