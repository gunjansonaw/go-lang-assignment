# User API - Go Backend

A RESTful API built with Go to manage users with their date of birth and dynamically calculated age.

## 🚀 Features

- ✅ CRUD operations for users
- ✅ Dynamic age calculation from date of birth
- ✅ Input validation using `go-playground/validator`
- ✅ Structured logging with Uber Zap
- ✅ Request ID middleware
- ✅ Request duration logging
- ✅ Pagination support for listing users
- ✅ Docker support
- ✅ Unit tests for age calculation

## 📋 Prerequisites

- Go 1.21 or higher
- MySQL 8.0 or higher
- Docker and Docker Compose (optional)

## 🛠️ Tech Stack

- **Framework**: [GoFiber](https://gofiber.io/)
- **Database**: MySQL
- **Code Generation**: SQLC (configured but using direct SQL queries for simplicity)
- **Logging**: [Uber Zap](https://github.com/uber-go/zap)
- **Validation**: [go-playground/validator](https://github.com/go-playground/validator)

## 📁 Project Structure

```
.
├── cmd/
│   └── server/
│       └── main.go          # Application entry point
├── config/
│   └── config.go            # Configuration management
├── db/
│   ├── migrations/          # Database migrations
│   └── queries/             # SQL queries (for SQLC)
├── internal/
│   ├── handler/             # HTTP handlers
│   ├── repository/          # Data access layer
│   ├── service/             # Business logic
│   ├── routes/              # Route definitions
│   ├── middleware/          # HTTP middleware
│   ├── models/              # Data models
│   └── logger/              # Logger setup
├── Dockerfile
├── docker-compose.yml
├── go.mod
└── README.md
```

## 🚀 Quick Start

### Using Docker Compose (Recommended)

1. Clone the repository:
```bash
git clone <repository-url>
cd gunjan
```

2. Start the services:
```bash
docker-compose up -d
```

The API will be available at `http://localhost:8080`

### Manual Setup

1. **Install dependencies:**
```bash
go mod download
```

2. **Set up MySQL database:**
```sql
CREATE DATABASE userdb;
```

3. **Run migrations:**
```bash
mysql -u root -p userdb < db/migrations/001_create_users_table.up.sql
```

4. **Set environment variables:**
```bash
export DB_HOST=localhost
export DB_PORT=3306
export DB_USER=root
export DB_PASSWORD=yourpassword
export DB_NAME=userdb
export SERVER_PORT=8080
export SERVER_HOST=0.0.0.0
```

5. **Run the application:**
```bash
go run cmd/server/main.go
```

## 📡 API Endpoints

### Create User
```http
POST /users
Content-Type: application/json

{
  "name": "Alice",
  "dob": "1990-05-10"
}
```

**Response:**
```json
{
  "id": 1,
  "name": "Alice",
  "dob": "1990-05-10"
}
```

### Get User by ID
```http
GET /users/:id
```

**Response:**
```json
{
  "id": 1,
  "name": "Alice",
  "dob": "1990-05-10",
  "age": 34
}
```

### Update User
```http
PUT /users/:id
Content-Type: application/json

{
  "name": "Alice Updated",
  "dob": "1991-03-15"
}
```

**Response:**
```json
{
  "id": 1,
  "name": "Alice Updated",
  "dob": "1991-03-15"
}
```

### Delete User
```http
DELETE /users/:id
```

**Response:** `204 No Content`

### List All Users (with Pagination)
```http
GET /users?page=1&page_size=10
```

**Response:**
```json
{
  "users": [
    {
      "id": 1,
      "name": "Alice",
      "dob": "1990-05-10",
      "age": 34
    }
  ],
  "page": 1,
  "page_size": 10,
  "total": 1,
  "total_pages": 1
}
```

### Health Check
```http
GET /health
```

## 🧪 Testing

Run unit tests:
```bash
go test ./internal/service/...
```

Run with coverage:
```bash
go test -cover ./internal/service/...
```

## 🔧 Configuration

The application uses environment variables for configuration:

| Variable | Default | Description |
|----------|---------|-------------|
| `DB_HOST` | `localhost` | Database host |
| `DB_PORT` | `3306` | Database port |
| `DB_USER` | `root` | Database user |
| `DB_PASSWORD` | `` | Database password |
| `DB_NAME` | `userdb` | Database name |
| `SERVER_PORT` | `8080` | Server port |
| `SERVER_HOST` | `0.0.0.0` | Server host |

## 📝 Notes

- Age is calculated dynamically based on the current date and the user's date of birth
- The age calculation accounts for whether the birthday has occurred this year
- Pagination defaults to page 1 with 10 items per page (max 100)
- All requests include a `X-Request-ID` header for tracing
- Request duration is logged for all endpoints

## 📄 License

This project is part of a backend development task.

