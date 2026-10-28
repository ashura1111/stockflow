# StockFlow - Inventory Management System

A complete, production-minded full-stack inventory management web application for construction materials and home supplies businesses. Built with Go, React, PostgreSQL, and Docker.

## Features

- **Product Management**: Create, update, delete products with SKU, pricing, and stock tracking
- **Stock Movements**: Track all inventory adjustments (stock in/out) with historical records
- **Sales Orders**: Create customer orders with automatic stock deduction and order cancellation
- **Dashboard**: Real-time statistics including total products, active products, total stock value, and low-stock alerts
- **Search & Filtering**: Advanced search by SKU/product name with status and low-stock filters
- **Pagination**: Efficient data loading with customizable page sizes
- **Real-time Updates**: All data refreshes after mutations
- **Error Handling**: Comprehensive validation and user-friendly error messages
- **Currency Formatting**: All prices displayed in Indonesian Rupiah (IDR)

## Technology Stack

### Backend
- **Language**: Go 1.22
- **HTTP Router**: Go standard `net/http` with custom routing
- **Database**: PostgreSQL 17
- **Database Driver**: pgx/v5
- **Validation**: go-playground/validator/v10

### Frontend
- **Framework**: React 19 with TypeScript
- **Build Tool**: Vite
- **Routing**: React Router v6
- **Icons**: Lucide React
- **Styling**: Tailwind-inspired CSS

### Infrastructure
- **Containerization**: Docker
- **Orchestration**: Docker Compose
- **Database Volume**: Persistent PostgreSQL data

## Project Structure

```
stockflow/
├── backend/
│   ├── cmd/api/
│   │   └── main.go                 # Application entry point
│   ├── internal/
│   │   ├── config/                 # Configuration management
│   │   ├── database/               # Database connection
│   │   ├── handlers/               # HTTP handlers
│   │   ├── middleware/             # HTTP middleware
│   │   ├── models/                 # Data models
│   │   ├── repositories/           # Data access layer
│   │   ├── services/               # Business logic
│   │   └── validation/             # Input validation
│   ├── migrations/                 # SQL migrations
│   ├── go.mod, go.sum             # Dependencies
│   └── Dockerfile
├── frontend/
│   ├── src/
│   │   ├── api/                    # API client
│   │   ├── components/             # Reusable components
│   │   ├── pages/                  # Page components
│   │   ├── types/                  # TypeScript types
│   │   ├── utils/                  # Utility functions
│   │   ├── App.tsx, main.tsx       # Application entry
│   │   ├── App.css, index.css      # Styling
│   ├── index.html, vite.config.ts
│   ├── package.json
│   └── Dockerfile
├── docker-compose.yml              # Service orchestration
└── README.md
```

## Prerequisites

- Docker and Docker Compose (for containerized setup)
- Node.js 20+ (for frontend development)
- Go 1.22+ (for backend development)
- PostgreSQL 17+ (for local development without Docker)

## Quick Start with Docker Compose

This is the recommended way to run the complete stack.

```bash
cd C:\stockflow
docker compose up --build
```

The application will be available at:
- **Frontend**: http://localhost:5173
- **Backend API**: http://localhost:8080
- **PostgreSQL**: localhost:5432
- **Health Check**: http://localhost:8080/api/health

The database will automatically initialize with schema and seed data on first run.

## Local Development Setup

### Prerequisites

Ensure you have Docker and Node.js installed.

### 1. Start PostgreSQL

```bash
docker run -d \
  --name stockflow-postgres \
  -e POSTGRES_DB=stockflow \
  -e POSTGRES_USER=stockflow \
  -e POSTGRES_PASSWORD=stockflow_dev_password \
  -p 5432:5432 \
  -v stockflow_data:/var/lib/postgresql/data \
  -v "$(pwd)/backend/migrations:/docker-entrypoint-initdb.d:ro" \
  postgres:17
```

### 2. Start Backend

```bash
cd backend
go run ./cmd/api
```

Backend runs on http://localhost:8080

### 3. Start Frontend

In a new terminal:

```bash
cd frontend
npm install --legacy-peer-deps
npm run dev
```

Frontend runs on http://localhost:5173

The Vite dev server proxies `/api/*` requests to `http://localhost:8080/api/*`.

## Database Schema

### Products
```sql
CREATE TABLE products (
    id SERIAL PRIMARY KEY,
    sku TEXT NOT NULL UNIQUE,
    name TEXT NOT NULL,
    price INTEGER NOT NULL CHECK (price >= 0),
    stock INTEGER NOT NULL DEFAULT 0 CHECK (stock >= 0),
    is_active BOOLEAN NOT NULL DEFAULT TRUE,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);
```

