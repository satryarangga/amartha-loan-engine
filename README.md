# Amartha Loan Management System

A Golang backend application for managing loans, borrowers, and payments with PostgreSQL database.

## Features

- **Borrower Management**: Create, read, update, and delete borrowers
- **Loan Management**: Create loans with automatic schedule generation
- **Payment Processing**: Generate payment links and handle payment webhooks
- **Database Migrations**: Using Goose for database schema management
- **Clean Architecture**: Controller-Service-Repository pattern

## Tech Stack

- **Language**: Go 1.24+
- **Framework**: Gin (HTTP web framework)
- **ORM**: GORM
- **Database**: PostgreSQL
- **Migrations**: Goose
- **Environment**: godotenv

## Prerequisites

- Go 1.24 or higher
- PostgreSQL 12 or higher
- Make (optional, for using Makefile commands)

## Installation

1. **Clone the repository**
   ```bash
   git clone <repository-url>
   cd amartha
   ```

2. **Install dependencies**
   ```bash
   go mod tidy
   ```

3. **Set up PostgreSQL database**
   ```bash
   # Create database
   createdb amartha_db
   ```

4. **Configure environment variables**
   ```bash
   # Copy and edit the config file
   cp config.env .env
   # Edit .env with your database credentials
   ```

5. **Run database migrations**
   ```bash
   # Install goose if not already installed
   go install github.com/pressly/goose/v3/cmd/goose@latest
   
   # Run migrations
   goose -dir migrations postgres "host=localhost user=postgres password=password dbname=amartha_db sslmode=disable" up
   ```

6. **Run the application**
   ```bash
   go run main.go
   ```

The server will start on `http://localhost:8080`

## API Endpoints

### Borrowers

- `GET /api/v1/borrowers` - Get all borrowers
- `GET /api/v1/borrowers/:id` - Get borrower by ID
- `POST /api/v1/borrowers` - Create new borrower
- `PUT /api/v1/borrowers/:id` - Update borrower
- `DELETE /api/v1/borrowers/:id` - Delete borrower

### Loans

- `POST /api/v1/loans` - Create new loan
- `GET /api/v1/loans` - Get all loans
- `GET /api/v1/loans/:id` - Get loan by ID
- `PUT /api/v1/loans/:id` - Update loan
- `DELETE /api/v1/loans/:id` - Delete loan

### Payments

- `POST /api/v1/loans/:id/payment-link` - Generate payment link
- `POST /api/v1/webhook/payment` - Handle payment webhook
- `GET /api/v1/loans/:id/payments` - Get payment history

## API Examples

### Create a Borrower

```bash
curl -X POST http://localhost:8080/api/v1/borrowers \
  -H "Content-Type: application/json" \
  -d '{
    "first_name": "John",
    "last_name": "Doe",
    "phone_number": "+1234567890"
  }'
```

### Create a Loan

```bash
curl -X POST http://localhost:8080/api/v1/loans \
  -H "Content-Type: application/json" \
  -d '{
    "borrower_id": "borrower-uuid-here",
    "amount": 10000.00,
    "repayment_cadence_days": 30,
    "interest_percentage": 5.5
  }'
```

### Generate Payment Link

```bash
curl -X POST http://localhost:8080/api/v1/loans/loan-uuid-here/payment-link
```

### Payment Webhook

```bash
curl -X POST http://localhost:8080/api/v1/webhook/payment \
  -H "Content-Type: application/json" \
  -d '{
    "loan_id": "loan-uuid-here",
    "amount": 500.00,
    "payment_method": "bank_transfer"
  }'
```

## Database Schema

### Borrowers Table
- `id` (UUID, Primary Key)
- `first_name` (VARCHAR)
- `last_name` (VARCHAR)
- `phone_number` (VARCHAR, Unique)
- `is_delinquent` (BOOLEAN)
- `created_at` (TIMESTAMP)
- `updated_at` (TIMESTAMP)

### Loans Table
- `id` (UUID, Primary Key)
- `borrower_id` (UUID, Foreign Key)
- `amount` (DECIMAL)
- `repayment_cadence_days` (INTEGER)
- `interest_percentage` (DECIMAL)
- `interest_amount` (DECIMAL)
- `status` (VARCHAR)
- `created_at` (TIMESTAMP)
- `updated_at` (TIMESTAMP)

### Loan Schedules Table
- `id` (UUID, Primary Key)
- `loan_id` (UUID, Foreign Key)
- `due_date` (TIMESTAMP)
- `basic_amount` (DECIMAL)
- `interest_amount` (DECIMAL)
- `total_payment` (DECIMAL)
- `status` (VARCHAR)
- `created_at` (TIMESTAMP)
- `updated_at` (TIMESTAMP)

### Loan Payments Table
- `id` (UUID, Primary Key)
- `loan_id` (UUID, Foreign Key)
- `loan_schedule_ids` (UUID[])
- `total_payment` (DECIMAL)
- `payment_method` (VARCHAR)
- `created_at` (TIMESTAMP)
- `updated_at` (TIMESTAMP)

## Project Structure

```
amartha/
├── config/
│   └── database.go
├── controllers/
│   ├── borrower_controller.go
│   ├── loan_controller.go
│   └── payment_controller.go
├── models/
│   └── models.go
├── repositories/
│   ├── borrower_repository.go
│   ├── loan_repository.go
│   ├── loan_schedule_repository.go
│   └── loan_payment_repository.go
├── services/
│   ├── borrower_service.go
│   ├── loan_service.go
│   └── payment_service.go
├── migrations/
│   ├── 20240101000001_create_borrowers_table.sql
│   ├── 20240101000002_create_loans_table.sql
│   ├── 20240101000003_create_loan_schedules_table.sql
│   └── 20240101000004_create_loan_payments_table.sql
├── main.go
├── go.mod
├── config.env
└── README.md
```

## Development

### Running Tests
```bash
go test ./...
```

### Database Migrations
```bash
# Create new migration
goose -dir migrations create migration_name

# Run migrations
goose -dir migrations postgres "connection-string" up

# Rollback migration
goose -dir migrations postgres "connection-string" down
```

### Code Formatting
```bash
go fmt ./...
```

### Linting
```bash
go vet ./...
```

## Environment Variables

Create a `.env` file with the following variables:

```env
DB_HOST=localhost
DB_PORT=5432
DB_USER=postgres
DB_PASSWORD=password
DB_NAME=amartha_db
DB_SSLMODE=disable
SERVER_PORT=8080
```

## Contributing

1. Fork the repository
2. Create a feature branch
3. Make your changes
4. Add tests if applicable
5. Submit a pull request

## License

This project is licensed under the MIT License. 