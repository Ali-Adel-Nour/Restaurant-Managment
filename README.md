# Restaurant Management System

A complete REST API for restaurant management built with Go, Gin framework, and MongoDB.

## Features

- 🔐 **Authentication & Authorization**: JWT-based authentication with refresh tokens
- 👥 **User Management**: User registration, login, and profile management
- 🍽️ **Menu Management**: Create and manage restaurant menus
- 🍕 **Food Management**: CRUD operations for food items with pagination
- 🪑 **Table Management**: Manage restaurant tables and seating
- 📝 **Order Management**: Complete order processing system
- 🧾 **Invoice Management**: Generate and manage invoices with payment tracking
- 📋 **Notes**: Special instructions and notes for orders

## Tech Stack

- **Language**: Go 1.24+
- **Web Framework**: Gin
- **Database**: MongoDB
- **Authentication**: JWT (JSON Web Tokens)
- **Password Hashing**: bcrypt
- **Validation**: go-playground/validator

## Prerequisites

- Go 1.24 or higher
- MongoDB 4.4 or higher
- Git

## Installation

### 1. Clone the repository

```bash
git clone https://github.com/ali-adel-nour/restaurant-management.git
cd restaurant-management
```

### 2. Install dependencies

```bash
go mod download
```

### 3. Set up environment variables

Create a `.env` file in the root directory (optional):

```env
MONGODB_URL=mongodb://localhost:27017
DB_NAME=restaurant
SECRET_KEY=your-secret-key-here
PORT=8080
```

### 4. Start MongoDB

Make sure MongoDB is running on your system:

```bash
# On Linux/Mac
sudo systemctl start mongod

# On Windows (if installed as service)
net start MongoDB
```

### 5. Run the application

```bash
go run main.go
```

The server will start on `http://localhost:8080`

## API Documentation

See [API_REFERENCE.md](API_REFERENCE.md) for complete API documentation.

## Project Structure

```
restaurant-management/
├── controllers/         # Request handlers
│   ├── collections.go
│   ├── foodController.go
│   ├── invoiceController.go
│   ├── menuController.go
│   ├── noteController.go
│   ├── orderController.go
│   ├── orderItemController.go
│   ├── tableController.go
│   └── userController.go
├── database/           # Database connection and setup
│   ├── collections.go
│   └── databaseConnection.go
├── helpers/            # Helper functions
│   └── tokenHelper.go
├── middleware/         # Middleware functions
│   └── authMiddleware.go
├── models/            # Data models
│   ├── foodModel.go
│   ├── inoviceModel.go
│   ├── menuModel.go
│   ├── noteModel.go
│   ├── orderItemModel.go
│   ├── orderModel.go
│   ├── tableModel.go
│   └── userModel.go
├── routes/            # Route definitions
│   ├── foodRouter.go
│   ├── invoiceRouter.go
│   ├── menuRouter.go
│   ├── noteRouter.go
│   ├── orderItemRouter.go
│   ├── orderRouter.go
│   ├── tableRouter.go
│   └── userRouter.go
├── .gitignore
├── go.mod
├── go.sum
├── main.go
└── README.md
```

## MongoDB Best Practices Implemented

✅ **Connection Pooling**: Configured with min/max pool sizes  
✅ **Context Timeouts**: All operations use appropriate timeouts  
✅ **Connection Management**: Proper initialization and graceful shutdown  
✅ **Environment Configuration**: Flexible connection string via env vars  
✅ **Health Checks**: Database ping on startup  
✅ **Idle Connection Handling**: Auto-close idle connections  
✅ **Centralized Collections**: Single source of truth for all collections  

## API Endpoints

### Authentication (Public)
- `POST /users/signup` - Register new user
- `POST /users/login` - User login
- `GET /users/:user_id` - Get user by ID

### Users (Protected)
- `GET /users` - Get all users
- `POST /users/logout` - User logout

### Menus (Protected)
- `GET /menus` - Get all menus
- `GET /menus/:menu_id` - Get menu by ID
- `POST /menus` - Create menu
- `PATCH /menus/:menu_id` - Update menu

### Foods (Protected)
- `GET /foods` - Get all foods (paginated)
- `GET /foods/:food_id` - Get food by ID
- `POST /foods` - Create food item
- `PATCH /foods/:food_id` - Update food item

### Tables (Protected)
- `GET /tables` - Get all tables
- `GET /tables/:table_id` - Get table by ID
- `POST /tables` - Create table
- `PATCH /tables/:table_id` - Update table

### Orders (Protected)
- `GET /orders` - Get all orders
- `GET /orders/:order_id` - Get order by ID
- `POST /orders` - Create order
- `PATCH /orders/:order_id` - Update order

### Order Items (Protected)
- `GET /orderItems` - Get all order items
- `GET /orderItems/:orderItem_id` - Get order item by ID
- `GET /orderItems/order/:order_id` - Get items by order
- `POST /orderItems` - Create order item
- `PATCH /orderItems/:orderItem_id` - Update order item

### Invoices (Protected)
- `GET /invoices` - Get all invoices
- `GET /invoices/:invoice_id` - Get invoice by ID
- `POST /invoices` - Create invoice
- `PATCH /invoices/:invoice_id` - Update invoice

### Notes (Protected)
- `GET /notes` - Get all notes
- `GET /notes/:note_id` - Get note by ID
- `POST /notes` - Create note
- `PATCH /notes/:note_id` - Update note

## Authentication

Protected endpoints require a JWT token in the header:

```
Header: token
Value: <your-jwt-token>
```

Tokens are valid for 24 hours and can be obtained from login/signup endpoints.

## Development

### Build

```bash
go build -o restaurant-server
```

### Run

```bash
./restaurant-server
```

### Test

```bash
go test ./...
```

## Security Features

- Password hashing using bcrypt (cost factor 14)
- JWT-based authentication
- Token refresh mechanism
- Request validation
- Secure password requirements (minimum 6 characters)

## Contributing

1. Fork the repository
2. Create a feature branch (`git checkout -b feature/amazing-feature`)
3. Commit your changes (`git commit -m 'Add some amazing feature'`)
4. Push to the branch (`git push origin feature/amazing-feature`)
5. Open a Pull Request

## License

This project is licensed under the MIT License.

## Author

Ali Adel Nour

## Acknowledgments

- Gin Web Framework
- MongoDB Go Driver
- JWT-Go library
- Go Community
