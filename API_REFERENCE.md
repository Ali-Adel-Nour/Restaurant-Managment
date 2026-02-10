# Restaurant Management API - Quick Reference

## Authentication Endpoints

| Method | Endpoint | Auth Required | Description |
|--------|----------|---------------|-------------|
| POST | `/users/signup` | ❌ | Create new user account |
| POST | `/users/login` | ❌ | Login and get JWT token |
| POST | `/users/logout` | ✅ | Logout user |

## User Endpoints

| Method | Endpoint | Auth Required | Description |
|--------|----------|---------------|-------------|
| GET | `/users` | ✅ | Get all users |
| GET | `/users/:user_id` | ❌ | Get user by ID |

## Menu Endpoints

| Method | Endpoint | Auth Required | Description |
|--------|----------|---------------|-------------|
| GET | `/menus` | ✅ | Get all menus |
| GET | `/menus/:menu_id` | ✅ | Get menu by ID |
| POST | `/menus` | ✅ | Create new menu |
| PATCH | `/menus/:menu_id` | ✅ | Update menu |

## Food Endpoints

| Method | Endpoint | Auth Required | Description |
|--------|----------|---------------|-------------|
| GET | `/foods` | ✅ | Get all foods (paginated) |
| GET | `/foods/:food_id` | ✅ | Get food by ID |
| POST | `/foods` | ✅ | Create new food item |
| PATCH | `/foods/:food_id` | ✅ | Update food item |

### Food Query Parameters
- `recordPerPage` - Number of records per page (default: 10)
- `page` - Page number (default: 1)

## Table Endpoints

| Method | Endpoint | Auth Required | Description |
|--------|----------|---------------|-------------|
| GET | `/tables` | ✅ | Get all tables |
| GET | `/tables/:table_id` | ✅ | Get table by ID |
| POST | `/tables` | ✅ | Create new table |
| PATCH | `/tables/:table_id` | ✅ | Update table |

## Order Endpoints

| Method | Endpoint | Auth Required | Description |
|--------|----------|---------------|-------------|
| GET | `/orders` | ✅ | Get all orders |
| GET | `/orders/:order_id` | ✅ | Get order by ID |
| POST | `/orders` | ✅ | Create new order |
| PATCH | `/orders/:order_id` | ✅ | Update order |

## Order Item Endpoints

| Method | Endpoint | Auth Required | Description |
|--------|----------|---------------|-------------|
| GET | `/orderItems` | ✅ | Get all order items |
| GET | `/orderItems/:orderItem_id` | ✅ | Get order item by ID |
| GET | `/orderItems/order/:order_id` | ✅ | Get all items for an order |
| POST | `/orderItems` | ✅ | Create new order item |
| PATCH | `/orderItems/:orderItem_id` | ✅ | Update order item |

## Invoice Endpoints

| Method | Endpoint | Auth Required | Description |
|--------|----------|---------------|-------------|
| GET | `/invoices` | ✅ | Get all invoices |
| GET | `/invoices/:invoice_id` | ✅ | Get invoice by ID |
| POST | `/invoices` | ✅ | Create new invoice |
| PATCH | `/invoices/:invoice_id` | ✅ | Update invoice |

## Note Endpoints

| Method | Endpoint | Auth Required | Description |
|--------|----------|---------------|-------------|
| GET | `/notes` | ✅ | Get all notes |
| GET | `/notes/:note_id` | ✅ | Get note by ID |
| POST | `/notes` | ✅ | Create new note |
| PATCH | `/notes/:note_id` | ✅ | Update note |

---

## Request Body Examples

### Sign Up
```json
{
  "first_name": "John",
  "last_name": "Doe",
  "email": "john.doe@example.com",
  "password": "password123",
  "phone": "+1234567890",
  "avatar": "https://example.com/avatar.jpg"
}
```

### Login
```json
{
  "email": "john.doe@example.com",
  "password": "password123"
}
```

### Create Menu
```json
{
  "name": "Lunch Menu",
  "category": "Main Course",
  "start_date": "2026-02-10T12:00:00Z",
  "end_date": "2026-12-31T23:59:59Z"
}
```

### Create Food
```json
{
  "name": "Margherita Pizza",
  "price": 12.99,
  "food_image": "https://example.com/pizza.jpg",
  "menu_id": "menu123"
}
```

### Create Table
```json
{
  "number_of_guests": 4,
  "table_number": 12
}
```

### Create Order
```json
{
  "order_date": "2026-02-10T18:30:00Z",
  "table_id": "table123"
}
```

### Create Order Item
```json
{
  "quantity": 2,
  "unit_price": 12.99,
  "food_id": "food123",
  "order_id": "order123"
}
```

### Create Invoice
```json
{
  "order_id": "order123",
  "payment_method": "CARD",
  "payment_status": "PENDING"
}
```

### Create Note
```json
{
  "title": "Special Instructions",
  "text": "No onions, extra cheese"
}
```

---

## Response Codes

| Code | Description |
|------|-------------|
| 200 | Success |
| 400 | Bad Request (Invalid input) |
| 401 | Unauthorized (Missing/invalid token) |
| 404 | Not Found |
| 409 | Conflict (Duplicate email/phone) |
| 500 | Internal Server Error |

---

## Authentication

All protected endpoints require a JWT token in the header:

```
Header: token
Value: <your-jwt-token>
```

Token is obtained from `/users/login` or `/users/signup` endpoints and is valid for 24 hours.

---

## Data Validation Rules

### User
- `first_name`: Required, 2-100 characters
- `last_name`: Required, 2-100 characters
- `email`: Required, valid email format
- `password`: Required, minimum 6 characters
- `phone`: Required

### Food
- `name`: Required, 2-100 characters
- `price`: Required, positive number
- `food_image`: Required, URL string
- `menu_id`: Required, valid menu ID

### Menu
- `name`: Required
- `category`: Required

### Table
- `number_of_guests`: Required
- `table_number`: Required

### Order
- `order_date`: Required, valid datetime
- `table_id`: Required, valid table ID

### Order Item
- `quantity`: Required, 1-5
- `unit_price`: Required
- `food_id`: Required, valid food ID
- `order_id`: Required, valid order ID

### Invoice
- `order_id`: Required, valid order ID
- `payment_method`: CARD | CASH | ""
- `payment_status`: Required, PENDING | PAID

---

## Tips

✅ Start by creating a Menu
✅ Add Food items to the Menu
✅ Create Tables
✅ Create Orders linked to Tables
✅ Add Order Items to Orders
✅ Generate Invoices for Orders
✅ Mark Invoices as PAID when payment received