### Stock Movements
```sql
CREATE TABLE stock_movements (
    id SERIAL PRIMARY KEY,
    product_id INTEGER NOT NULL REFERENCES products(id),
    type TEXT NOT NULL CHECK (type IN ('in', 'out')),
    qty INTEGER NOT NULL CHECK (qty > 0),
    note TEXT,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);
```

### Sales Orders
```sql
CREATE TABLE sales_orders (
    id SERIAL PRIMARY KEY,
    customer TEXT NOT NULL,
    status TEXT NOT NULL DEFAULT 'CREATED'
        CHECK (status IN ('CREATED', 'CANCELLED')),
    total INTEGER NOT NULL CHECK (total >= 0),
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    cancelled_at TIMESTAMPTZ
);
```

### Sales Order Items
```sql
CREATE TABLE sales_order_items (
    id SERIAL PRIMARY KEY,
    order_id INTEGER NOT NULL REFERENCES sales_orders(id),
    sku TEXT NOT NULL REFERENCES products(sku),
    qty INTEGER NOT NULL CHECK (qty > 0),
    unit_price INTEGER NOT NULL CHECK (unit_price >= 0)
);
```

## API Endpoints

### Health Check

```bash
GET /api/health
```

Returns API and database health status.

### Dashboard

```bash
GET /api/dashboard
```

Returns summary statistics and recent activity.

### Products

```bash
# List products (with search, filtering, pagination)
GET /api/products?search=keramik&page=1&limit=10&status=active&low_stock=true

# Get single product
GET /api/products/{id}

# Create product
POST /api/products
{
  "sku": "CRM-100",
  "name": "Product Name",
  "price": 100000,
  "stock": 20,
  "is_active": true
}

# Update product
PUT /api/products/{id}
{
  "sku": "CRM-100",
  "name": "Updated Name",
  "price": 120000,
  "is_active": true
}

# Delete product (only if no history)
DELETE /api/products/{id}

# Adjust stock
POST /api/products/{id}/stock
{
  "type": "in",
  "qty": 10,
  "note": "Restock from supplier"
}

# Get movement history
GET /api/products/{id}/movements?page=1&limit=10
```

### Stock Movements

All stock movements are created automatically when products are created with initial stock, or when stock is adjusted through `/products/{id}/stock`.

### Sales Orders

```bash
# List orders (with search, pagination)
GET /api/orders?search=customer&page=1&limit=10

# Get single order
GET /api/orders/{id}

# Create order
POST /api/orders
{
  "customer": "PT Example Company",
  "items": [
    { "sku": "CRM-001", "qty": 3 },
    { "sku": "SAN-003", "qty": 2 }
  ]
}

# Cancel order (restores stock)
POST /api/orders/{id}/cancel
```

## API Error Handling

All errors follow a consistent JSON structure:

```json
{
  "error": {
    "code": "INSUFFICIENT_STOCK",
    "message": "Insufficient stock for the requested operation.",
    "sku": "CRM-001",
    "available": 2
  }
}
```

Common error codes:
- `VALIDATION_ERROR`: Invalid request data
- `SKU_ALREADY_EXISTS`: Duplicate SKU
- `PRODUCT_HAS_HISTORY`: Cannot delete product with history
- `INSUFFICIENT_STOCK`: Not enough stock available
- `NOT_FOUND`: Resource not found
- `ALREADY_CANCELLED`: Order already cancelled
- `INTERNAL_ERROR`: Server error

## Seed Data

The database is automatically seeded with 12 realistic construction and home supply products:

| SKU | Product | Price (IDR) | Initial Stock |
|---|---|---:|---:|
| CRM-001 | Keramik Lantai 40x40 cm | 85,000 | 120 |
| CRM-002 | Keramik Lantai 60x60 cm | 145,000 | 75 |
| CRM-003 | Keramik Dinding 25x40 cm | 95,000 | 60 |
| CRM-004 | Granit Tile 60x60 cm | 185,000 | 45 |
| CRM-005 | Granit Tile 80x80 cm | 275,000 | 30 |
| CRM-006 | Perekat Keramik 25 kg | 75,000 | 50 |
| CRM-007 | Nat Keramik 1 kg | 18,000 | 80 |
| SAN-001 | Kloset Duduk | 1,250,000 | 12 |
| SAN-002 | Wastafel Keramik | 675,000 | 15 |
| SAN-003 | Keran Air Stainless | 185,000 | 25 |
| SAN-004 | Shower Set Kamar Mandi | 450,000 | 0 |
| HOM-001 | Kitchen Sink Stainless | 550,000 | 10 |

