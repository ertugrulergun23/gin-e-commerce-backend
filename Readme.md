# 🛒 Go Gin E-Commerce RESTful Backend API

Modern, secure, and production-ready E-Commerce Backend API built with **Go (Golang)**, **Gin Web Framework**, **GORM ORM**, **PostgreSQL**, and **Caddy Reverse Proxy**.

---

## 🚀 Tech Stack & Features

- **Language & Framework:** Go 1.26+ & [Gin Web Framework](https://github.com/gin-gonic/gin)
- **Database & ORM:** PostgreSQL 16 & [GORM](https://gorm.io/)
- **Database Migrations:** [golang-migrate](https://github.com/golang-migrate/migrate) with Go `embed.FS` (automatically runs migrations on startup)
- **Authentication & Authorization:** JWT (JSON Web Tokens) with Role-Based Access Control (**RBAC**: `admin`, `seller`, `buyer`)
- **Reverse Proxy & SSL/TLS:** [Caddy Server](https://caddyserver.com/) with automatic HTTPS
- **Containerization & Live Reload:** Docker, Docker Compose & [Air](https://github.com/air-verse/air)
- **Security Middlewares:**
  - 🛡️ **Security Headers:** HSTS, CSP, X-Frame-Options, X-XSS-Protection, X-Content-Type-Options, Referrer-Policy, Permissions-Policy
  - 🌐 **CORS:** Cross-Origin Resource Sharing configuration with automatic preflight handling
  - ⏱️ **Rate Limiting:** In-memory token bucket rate limiter per client IP (100 req/min)
  - 📦 **Request Body Limiter:** 1MB maximum payload enforcement (`413 Payload Too Large`)
  - 🔍 **Content-Type Validation:** Enforces `application/json` on state-changing requests (POST/PUT/PATCH)
- **Unit Testing:** Comprehensive test suite with [go-sqlmock](https://github.com/DATA-DOG/go-sqlmock) covering all routes, handlers, and middlewares.

---

## 📁 Project Structure

```text
├── auth/                 # JWT token generation, parsing, and claims
├── database/             # Database initialization and migration runner
│   └── migrations/       # SQL migration scripts (.up.sql & .down.sql)
├── handlers/             # HTTP controller handlers (Cart, Comment, Login, Order, Product, User)
├── middleware/           # Auth and Security middlewares (RateLimiter, CORS, Headers, BodySize)
├── models/               # GORM database models & validation bindings
├── routers/              # Gin engine setup & endpoint routing
├── tests/                # Unit & integration tests with SQLMock
├── Caddyfile             # Caddy reverse proxy configuration
├── Dockerfile            # Multi-stage/Air container definition
├── docker-compose.yml    # Full stack definition (App + Postgres + Caddy)
├── .env.example          # Environment variables template
├── main.go               # Application entry point
└── README.md
```

---

## ⚙️ Getting Started

### Prerequisites

- [Docker](https://www.docker.com/) & [Docker Compose](https://docs.docker.com/compose/)
- *(Optional)* [Go 1.26+](https://go.dev/dl/) if running locally without Docker.

### 1. Clone the Repository

```bash
git clone https://github.com/your-username/gin-e-commerce-backend.git
cd gin-e-commerce-backend
```

### 2. Configure Environment Variables

Copy the example environment file:

```bash
cp .env.example .env
```

Edit `.env` to configure your `JWT_SECRET_KEY` and optional database credentials.

### 3. Run with Docker Compose

Start the entire stack (PostgreSQL + App with Air live-reload + Caddy HTTPS):

```bash
docker-compose up -d --build
```

The application will be accessible at:
- **HTTPS (via Caddy):** `https://localhost`
- **HTTP (via Caddy):** `http://localhost`

---

## 🔌 API Endpoints

### 🔓 Public Endpoints

| Method | Endpoint | Description |
|---|---|---|
| `GET` | `/ping` | Health check endpoint |
| `POST` | `/login` | Authenticate user and receive JWT Token |
| `POST` | `/user/create` | Register a new user (`buyer`, `seller`, `admin`) |
| `GET` | `/products` | List products with query filters (`name`, `down_price`, `up_price`, `point`) |
| `GET` | `/product/:id` | Get single product details by ID |
| `GET` | `/comments` | List comments with query filters (`owner`, `product`) |
| `GET` | `/comment/:id` | Get single comment by ID |

---

### 🔒 Protected Endpoints (`Authorization: Bearer <token>`)

#### 👤 User Management
| Method | Endpoint | Required Role | Description |
|---|---|---|---|
| `GET` | `/api/user` | Authenticated User | Get current profile (password masked) |
| `PUT` | `/api/user/update` | Authenticated User | Update profile (username, email, password) |
| `DELETE` | `/api/user/delete` | Authenticated User | Delete own account |

#### 📦 Product Management
| Method | Endpoint | Required Role | Description |
|---|---|---|---|
| `POST` | `/api/product/create` | `seller`, `admin` | Create a new product |
| `PUT` | `/api/product/:id/update` | Seller (Owner) | Update own product details |
| `DELETE` | `/api/product/:id/delete` | Seller (Owner) | Delete own product |

#### 💬 Comments & Reviews
| Method | Endpoint | Required Role | Description |
|---|---|---|---|
| `POST` | `/api/comment/create/:product_id` | Authenticated User | Write a review/comment with rating (1-5) |
| `PUT` | `/api/comment/:id/update` | Comment Author | Update own comment |
| `DELETE` | `/api/comment/:id/delete` | Comment Author | Delete own comment |

#### 🛒 Shopping Cart
| Method | Endpoint | Required Role | Description |
|---|---|---|---|
| `GET` | `/api/cart` | Authenticated User | View items in current user's cart |
| `POST` | `/api/cart/add` | Authenticated User | Add product and quantity to cart |
| `PUT` | `/api/cart/:id/update` | Cart Owner | Update item quantity (deletes if quantity becomes 0) |
| `DELETE` | `/api/cart/:id/delete` | Cart Owner | Remove item from cart |

#### 📋 Orders & Checkout
| Method | Endpoint | Required Role | Description |
|---|---|---|---|
| `GET` | `/api/order` | Authenticated User | Get user's orders and order items |
| `POST` | `/api/order/cart` | Authenticated User | Checkout cart into a new order (Atomic DB Transaction) |
| `POST` | `/api/order/create` | Authenticated User | Create a direct single-item order |
| `PUT` | `/api/order/:id/update` | `admin` | Update order status (`pending`, `shipped`, `delivered`, etc.) |

---

## 🧪 Running Tests

To run the complete automated test suite:

```bash
go test -v ./...
```

To run without cache:

```bash
go test -v -count=1 ./...
```

To view test coverage:

```bash
go test -cover ./...
```

---

## 🔒 Security Implementations

1. **Password Hashing:** Passwords are never stored in plaintext; hashed using `bcrypt` (DefaultCost).
2. **JWT Authentication:** Cryptographically signed tokens with HMAC-SHA256 and 24-hour expiration.
3. **RBAC:** Strict role checks preventing buyers from creating products or non-admins from modifying order statuses.
4. **Data Isolation:** User ID is extracted directly from the verified JWT claims to prevent IDOR vulnerabilities.
5. **Database Transactions:** Multi-step order creation and cart clearing run inside atomic GORM DB transactions.
6. **SQL Injection Prevention:** Parameterized SQL queries across all GORM statements.
7. **Rate Limiting & DoS Protection:** Built-in IP rate limiter and 1MB request body payload limit.
