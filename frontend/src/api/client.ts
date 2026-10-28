import type {
  CreateOrderRequest,
  CreateProductRequest,
  DashboardStats,
  ErrorResponse,
  OrderList,
  ProductList,
  SalesOrder,
  StockMovementList,
  StockMovementRequest,
  UpdateProductRequest,
} from '../types/index';

const API_BASE = '/api';

class ApiClient {
  async request<T>(
    method: string,
    endpoint: string,
    body?: unknown
  ): Promise<T> {
    const url = `${API_BASE}${endpoint}`;
    const options: RequestInit = {
      method,
      headers: {
        'Content-Type': 'application/json',
      },
    };

    if (body) {
      options.body = JSON.stringify(body);
    }

    const response = await fetch(url, options);

    if (!response.ok) {
      const errorData = await response.json() as ErrorResponse;
      throw new ApiError(errorData.error.message, errorData.error.code, errorData.error);
    }

    if (response.status === 204) {
      return undefined as T;
    }

    return response.json() as Promise<T>;
  }

  async get<T>(endpoint: string): Promise<T> {
    return this.request<T>('GET', endpoint);
  }

  async post<T>(endpoint: string, body?: unknown): Promise<T> {
    return this.request<T>('POST', endpoint, body);
  }

  async put<T>(endpoint: string, body?: unknown): Promise<T> {
    return this.request<T>('PUT', endpoint, body);
  }

  async delete<T>(endpoint: string): Promise<T> {
    return this.request<T>('DELETE', endpoint);
  }

  health() {
    return this.get('/health');
  }

  dashboard() {
    return this.get<DashboardStats>('/dashboard');
  }

  getProducts(search: string = '', status: string = '', lowStock: boolean = false, page: number = 1, limit: number = 10) {
    const params = new URLSearchParams();
    if (search) params.append('search', search);
    if (status) params.append('status', status);
    if (lowStock) params.append('low_stock', 'true');
    params.append('page', page.toString());
    params.append('limit', limit.toString());
    return this.get<ProductList>(`/products?${params}`);
  }

  getProduct(id: number) {
    return this.get(`/products/${id}`);
  }

  createProduct(req: CreateProductRequest) {
    return this.post('/products', req);
  }

  updateProduct(id: number, req: UpdateProductRequest) {
    return this.put(`/products/${id}`, req);
  }

  deleteProduct(id: number) {
    return this.delete(`/products/${id}`);
  }

  adjustStock(id: number, req: StockMovementRequest) {
    return this.post(`/products/${id}/stock`, req);
  }

  getMovementHistory(id: number, page: number = 1, limit: number = 10) {
    const params = new URLSearchParams();
    params.append('page', page.toString());
    params.append('limit', limit.toString());
    return this.get<StockMovementList>(`/products/${id}/movements?${params}`);
  }

  getMovements(search: string = '', type: string = '', page: number = 1, limit: number = 20) {
    const params = new URLSearchParams();
    if (search) params.append('search', search);
    if (type) params.append('type', type);
    params.append('page', page.toString());
    params.append('limit', limit.toString());
    return this.get<StockMovementList>(`/movements?${params}`);
  }

  getOrders(search: string = '', status: string = '', page: number = 1, limit: number = 10) {
    const params = new URLSearchParams();
    if (search) params.append('search', search);
    if (status) params.append('status', status);
    params.append('page', page.toString());
    params.append('limit', limit.toString());
    return this.get<OrderList>(`/orders?${params}`);
  }

  getOrder(id: number) {
    return this.get<SalesOrder>(`/orders/${id}`);
  }

  createOrder(req: CreateOrderRequest) {
    return this.post<SalesOrder>('/orders', req);
  }

  cancelOrder(id: number) {
    return this.post(`/orders/${id}/cancel`);
  }
}

export class ApiError extends Error {
  code: string;
  details: {
    message: string;
    code: string;
    sku?: string;
    available?: number;
  };

  constructor(
    message: string,
    code: string,
    details: { message: string; code: string; sku?: string; available?: number }
  ) {
    super(message);
    this.name = 'ApiError';
    this.code = code;
    this.details = details;
  }
}

export const api = new ApiClient();
