# GO-CRUD
# Go REST API with PostgreSQL and Docker

A simple RESTful API built using **Golang**, **Gorilla Mux**, and **PostgreSQL**. The app is containerized using **Docker** for seamless development and deployment.

## 📦 Features

- RESTful API to manage users (CRUD operations)
- PostgreSQL integration
- Dockerized setup with `docker-compose`
- Health check for PostgreSQL container
- Auto-creation of `people` table if not present
- Middleware to set JSON content type

## 🚀 Getting Started

### Prerequisites

- Docker
- Docker Compose

### Run the App

1. Clone this repository:

```bash
git clone https://github.com/your-username/go-rest-api.git
cd go-rest-api
```
---
Build and run the services
---
docker-compose up --build
---
API Endpoints
---
| Method | Endpoint       | Description             |
| ------ | -------------- | ----------------------- |
| GET    | `/people`      | Get all users           |
| GET    | `/people/{id}` | Get a user by ID        |
| POST   | `/people`      | Create a new user       |
| PUT    | `/people/{id}` | Update an existing user |
| DELETE | `/people/{id}` | Delete a user by ID     |

---
Sample JSON for POST/PUT:
---
```bash
{
  "name": "atul nath",
  "email": "atul@gmail.com"
}
```
---
Docker Overview
---
go-app: Go backend server (built from source)

go_db: PostgreSQL 12 with persistent volume

---
Environment Variables
---
| Variable           | Value                                                                       |
| ------------------ | --------------------------------------------------------------------------- |
| POSTGRES\_USER     | postgres                                                                    |
| POSTGRES\_PASSWORD | postgres                                                                    |
| POSTGRES\_DB       | postgres                                                                    |
| DATABASE\_URL      | host=go\_db user=postgres password=postgres dbname=postgres sslmode=disable |
---
Tech Stack
---
. Go

. Gorilla Mux

. PostgreSQL

. Docker

. Docker Compose
