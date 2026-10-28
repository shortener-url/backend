# URL Shortener - Backend

Welcome to the **URL Shortener Backend**, a robust and efficient API service built with **Go (Golang)**. This project handles the core logic for shortening URLs, managing redirections, and providing analytics data, following a clean **Hexagonal Architecture**.

## 🚀 Features

- **Efficient URL Shortening**: Generates unique, compact keys for any long URL.
- **Fast Redirection**: High-performance HTTP redirection to original links.
- **Collision Handling**: Built-in logic to handle key collisions during generation.
- **Swagger Documentation**: Interactive API documentation for easy integration.
- **Hexagonal Architecture**: Clean separation of concerns for better maintainability and testing.
- **Docker Support**: Ready for containerized environments.

## 🛠️ Tech Stack

- **Language**: [Go 1.22+](https://go.dev/)
- **Web Framework**: [Gin Gonic](https://gin-gonic.com/)
- **ORM**: [GORM](https://gorm.io/)
- **Database**: [PostgreSQL](https://www.postgresql.org/)
- **Documentation**: [Swagger (swag)](https://github.com/swaggo/swag)
- **Hot Reload**: [Air](https://github.com/air-verse/air)
- **Containerization**: Docker

## 🏁 Getting Started

### Prerequisites

- Go 1.22+
- PostgreSQL database
- [Swag CLI](https://github.com/swaggo/swag) (for documentation updates)
- [Air](https://github.com/air-verse/air) (for local development with hot reload)

### Local Development

1.  **Clone the repository**:
    ```bash
    git clone <repository-url>
    cd short-url/backend
    ```

2.  **Set up environment variables**:
    Create a `.env` file in the root directory (or use the existing one) with your database credentials:
    ```env
    DB_POSTGRES_HOST=your-host
    DB_POSTGRES_PORT=5432
    DB_POSTGRES_USER=your-user
    DB_POSTGRES_PASS=your-password
    DB_POSTGRES_NAME=your-db-name
    DB_POSTGRES_SCHEMA=shortener
    
    KEY_LEN=5
    KEY_TRY=10
    ```

3.  **Run the application**:
    You can use the provided `Makefile`:
    ```bash
    make run
    ```
    This will initialize Swagger docs and start the server with Air for live reloading.

4.  **Running Tests**:
    ```bash
    make test
    ```

### API Documentation

Once the server is running, you can access the interactive Swagger UI at:
[http://localhost:8080/swagger/index.html](http://localhost:8080/swagger/index.html)

The database health check is available at `GET /health`. It returns `200` when PostgreSQL is reachable and `503` otherwise.

## 🐳 Docker Deployment

You can run the service using Docker and Docker Compose:

```bash
# Start all services (Backend + DB)
docker-compose up -d
```

## ⚙️ Environment Variables

| Variable | Description | Default |
| :--- | :--- | :--- |
| `DB_POSTGRES_HOST` | Database host address. | - |
| `DB_POSTGRES_PORT` | Database port. | `5432` |
| `DB_POSTGRES_USER` | Database username. | - |
| `DB_POSTGRES_PASS` | Database password. | - |
| `DB_POSTGRES_NAME` | Database name. | - |
| `DB_POSTGRES_SCHEMA` | Database schema name. | `shortener` |
| `KEY_LEN` | Length of the generated short key. | `5` |
| `KEY_TRY` | Max attempts to generate a unique key. | `10` |

## 📂 Project Structure

- `cmd/`: Application entry points (`main.go`).
- `internal/adapters/`: Implementation of external drivers (HTTP handlers, Repositories).
- `internal/core/`: Business logic, services, and domain models.
- `docs/`: Generated Swagger documentation.
- `utils/`: Helper functions and utilities.

---

Developed with ❤️ for a scalable web.
