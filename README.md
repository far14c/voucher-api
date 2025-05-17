# Voucher API

A RESTful API service for managing vouchers, customer points, and redemptions built with Go.

## Features

- Brand management
- Voucher creation and management
- Customer points tracking
- Voucher redemption system
- MySQL database integration

## Prerequisites

- Go 1.21 or higher
- MySQL 5.7 or higher (XAMPP or standalone)
- Make sure you have the following environment variables set up in your `.env` file:
  ```
  DB_HOST=localhost
  DB_PORT=3306
  DB_USER=root
  DB_PASSWORD=
  DB_NAME=voucher_db
  ```

## Installation

1. Clone the repository
   ```bash
   git clone https://github.com/yourusername/voucher-api.git
   cd voucher-api
   ```

2. Install dependencies
   ```bash
   go mod download
   ```

3. Set up the database
   - Create a new MySQL database named `voucher_db`
   - Import the schema using `schema.sql`

4. Create and configure `.env` file
   ```bash
   cp .env.example .env
   # Edit .env with your database credentials
   ```

5. Run the application
   ```bash
   go run main.go
   ```

## API Endpoints

### Brands
- `GET /api/brands` - List all brands
- `POST /api/brands` - Create a new brand
- `GET /api/brands/{id}` - Get brand details
- `PUT /api/brands/{id}` - Update brand
- `DELETE /api/brands/{id}` - Delete brand

### Vouchers
- `GET /api/vouchers` - List all vouchers
- `POST /api/vouchers` - Create a new voucher
- `GET /api/vouchers/{id}` - Get voucher details
- `PUT /api/vouchers/{id}` - Update voucher
- `DELETE /api/vouchers/{id}` - Delete voucher
- `GET /api/vouchers/brand/{id}` - Get vouchers by brand

### Customers
- `GET /api/customers` - List all customers
- `POST /api/customers` - Create a new customer
- `GET /api/customers/{id}` - Get customer details
- `PUT /api/customers/{id}` - Update customer
- `DELETE /api/customers/{id}` - Delete customer

### Redemptions
- `POST /api/redemptions` - Create a new redemption
- `GET /api/redemptions/{id}` - Get redemption details
- `GET /api/customers/{id}/redemptions` - Get customer's redemptions

## Deployment

### Free Hosting Options

#### Database Hosting (db4free.net)
1. Go to [db4free.net](https://db4free.net)
2. Create a new account and database
3. Note down your database credentials
4. Update `config.production.yaml` with your db4free.net credentials

#### Application Hosting (Render.com)
1. Create an account on [Render](https://render.com)
2. Create a new Web Service
3. Connect your GitHub repository
4. Configure the build settings:
   - Build Command: `go build -o main .`
   - Start Command: `./main`
5. Add environment variables from your `config.production.yaml`

### Docker Deployment
1. Build the Docker image:
   ```bash
   docker build -t voucher-api .
   ```

2. Run the container:
   ```bash
   docker run -p 8080:8080 voucher-api
   ```

## Database Schema

The application uses the following tables:
- `brands` - Store brand information
- `vouchers` - Store voucher details
- `customers` - Store customer information and points balance
- `redemptions` - Store redemption transactions
- `redemption_items` - Store individual items in a redemption

## Contributing

1. Fork the repository
2. Create your feature branch (`git checkout -b feature/amazing-feature`)
3. Commit your changes (`git commit -m 'Add some amazing feature'`)
4. Push to the branch (`git push origin feature/amazing-feature`)
5. Open a Pull Request

## License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.