Note: SAN-004 is inactive, and multiple products have stock ≤ 10 for demonstration.

## Testing

### Backend Tests

```bash
cd backend
go test ./...
```

Tests cover:
- Product validation and creation
- Stock adjustment (in/out)
- Stock insufficiency detection
- Sales order creation with concurrent safety
- Order cancellation and stock restoration
- Transaction integrity

### Frontend Tests

```bash
cd frontend
npm run test
```

Tests cover:
- Form validation
- Number formatting
- Relative time formatting
- Component rendering

### Build Verification

```bash
cd frontend
npm run build

cd backend
go build ./cmd/api
```

## Backend Validation and Transaction Safety

The application ensures data integrity through:

1. **Database Constraints**: Foreign keys, check constraints, unique constraints
2. **Transaction Isolation**: ACID compliance for multi-step operations
3. **Row Locking**: Pessimistic locking for stock updates to prevent race conditions
4. **Validation**: Input validation on all endpoints
5. **Error Handling**: Comprehensive error responses with actionable information

### Stock Adjustment Safety

Stock adjustments use PostgreSQL transactions with row locking:

```go
// Atomically check and update stock
UPDATE products
SET stock = stock + $1
WHERE id = $2 AND stock + $1 >= 0;
```

This prevents overselling even under concurrent requests.

### Order Creation Safety

Order creation is atomic:
1. Validate all items exist and are active
2. Lock product rows in deterministic order
3. Check stock for each item
4. Create order header
5. Create order items
6. Decrement product stock
7. Create movement records
8. Commit or rollback entire transaction

## Performance Considerations

- **Indexes**: Optimized indexes on SKU, product activity status, timestamps
- **Pagination**: All list endpoints support configurable pagination (max 100 per page)
- **Connection Pooling**: PostgreSQL connection pool (5-25 connections)
- **Lazy Loading**: Dashboard statistics computed on-demand
- **Caching**: Frontend caches product list during order creation

## Security

- **SQL Injection Prevention**: All queries use parameterized statements
- **CORS**: Configured to allow frontend origin
- **Input Validation**: All user inputs validated server-side
- **No Secrets in Code**: Database credentials via environment variables
- **No Sensitive Data Leaks**: SQL errors and stack traces not exposed to clients

## Environment Configuration

Backend configuration via environment variables (defaults provided for development):

```
DB_HOST=localhost              # PostgreSQL hostname
DB_PORT=5432                   # PostgreSQL port
DB_USER=stockflow              # PostgreSQL username
DB_PASSWORD=stockflow_dev_password  # PostgreSQL password
DB_NAME=stockflow              # PostgreSQL database name
SERVER_PORT=8080               # HTTP server port
```

See `backend/.env.example` for all options.

## Known Limitations

- No authentication or authorization (multi-user support)
- No audit logging (who/when changes)
- No role-based access control
- No customer management (only free-form customer names)
- No supplier management or purchase orders
- No inventory forecasting or analytics
- No barcode scanning
- No mobile app

## Future Enhancements

- User authentication and RBAC
- Customer and supplier management
- Purchase order workflow
- Inventory aging and expiration tracking
- Barcode/QR code scanning
- Advanced reporting and analytics
- Warehouse location tracking
- Multi-warehouse support
- API rate limiting and throttling
- Search result caching
- Real-time WebSocket updates

## Troubleshooting

### PostgreSQL Connection Failed

Ensure PostgreSQL is running and accessible:

```bash
# Check if container is running
docker ps | grep postgres

# Check logs
docker logs stockflow-postgres

# Test connection
psql -h localhost -U stockflow -d stockflow
```

### Frontend Cannot Connect to Backend

1. Ensure backend is running on port 8080
2. Check Vite proxy configuration in `vite.config.ts`
3. Verify no CORS errors in browser console
4. Check backend logs for connection issues

### Database Already Exists

The migration is idempotent. Seed data uses `ON CONFLICT DO NOTHING` to prevent duplicates. Existing data is preserved.

### Build Failures

```bash
# Clean and rebuild
cd frontend && rm -rf node_modules dist && npm install --legacy-peer-deps && npm run build

cd backend && rm main && go clean && go build ./cmd/api
```

## Contributing

To extend the application:

1. Follow existing code patterns and naming conventions
2. Add database migrations in `backend/migrations/`
3. Update API types in `frontend/src/types/`
4. Add tests for new features
5. Update README documentation

## License

This is a demonstration project for inventory management systems.

## Support

For issues or questions, refer to:
- Backend error codes and messages (API responses)
- Frontend error notifications and validation messages
- Docker Compose logs (`docker compose logs`)
- PostgreSQL logs (`docker logs stockflow-postgres`)
