# StockFlow

**Inventory management for construction materials and home supplies**
**Aplikasi manajemen inventaris untuk bahan bangunan dan perlengkapan rumah**

StockFlow is a full-stack demo application for managing products, stock movements, and sales orders. It is built with Go, React, TypeScript, PostgreSQL, and Docker Compose.

StockFlow adalah aplikasi demo full-stack untuk mengelola produk, pergerakan stok, dan pesanan penjualan. Aplikasi ini dibangun menggunakan Go, React, TypeScript, PostgreSQL, dan Docker Compose.

## Features | Fitur

- Product listing, search, filters, pagination, and product maintenance.
  Daftar produk, pencarian, filter, pagination, serta pengelolaan produk.
- Stock-in and stock-out adjustments with a movement history.
  Penyesuaian stok masuk dan keluar beserta riwayat pergerakannya.
- Sales orders that deduct stock and can be cancelled to restore it.
  Pesanan penjualan yang mengurangi stok dan dapat dibatalkan untuk mengembalikan stok.
- Dashboard with stock and order summaries.
  Dashboard berisi ringkasan stok dan pesanan.
- PostgreSQL initialization with schema and sample data on first startup.
  Inisialisasi PostgreSQL dengan skema dan data contoh saat pertama kali dijalankan.

## Technology Stack | Teknologi

| Area | Technology |
|---|---|
| Backend | Go 1.27+, standard-library `net/http`, pgx/v5 |
| Frontend | React 19, TypeScript, Vite |
| Database | PostgreSQL 17 |
| Local orchestration | Docker Compose |

## Project Structure | Struktur Proyek

```text
stockflow/
â”œâ”€â”€ backend/
â”‚   â”œâ”€â”€ cmd/api/              # API entry point
â”‚   â”œâ”€â”€ internal/             # Config, database, handlers, models, repositories, services
â”‚   â”œâ”€â”€ migrations/           # SQL schema and seed scripts
â”‚   â”œâ”€â”€ Dockerfile
â”‚   â””â”€â”€ go.mod
â”œâ”€â”€ frontend/
â”‚   â”œâ”€â”€ src/                  # React application
â”‚   â”œâ”€â”€ Dockerfile
â”‚   â””â”€â”€ package.json
â”œâ”€â”€ docker-compose.yml
â””â”€â”€ README.md
```

## Requirements | Prasyarat

### Docker Compose (recommended | disarankan)

- Docker Desktop with Docker Compose v2.
- Ports `5173`, `8080`, and `5434` available on the host.

### Running services locally | Menjalankan service secara lokal

- Go 1.27 or newer.
- Node.js 20 or newer and npm.
- Docker, for running the PostgreSQL service.

## Run with Docker Compose | Menjalankan dengan Docker Compose

From the project root, build and start the application:

Dari root proyek, build dan jalankan aplikasi:

```powershell
cd C:\stockflow
docker compose up --build
```

Wait until PostgreSQL reports healthy and the backend starts. Open:

Tunggu sampai PostgreSQL berstatus healthy dan backend berjalan. Buka:

- Frontend: <http://localhost:5173>
- Backend health check: <http://localhost:8080/api/health>
- PostgreSQL from the host: `localhost:5434` (inside Compose, the port is `5432`)

The Compose setup mounts `backend/migrations` into PostgreSQL's initialization directory. PostgreSQL runs `001_init.sql` and then `002_seed.sql` **only when the database volume is initialized for the first time**. The named volume `stockflow_data` preserves data across restarts.

Konfigurasi Compose memasang `backend/migrations` ke direktori inisialisasi PostgreSQL. PostgreSQL menjalankan `001_init.sql` lalu `002_seed.sql` **hanya ketika volume database pertama kali dibuat**. Named volume `stockflow_data` menyimpan data saat container dimulai ulang.

To stop the services, press `Ctrl+C`, or run `docker compose down` in another terminal. `docker compose down` keeps the database volume. Removing the volume also deletes its database data; only do so when you intend to reset the local database.

Untuk menghentikan service, tekan `Ctrl+C` atau jalankan `docker compose down` dari terminal lain. Perintah tersebut tidak menghapus volume database. Menghapus volume juga akan menghapus data database lokal; lakukan hanya jika memang ingin mereset database.

## Run Locally | Menjalankan Secara Lokal

Use this option when you want hot reload for the frontend and run the Go API directly on the host. PostgreSQL still runs in Docker.

Gunakan opsi ini jika ingin memakai hot reload pada frontend dan menjalankan API Go langsung di host. PostgreSQL tetap dijalankan melalui Docker.

### 1. Start PostgreSQL | Jalankan PostgreSQL

From the project root:

Dari root proyek:

```powershell
docker compose up -d postgres
```

Compose exposes PostgreSQL on host port `5434` and initializes schema and seed data on a fresh volume.

Compose membuka PostgreSQL pada port host `5434` serta menginisialisasi skema dan data contoh jika volume masih baru.

### 2. Start the backend | Jalankan backend

In a PowerShell terminal, configure the database port to match the Compose host mapping, then start the API:

