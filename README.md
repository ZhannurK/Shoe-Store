# Shoe‑Store Micro‑services Project

A **cloud‑native online sneaker shop** built as a constellation of lightweight Go micro‑services.  The project demonstrates clean service boundaries, event‑driven communication, and full production tooling (observability, CI, containerisation).

---

## 📚 Project Overview & Topic

This repository is a course capstone showcasing how to design and implement an **e‑commerce platform** (focused on sneaker sales) using modern micro‑service patterns. Its learning goals are:

* Practise gRPC + Protocol Buffers for service‑to‑service APIs.
* Apply MongoDB + Redis in a CQRS style (document store + cache/queue).
* Orchestrate services with **Docker Compose** for local dev and CI.
* Demonstrate operational concerns: metrics, health checks, graceful shutdowns.

---

## 🛠️ Technologies Used

| Layer              | Choices                               |
| ------------------ | ------------------------------------- |
| Language           | **Go 1.23**                           |
| Service Frameworks | Gin (Gateway), gRPC, Protocol Buffers |
| Data Stores        | MongoDB 6 / Atlas, Redis 7            |
| Messaging          | NATS JetStream                        |
| Observability      | Prometheus v2, Grafana v10            |
| Containerisation   | Docker 24, Docker‑Compose v3.8        |
| CI / Lint          | GitHub Actions, golangci‑lint         |

---

## ✨ Implemented Features

* **User Authentication** – e‑mail signup/login, JWT issuance & validation, password hashing (bcrypt), token blacklist cache.
* **Product Catalogue & Inventory** – CRUD for sneaker models, stock reservation, low‑stock alerts via NATS.
* **Order Processing** – two‑phase reservation + payment simulation with idempotent saga log stored in Redis/Mongo.
* **API Gateway** – single public REST front; routes requests to internal gRPC services and exposes an OpenAPI (Swagger) UI.
* **Asynchronous Events** – `order.created`, `order.paid`, `low_stock` subjects published on NATS JetStream.
* **Observability Suite** – unified Prometheus scrape config, Grafana dashboard (latency, error rates, business KPIs).
* **Local & CI Tooling** – Docker Compose stack, makefile helpers (`make proto`, `make dev-logs`), GitHub Actions pipeline.

---

## 🗺️ High‑Level Architecture

```text
┌────────────┐    HTTP/JSON     ┌──────────────┐
│  Frontend  │  ─────────────▶ │  API Gateway │
└────────────┘                  └─────┬────────┘
                                      │ gRPC
             ┌─────────────────────────┼────────────────────────┐
             ▼                         ▼                        ▼
      ┌──────────────┐          ┌───────────────┐        ┌────────────────┐
      │ Auth Service │          │ Inventory Svc │        │ Transaction Svc│
      └──────┬───────┘          └──────┬────────┘        └──────┬─────────┘
             │ JWT/Redis              │ Mongo/Redis           │ Mongo/Redis
             ▼                         ▼                       ▼
         ┌────────┐              ┌────────┐              ┌────────┐
         │ Redis  │              │MongoDB │              │ Redis  │
         └────────┘              └────────┘              └────────┘
                   ▲                  ▲                        ▲
                   │ NATS events      │ Low‑stock alerts       │ Order saga
                   └──────────────────┴────────────────────────┘
```

---

## 🏃‍♀️ Running Locally

### With Docker (recommended)

```bash
# Clone & enter repo
$ git clone https://github.com/your‑org/shoe‑store.git && cd Shoe‑Store‑master

# Spin up full stack
$ docker compose up --build -d

# Tail logs
$ make dev-logs
```

Open‑ports cheat‑sheet:

| Service          | URL                                            | Notes                              |
| ---------------- | ---------------------------------------------- | ---------------------------------- |
| Gateway (REST)   | [http://localhost:8181](http://localhost:8181) | Swagger at `/swagger/index.html`   |
| Auth gRPC        | `localhost:8087`                               | reflection enabled                 |
| Inventory gRPC   | `localhost:5052`                               | Prometheus metrics `:5053/metrics` |
| Transaction gRPC | `localhost:8088`                               |                                    |
| Prometheus       | [http://localhost:9090](http://localhost:9090) |                                    |
| Grafana          | [http://localhost:3000](http://localhost:3000) | admin / admin                      |

### Manual (no Docker)

1. Start MongoDB & Redis locally.
2. Export env vars (see `.env.sample`).
3. Run services in separate terminals:

```bash
$ go run cmd/main.go   # inside each service folder
```

---

## 🧪 Running Tests

```bash
# Run all unit tests
$ go test ./...

# Lint (requires golangci‑lint)
$ golangci-lint run ./...
```

CI executes the same commands on every pull request.

---

## 🔌 gRPC API Reference

Below is a condensed list of RPCs; for full detail open the generated **Swagger** (Gateway) or run `grpcurl list <svc>`.

### Auth Service (`auth.proto`)

| RPC        | Request                                   | Response                           | Purpose               |
| ---------- | ----------------------------------------- | ---------------------------------- | --------------------- |
| `Signup`   | `SignupRequest { email, password, name }` | `AuthReply { user_id }`            | Register new user     |
| `Login`    | `LoginRequest { email, password }`        | `TokenReply { jwt }`               | Issue JWT             |
| `Confirm`  | `ValidateRequest { jwt }`                 | `ValidateReply { valid, user_id }` | Verify token validity |

### Inventory Service (`inventory.proto`)

| RPC            | Request                               | Response                   | Purpose                |
| -------------- | ------------------------------------- | -------------------------- | ---------------------- |
| `ListSneakers` | `PageQuery { page, size }`            | `SneakerList`              | Paginated catalogue    |
| `GetSneaker`   | `SneakerID`                           | `Sneaker`                  | Single product detail  |
| `ReserveStock` | `ReserveRequest { sneaker_id, qty }`  | `ReserveReply { success }` | Tentatively hold stock |
| `AdjustStock`  | `AdjustRequest { sneaker_id, delta }` | `AdjustReply { new_qty }`  | Admin stock update     |

### Transaction Service (`order.proto`)

| RPC           | Request                             | Response                  | Purpose                |
| ------------- | ----------------------------------- | ------------------------- | ---------------------- |
| `CreateOrder` | `OrderRequest { user_id, items[] }` | `OrderReply { order_id }` | Start order saga       |
| `GetOrder`    | `OrderID`                           | `Order`                   | Order status & details |
| `ListOrders`  | `UserID`                            | `OrderList`               | All orders for user    |

---

## 📈 Monitoring & Dashboards

* Prometheus scrapes `/metrics` every 15 s (config in `deploy/prometheus.yml`).
* Grafana is pre‑provisioned with **Shoe‑Store Overview** dashboard (ID 1). Import your own JSON to extend.

---

## 🤝 Contributing

1. **Fork** → `git checkout -b feature/xyz`
2. **Commit** using conventional messages.
3. **Push & PR** – CI must pass.
