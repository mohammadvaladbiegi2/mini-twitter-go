# Mini Twitter Backend (Go)

Chirper is a minimal yet production-oriented Twitter clone backend written in Go. It aims to demonstrate building a scalable microblogging REST API using modern best-practices: clean architecture, JWT auth, PostgreSQL for persistence, Redis for caching, background workers for notifications, Dockerized deployment, and CI pipelines. Ideal as a portfolio project and interview demo.

## Key features
- User signup/login with JWT (access + refresh)
- Create/read tweets, like tweets
- Follow / unfollow relationships
- Timeline generation (fan-out/fan-in options)
- Redis caching for hot timelines
- Docker Compose 


## 🚀 How to Run Locally

### Prerequisites
Before starting, make sure you have installed on your system:
- **Go 1.24+**
- **Docker**

---

### 1. Clone the project


```bash
git clone https://github.com/mohammadvaladbiegi2/mini-twitter-go.git
cd mini-twitter-go
```
2. Install Go dependencies


```bash
go mod tidy
```
3. Start services with Docker
This project includes PostgreSQL, Redis, MinIO, and pgAdmin via Docker Compose.



```bash
docker compose up -d
```
4. Prepare the database schema
A ready schema.sql file is provided in the project root.
Import it into the database (twitter_clone by default):



```bash
docker exec -i postgres_db psql -U admin -d twitter_clone < schema.sql
```
This will create all tables, relations, and indexes.

5. Run the backend API
After containers are running and the schema is imported, run the API:



```bash
go run cmd/api/main.go
```

6. Explore the API with Swagger
Open your browser at:

👉 http://localhost:7080/swagger/index.html

Here you can see all endpoints, request/response formats, and try out the APIs interactively.



🛠 Tech Stack
Language: Go (1.24+)

Database: PostgreSQL

Cache: Redis

Object Storage: MinIO (S3-compatible)

Auth: JWT (access + refresh tokens)

Docs: Swagger (OpenAPI)

📌 Notes
The .env file is included with all necessary configuration (no extra setup needed).

schema.sql ensures database structure is consistent — no need to run migrations manually.

Default services (Postgres, Redis, MinIO, pgAdmin) are all handled with Docker Compose.
