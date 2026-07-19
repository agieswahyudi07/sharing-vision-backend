# Sharing Vision - Go Backend (Gin)

This is the backend service for the Sharing Vision recruitment test, built using Go (Golang) and the **Gin** HTTP web framework. It uses **GORM** to interact with MySQL and automatically manage database schemas.

## Features
- **Auto Database & Schema Migration**: At startup, it automatically creates the database if it doesn't exist and migrates the `posts` table structure.
- **Strict Validations**: Rejects request bodies that violate character length specifications:
  - `Title`: required, min 20 characters
  - `Content`: required, min 200 characters
  - `Category`: required, min 3 characters
  - `Status`: required, must choose between `publish`, `draft`, or `thrash`
- **Ambiguous Route Dispatching**: Solves Gin router wildcard path conflicts (`/article/:id` vs `/article/:limit/:offset`) by implementing a single dynamic path dispatcher.
- **CORS Configured**: Configured CORS middleware to allow requests from the Vue dev client.

## Requirements
- Go 1.22+
- MySQL Server (running at port `3306`)

## Quickstart

1. **Environment Configuration**:
   Copy `.env.example` to `.env` and fill in your MySQL credentials:
   ```bash
   cp .env.example .env
   ```
   *Make sure your MySQL server is running before executing the next step.*

2. **Run Server**:
   Download dependencies and start the hot-reload/dev server:
   ```bash
   go run main.go
   ```
   The backend server will run at `http://localhost:8080`.

---

## API Endpoint Reference

| Method | Endpoint | Request Body | Description |
| :--- | :--- | :--- | :--- |
| **POST** | `/article` | `{title, content, category, status}` | Create a new article |
| **GET** | `/article/:limit/:offset` | *None* | Paginated listing. Returns list and `X-Total-Count` header |
| **GET** | `/article/:id` | *None* | Get details of a single article |
| **PUT/PATCH** | `/article/:id` | `{title, content, category, status}` | Update an article |
| **DELETE** | `/article/:id` | *None* | Delete an article permanently |

*Note: You can also use `POST /article/:id?action=delete` or pass empty payloads to trigger deletion.*

## Postman Collection
Import the file `Postman_Collection.json` located at the root of the backend directory into Postman to test all endpoints.