Di terminal PowerShell, atur port database agar sesuai dengan pemetaan port Compose, lalu jalankan API:

```powershell
$env:DB_PORT = "5434"
cd C:\stockflow\backend
go run ./cmd/api
```

The API listens on <http://localhost:8080>. The backend reads `DB_HOST`, `DB_PORT`, `DB_USER`, `DB_PASSWORD`, `DB_NAME`, and `SERVER_PORT`; development defaults are documented in `backend/.env.example`. The application does not automatically load a `.env` file, so set environment variables in the shell before starting it.

API berjalan di <http://localhost:8080>. Backend membaca `DB_HOST`, `DB_PORT`, `DB_USER`, `DB_PASSWORD`, `DB_NAME`, dan `SERVER_PORT`; nilai default development tersedia di `backend/.env.example`. Aplikasi tidak memuat file `.env` secara otomatis, jadi atur environment variable di shell sebelum menjalankannya.

### 3. Start the frontend | Jalankan frontend

In a second terminal:

Di terminal kedua:

```powershell
cd C:\stockflow\frontend
npm install --legacy-peer-deps
npm run dev
```

Open <http://localhost:5173>. Vite proxies `/api` requests to `http://localhost:8080`.

Buka <http://localhost:5173>. Vite meneruskan request `/api` ke `http://localhost:8080`.

## API Reference and cURL Examples | Referensi API dan Contoh cURL

The examples below use `curl.exe`, which avoids PowerShell's `curl` alias. Run them in PowerShell after starting the backend. Request and response bodies use JSON; prices are integer Indonesian rupiah (IDR).

Contoh berikut menggunakan `curl.exe` agar tidak tertukar dengan alias `curl` di PowerShell. Jalankan setelah backend aktif. Request dan response menggunakan JSON; harga berupa bilangan bulat dalam rupiah (IDR).

### Health and dashboard | Health check dan dashboard

```powershell
curl.exe http://localhost:8080/api/health
curl.exe http://localhost:8080/api/dashboard
```

### List and search products | Daftar dan pencarian produk

```powershell
curl.exe "http://localhost:8080/api/products?search=keramik&page=1&limit=10&status=active"
curl.exe "http://localhost:8080/api/products?low_stock=true"
curl.exe http://localhost:8080/api/products/1
```

Product list query parameters include `search`, `status`, `low_stock`, `page`, and `limit`. Page size is capped at 100.

Parameter query untuk daftar produk mencakup `search`, `status`, `low_stock`, `page`, dan `limit`. Ukuran halaman dibatasi maksimal 100.

### Create a product | Membuat produk

```powershell
curl.exe -X POST http://localhost:8080/api/products `
  -H "Content-Type: application/json" `
  -d '{"sku":"DEMO-001","name":"Demo Tile","price":85000,"stock":20,"is_active":true}'
```

### Update a product | Memperbarui produk

`PUT` updates product details; stock is adjusted through the stock endpoint instead.

`PUT` memperbarui detail produk; perubahan stok dilakukan melalui endpoint stok.

```powershell
curl.exe -X PUT http://localhost:8080/api/products/1 `
  -H "Content-Type: application/json" `
  -d '{"sku":"CRM-001","name":"Keramik Lantai 40x40 cm","price":90000,"is_active":true}'
```

### Adjust stock and view movement history | Menyesuaikan stok dan melihat riwayat

Use `type: "in"` to add stock or `type: "out"` to remove stock. Quantity must be positive.

Gunakan `type: "in"` untuk menambah stok atau `type: "out"` untuk mengurangi stok. Jumlah harus positif.

```powershell
curl.exe -X POST http://localhost:8080/api/products/1/stock `
  -H "Content-Type: application/json" `
  -d '{"type":"in","qty":5,"note":"Restock from supplier"}'

curl.exe "http://localhost:8080/api/products/1/movements?page=1&limit=10"
curl.exe "http://localhost:8080/api/movements?type=in&page=1&limit=20"
```

### Create, list, and cancel a sales order | Membuat, melihat, dan membatalkan pesanan

The order endpoint calculates line totals from the current product prices. Products must be active and have sufficient stock.

Endpoint pesanan menghitung subtotal berdasarkan harga produk saat ini. Produk harus aktif dan stok harus mencukupi.

```powershell
curl.exe -X POST http://localhost:8080/api/orders `
  -H "Content-Type: application/json" `
  -d '{"customer":"PT Example","items":[{"sku":"CRM-001","qty":2},{"sku":"SAN-003","qty":1}]}'

curl.exe "http://localhost:8080/api/orders?search=Example&page=1&limit=10"
curl.exe http://localhost:8080/api/orders/1
curl.exe -X POST http://localhost:8080/api/orders/1/cancel
```

Use the order `id` returned by the create/list response in place of `1` when retrieving or cancelling a specific order. Cancellation restores the ordered quantities to stock and records stock movements.

Gunakan `id` pesanan dari response create/list sebagai pengganti `1` saat melihat atau membatalkan pesanan. Pembatalan mengembalikan jumlah barang ke stok dan mencatat pergerakan stok.

