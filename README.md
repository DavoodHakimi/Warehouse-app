# Warehouse Management System

A web-based warehouse management application that handles inventory, orders, business partners, and users — with role-based access control. Built as a monorepo with a Go backend and a Next.js frontend.

## Features

- **Order Management** — Create and track purchase and sale orders through a full lifecycle: Pending → Approved → Packed → Shipped / Received. Order numbers are auto-generated.
- **Product Management** — Add products with auto-generated product numbers, set pricing, and track stock changes as orders move through their lifecycle.
- **Business Partners** — Manage suppliers, customers, or both. Store contact info and associate them with orders.
- **User Management & RBAC** — Multi-tenant user system with five roles (CEO, Warehouse Manager, Storeman-Full, Storeman-EnterOnly, Storeman-ExitOnly). Each role has granular permissions that control access to every feature.
- **Stock Tracking** — Available and reserved stock levels update automatically when orders are approved, shipped, received, or canceled.
- **Audit Trail** — All updates to users, products, partners, and orders are logged with before/after values.
- **Dashboard** — A clean UI showing recent orders, products, and partners with quick access to all sections.

## Tech Stack

| Layer | Technology |
|-------|-----------|
| Backend | Go, Gin, GORM |
| Frontend | Next.js, React, TypeScript, Tailwind CSS |
| Database | PostgreSQL 16 |
| Reverse Proxy | Nginx |
| Auth | JWT (HS256, 24h expiry), bcrypt |
| Containerization | Docker, docker-compose |

## Quick Start

### Docker (recommended)

```bash
docker compose up --build
```

Open [http://localhost](http://localhost) — nginx serves the frontend and proxies `/api/` requests to the backend.

### Manual

#### Backend

```bash
cd backend
DBUSER="YOUR_USERNAME" DBPASS="YOUR_PASSWORD" DBHOST="localhost" DBPORT="PORT" DBNAME="DB_NAME" JWT_SECRET="your-secret" go run cmd/api/main.go
```

The server starts on `http://localhost:8080` and auto-migrates the database schema and seeds reference data (currencies, roles, permissions).

#### Frontend

```bash
cd frontend
npm install
npm run dev
```

The app runs on `http://localhost:3000`. Create an account, and you're in.

## System Architecture

```
                      ┌─────────┐
                      │  Nginx  │  (port 80)
                      │  :80    │
                      └────┬────┘
                           │
               ┌───────────┴───────────┐
               ▼                       ▼
        ┌──────────┐           ┌──────────────┐
        │ Backend  │           │   Frontend   │
        │ Go :8080 │           │ Next.js :3000│
        └────┬─────┘           └──────────────┘
             │
             ▼
      ┌──────────────┐
      │  PostgreSQL  │
      │  :PORT       │
      └──────────────┘
```

Nginx routes `/api/` requests to the Go backend and everything else to the Next.js frontend, so from the browser's perspective everything comes from the same origin — no CORS issues.

## Project Structure

```
backend/
├── cmd/api/          — Entry point
├── internal/
│   ├── auth/         — Authentication (login, signup, JWT)
│   ├── users/        — User management
│   ├── partners/     — Business partner CRUD
│   ├── products/     — Product CRUD
│   ├── orders/       — Order CRUD + status transitions
│   ├── middleware/    — Auth & RBAC middleware
│   └── database/     — Migrations, seed data
frontend/
└── app/              — Next.js pages (login, signup, dashboard)
nginx/
└── default.conf      — Nginx config
docker-compose.yaml   — Orchestrates all services
```

## License

MIT
