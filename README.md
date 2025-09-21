# Student Management API

![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)

A simple **Student Management REST API** built with [Go](https://go.dev/) and [SQLite](https://www.sqlite.org/).
This project demonstrates clean architecture, structured logging, graceful shutdown, and modular design.

---

## ✨ Features

- Student CRUD endpoints (Create, Read, Update, Delete)
- Pagination support for listing students
- Graceful shutdown with signal handling
- Configurable via environment variables
- File + console logging (development) / file-only logging (production)
- Lightweight SQLite backend (easy to run anywhere)

---

## 📂 Project Structure

```
student_api/
├── cmd/               # Application entry point
├── config/            # Configuration (ignored in Git)
├── db/                # Database (ignored in Git)
├── internal/
│   ├── config/        # Configuration loading (env, YAML)
│   ├── db/            # SQLite database logic
│   ├── logger/        # Centralized logging
│   ├── response/      # JSON response helpers
│   ├── student/       # Student handlers
│   └── types/         # Domain models
├── logs/              # Log output (ignored in Git)
└── go.mod
```

---

## ⚙️ Configuration

The app is configured via a YAML configuration file (e.g., config.yaml).

Example config.yaml:

```yaml
env: "development"
db_path: "db/db.db"
http_server:
  host: "localhost"
  port: 8080
logging:
  level: "debug"
  file: "logs/app.log"
```

Example `.env` file:

```env
CONFIG_PATH="config/config.yaml"
```

---

## 🚀 Running the Project

### 1. Clone the repo

```bash
git clone https://github.com/MahfujulSagor/student_api.git
cd student_api
```

### 2. Install dependencies

```bash
go mod tidy
```

### 3. Run the server

```bash
go run cmd/student_api/main.go
```

Server runs at:
👉 `http://localhost:8080`

---

## If you dont setup `.env`

### Run the server like this

```bash
go run cmd/student_api/main.go -config config/config.yaml
```

Server runs at:
👉 `http://localhost:8080`

---

## 📡 API Endpoints

### Health Check

```bash
GET /
```

### Students

| Method   | Endpoint             | Description                     |
| -------- | -------------------- | ------------------------------- |
| `POST`   | `/api/students`      | Create a new student            |
| `GET`    | `/api/students`      | List students (with pagination) |
| `GET`    | `/api/students/{id}` | Get student by ID               |
| `PUT`    | `/api/students/{id}` | Update student by ID            |
| `DELETE` | `/api/students/{id}` | Delete student by ID            |

---

## 📖 Example Request / Response

### Create Student

```http
POST /api/students
Content-Type: application/json

{
  "name": "Alice",
  "email": "alice@example.com",
  "age": 20
}
```

Response:

```json
{
  "success": "OK",
  "message": "Student created with ID: 1"
}
```

---

## 🛠 Development Notes

- Logs are stored in `logs/app.log`
- In **development**, logs also print to console
- SQLite DB file defaults to `db.db` in the project root
- Graceful shutdown ensures ongoing requests complete within 10s

---

## 📜 License

MIT License. Free to use for learning and practice.