### Endpoint summary | Ringkasan endpoint

| Method | Path | Purpose |
|---|---|---|
| `GET` | `/api/health` | Health check, including database connectivity |
| `GET` | `/api/dashboard` | Dashboard statistics and recent activity |
| `GET` | `/api/products` | List and filter products |
| `GET` | `/api/products/{id}` | Get a product |
| `POST` | `/api/products` | Create a product |
| `PUT` | `/api/products/{id}` | Update product details |
| `DELETE` | `/api/products/{id}` | Delete a product when allowed |
| `POST` | `/api/products/{id}/stock` | Adjust product stock |
| `GET` | `/api/products/{id}/movements` | List movements for a product |
| `GET` | `/api/movements` | List and filter all movements |
| `GET` | `/api/orders` | List and filter sales orders |
| `GET` | `/api/orders/{id}` | Get a sales order |
| `POST` | `/api/orders` | Create a sales order |
| `POST` | `/api/orders/{id}/cancel` | Cancel an order and restore stock |

## Design Decisions | Keputusan Desain

- **Go standard-library HTTP routing.** The API uses `net/http`'s `ServeMux`, keeping routing dependencies small and relying on Go's built-in method/path patterns.
  **Routing HTTP standard library Go.** API menggunakan `ServeMux` dari `net/http`, sehingga dependensi routing tetap sederhana dan memanfaatkan pola method/path bawaan Go.
- **PostgreSQL for relational integrity.** Products, stock movements, order headers, and order items use relational tables, foreign keys, and check constraints.
  **PostgreSQL untuk integritas relasional.** Produk, pergerakan stok, header pesanan, dan item pesanan menggunakan tabel relasional, foreign key, serta check constraint.
- **Transactions for stock-changing workflows.** Stock adjustment, order creation, and cancellation group related database writes into transactions. Stock updates also use row locking in the repository.
  **Transaksi untuk proses yang mengubah stok.** Penyesuaian stok, pembuatan pesanan, dan pembatalan mengelompokkan penulisan database terkait dalam transaksi. Update stok juga menggunakan row lock di repository.
- **Integer IDR values.** Prices are stored as integers to avoid floating-point rounding for the project's whole-rupiah currency values.
  **Nilai IDR berupa integer.** Harga disimpan sebagai bilangan bulat untuk menghindari pembulatan floating-point pada nominal rupiah utuh di proyek ini.
- **Compose-based local setup.** Docker Compose provides a repeatable local database and app environment; PostgreSQL's initialization scripts seed only a new database volume.
  **Setup lokal berbasis Compose.** Docker Compose menyediakan lingkungan aplikasi dan database lokal yang konsisten; skrip inisialisasi PostgreSQL hanya melakukan seed pada volume database baru.

## Testing and Build | Pengujian dan Build

Run backend tests:

Jalankan test backend:

```powershell
cd C:\stockflow\backend
go test ./...
```

Run frontend tests once, then build:

Jalankan test frontend sekali, lalu build:

```powershell
cd C:\stockflow\frontend
npm run test -- --run
npm run build
```

The frontend test command uses Vitest; `--run` makes it exit after the test run instead of staying in watch mode.

Perintah test frontend menggunakan Vitest; `--run` membuat proses selesai setelah pengujian, bukan tetap berjalan dalam mode watch.

## Known Limitations and What I Would Improve | Batasan dan Rencana Perbaikan

This is a demo application, not a production-ready inventory platform. It currently has no authentication or authorization, and it has no customer/supplier management, purchase-order workflow, or multi-warehouse support.

Ini adalah aplikasi demo, bukan platform inventaris siap produksi. Saat ini belum ada autentikasi maupun otorisasi, serta belum tersedia pengelolaan pelanggan/pemasok, alur purchase order, atau dukungan multi-gudang.

If I continued development, I would prioritize:

Jika pengembangan dilanjutkan, prioritas saya adalah:

1. **Authentication, authorization, and audit history** so access and important stock changes can be attributed to users.
   **Autentikasi, otorisasi, dan audit history** agar akses serta perubahan stok penting dapat ditelusuri ke pengguna.
2. **Versioned migrations** with a migration runner, rather than relying on PostgreSQL's first-start initialization directory for schema changes.
   **Migration berversi** dengan migration runner, bukan mengandalkan direktori inisialisasi PostgreSQL saat pertama kali berjalan untuk perubahan skema.
3. **More integration and end-to-end tests**, especially for database concurrency, rollback behavior, and the main UI workflows.
   **Lebih banyak integration dan end-to-end test**, terutama untuk konkurensi database, perilaku rollback, dan alur utama UI.
4. **Operational hardening**, including production secret management, stricter configuration validation, structured logging, and health/readiness checks.
   **Penguatan operasional**, termasuk pengelolaan secret untuk production, validasi konfigurasi yang lebih ketat, structured logging, serta health/readiness check.
5. **Business workflows** for suppliers, purchase orders, and warehouse locations once the core inventory requirements are clear.
   **Alur bisnis** untuk pemasok, purchase order, dan lokasi gudang setelah kebutuhan inti inventaris ditetapkan.
