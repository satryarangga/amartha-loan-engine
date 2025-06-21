# Amartha Loan Management System

A Golang backend application for managing loans, borrowers, and payments with PostgreSQL database.

## Assumptions
- Borrower can only have 1 active loan, which means the loan needs to be fully repaid before he / she can make another loan
- Borrower can start to pay loan schedule 3 days before due date
- Borrower will do repayment through app or web where they can choose a payment method and click a button to pay, once its clicked the borrower can see total amount they need to pay and link to make a payment (Payment Link retrieved from payment gateway API)

## Features

- **Borrower Management**: Create and Get Detail Borrower
- **Loan Management**: Create loans with automatic schedule generation and Get loan detail
- **Payment Processing**: Generate payment links and handle payment webhooks
- **Database Migrations**: Using Goose for database schema management
- **Clean Architecture**: Controller-Service-Repository pattern
- **API Documentation**: Swagger/OpenAPI documentation

## Tech Stack
- **Language**: Go 1.24+
- **Framework**: Gin (HTTP web framework)
- **ORM**: GORM
- **Database**: PostgreSQL
- **Migrations**: Goose
- **Config Environment**: Viper
- **API Documentation**: Swagger/OpenAPI

## Prerequisites

- Go 1.24 or higher
- PostgreSQL 12 or higher
- Make (optional, for using Makefile commands)

## Installation

1. **Clone the repository**
   ```bash
   git clone https://github.com/satryarangga/amartha-loan-engine.git
   cd amartha-loan-engine
   ```

2. **Set up PostgreSQL database**
   ```bash
   # Create database
   createdb amartha
   ```

3. **Install dependencies and setup project**
   ```bash
   make setup
   ```
   This will:
   - Install all dependencies
   - Generate Swagger documentation
   - Run database migrations

4. **Configure environment variables**
   ```bash
   # Copy and edit the config file
   cp app.env.example app.env
   # Edit app.env with your database credentials
   ```

5. **Run the application**
   ```bash
   make dev
   ```

The server will start on `http://localhost:8080`

## API Documentation

### Swagger UI
Once the application is running, you can access the interactive API documentation at:
```
http://localhost:8080/swagger/index.html
```

### Generate Swagger Documentation
To regenerate the Swagger documentation after making changes:
```bash
make swagger
```

## API Endpoints

### Borrowers

- `GET /api/v1/borrowers/:id` - Get borrower by ID
- `POST /api/v1/borrowers` - Create new borrower

### Loans

- `GET /api/v1/loans` - Get all loans
- `GET /api/v1/loans/:id` - Get loan by ID

### Payments

- `POST /api/v1/payments/link` - Generate payment link
- `POST /api/v1/payment/webhook` - Handle payment webhook

## Project Structure

```
amartha/
├── cmd/
│   └── migration/
│       └── main.go
├── config/
│   ├── config.go
│   ├── database.go
│   └── logger.go
├── controllers/
│   ├── borrower_controller.go
│   ├── loan_controller.go
│   └── payment_controller.go
├── database/
│   ├── migration/
│   │   ├── main.go
│   │   └── sql/
│   │       ├── 20240101000001_create_borrowers_table.sql
│   │       ├── 20240101000002_create_loans_table.sql
│   │       ├── 20240101000003_create_loan_schedules_table.sql
│   │       └── 20240101000004_create_loan_payments_table.sql
│   └── seeder/
│       ├── main.go
│       └── sql/
│           ├── borrower.sql
│           └── loan_schedules.sql
├── docs/
│   ├── docs.go
│   ├── swagger.json
│   └── swagger.yaml
├── models/
│   ├── entity.go
│   ├── repository.go
│   ├── request.go
│   └── response.go
├── repositories/
│   ├── borrower_repository.go
│   ├── borrower_repository_impl.go
│   ├── common_repository.go
│   ├── common_repository_impl.go
│   ├── loan_payment_repository.go
│   ├── loan_payment_repository_impl.go
│   ├── loan_repository.go
│   ├── loan_repository_impl.go
│   ├── loan_schedule_repository.go
│   └── loan_schedule_repository_impl.go
├── scripts/
│   └── swagger.sh
├── services/
│   ├── borrower_service.go
│   ├── borrower_service_impl.go
│   ├── loan_service.go
│   ├── loan_service_impl.go
│   ├── payment_service.go
│   └── payment_service_impl.go
├── go.mod
├── go.sum
├── main.go
├── Makefile
└── README.md
```

## Development

### Available Make Commands
```bash
make help      # Show available commands
make deps      # Install dependencies
make build     # Build the application
make dev       # Run the application
make test      # Run tests
make clean     # Clean build artifacts
make mocks   # Generate mock for test
make swagger   # Generate Swagger documentation
make mig-up   # Run database migrations
make mig-down   # Rollback database migrations
make mig-reset   # DANGEROUS - Reset migration
make seed     # Insert seed data (for development)
make setup     # Complete project setup
```
