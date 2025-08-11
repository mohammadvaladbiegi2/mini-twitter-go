# Mini Twitter Backend (Go)

Chirper is a minimal yet production-oriented Twitter clone backend written in Go. It aims to demonstrate building a scalable microblogging REST API using modern best-practices: clean architecture, JWT auth, PostgreSQL for persistence, Redis for caching, background workers for notifications, Dockerized deployment, and CI pipelines. Ideal as a portfolio project and interview demo.

## Key features
- User signup/login with JWT (access + refresh)
- Create/read tweets, like tweets
- Follow / unfollow relationships
- Timeline generation (fan-out/fan-in options)
- Redis caching for hot timelines
- Docker Compose + GitHub Actions for CI
- Tests: unit + integration

