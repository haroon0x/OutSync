# OutSync Architecture Documentation

Welcome to the **OutSync** architecture documentation. This document provides an in-depth breakdown of the current system design, codebase layout, data flows, and technical implementation of the Transactional Outbox Pattern in Go.

---

## 1. System Overview

**OutSync** is a reliable, asynchronous event-forwarding engine built to guarantee **at-least-once event delivery**. In distributed systems, microservices, and AI agent architectures, performing a database write and publishing a corresponding notification event to a message queue or broker must be atomic. 

If a system writes to the database but crashes before publishing the event, downstream systems (like email notify servers, indexers, or AI agents) will never know the change occurred, leaving the system in an inconsistent state.

OutSync solves this problem using the **Transactional Outbox Pattern**:
1. **Atomicity:** The primary business data (e.g., user profiles) and the corresponding event payload (the "outbox event") are saved to the same PostgreSQL database inside a single ACID-compliant database transaction.
2. **Reliable Dispatching:** An independent background worker polls the database for unprocessed events, publishes them to a queue/broker (in this version, Redis), and marks them as processed upon success.

---

## 2. System Architecture Diagram

The diagram below represents the exact operational architecture and event flow implemented in the Go codebase.

```mermaid
flowchart TB
    subgraph ClientContext ["Client & Application Context"]
        Client["Client / main.go Test Executer"]
    end

    subgraph RelationalDB ["PostgreSQL Database (ACID Boundaries)"]
        direction TB
        subgraph DBTransaction ["PostgreSQL Transaction (Begin / Commit)"]
            UsersTable[("users table<br>(JSONB storage)")]
            OutboxTable[("outbox_events table<br>(aggregate_id, payload, status)")]
        end
    end

    subgraph Dispatcher ["OutSync Background Dispatcher"]
        Worker["worker.Poll() Loop<br>(internal/worker/worker.go)"]
    end

    subgraph QueueBroker ["Queue Broker (Redis)"]
        RedisList[("Redis Server<br>(List: 'events_queue')")]
    end

    %% Event Flow Links
    Client -->|1. Call CreateUserWithEvent()| DBTransaction
    DBTransaction -->|2a. Save User Data| UsersTable
    DBTransaction -->|2b. Write Event status='pending'| OutboxTable
    
    Worker -->|3. GetPendingEvents (LIMIT 10)| OutboxTable
    Worker -->|4. RPush event.Payload| RedisList
    Worker -->|5. MarkEventProcessed status='processed'| OutboxTable

    %% Styling
    classDef client fill:#e0f2fe,stroke:#0284c7,stroke-width:2px;
    classDef db fill:#fef3c7,stroke:#d97706,stroke-width:2px;
    classDef worker fill:#f3e8ff,stroke:#7c3aed,stroke-width:2px;
    classDef queue fill:#dcfce7,stroke:#16a34a,stroke-width:2px;
    
    class Client client;
    class UsersTable,OutboxTable db;
    class Worker worker;
    class RedisList queue;
```

---

## 3. Project & Package Directory Structure

The Go codebase is designed with clean architecture principles, storing private application code under Go's special `internal/` directory.

```
OutSync/
├── cmd/
│   └── outsync/
│       └── main.go         # Application entry point. Loads config, runs database migrations, and performs a single test database write transaction.
├── internal/
│   ├── api/
│   │   └── models.go       # Placeholder for API data models (currently empty).
│   ├── config/
│   │   └── config.go       # System configuration struct. Loads environment variables (.env) using godotenv.
│   ├── database/
│   │   ├── db.go           # Handles Postgres database connections and executes schema.sql migrations.
│   │   ├── queries.go      # Business queries: implements Transactional Outbox insertion, polling, and status updating.
│   │   └── schema.sql      # Database schema defining the SQL tables for users and outbox events.
│   └── worker/
│       └── worker.go       # Background loop querying Postgres pending events, pushing to Redis list, and updating status.
├── docs/
│   └── intro.md            # Concept document introducing the transactional outbox pattern.
├── docker-compose.yml      # Starts PostgreSQL (postgres:16-alpine) and Redis (redis:8-alpine) for local development.
├── go.mod                  # Go module definition (Go 1.25.1).
├── go.sum                  # Go module checksum tracking.
├── task.md                 # Project implementation plan and checklist indicating completed tasks.
└── README.md               # [OUTDATED] Placeholder README reflecting a different implementation (Python/FastAPI/Kafka).
```

---

## 4. Package Responsibilities & Code Breakdown

### `cmd/outsync/main.go`
* **Purpose:** The system bootstrap.
* **Execution:**
  1. Bootstraps the application by loading environment configuration.
  2. Runs migrations by executing `internal/database/schema.sql` against the database.
  3. Connects to PostgreSQL using `pgx`.
  4. Triggers a test insertion by executing `database.CreateUserWithEvent` with a sample prompt.

### `internal/config/config.go`
* **Struct:** `Config { DatabaseUrl string, GeminiAPIKey string }`
* **Responsibility:** Loads configuration variables from `.env` using the `github.com/joho/godotenv` package. 
* **Note:** `GEMINI_API_KEY` is loaded but currently unused.

### `internal/database`
* **`db.go`:** Configures database connections via `github.com/jackc/pgx/v5`. Implements `ApplySchema()`, which reads `schema.sql` and updates the Postgres database structure.
* **`schema.sql`:**
  * `users`: `id UUID PRIMARY KEY`, `data JSONB NOT NULL`.
  * `outbox_events`: `id UUID PRIMARY KEY`, `aggregate_id UUID NOT NULL`, `aggregate_type VARCHAR(255) NOT NULL`, `payload JSONB NOT NULL`, `status VARCHAR(255) NOT NULL` (e.g. `pending`, `processed`), `created_at TIMESTAMP`.
* **`queries.go`:** Implements core SQL logic:
  * `CreateUserWithEvent()`: Performs transaction-bound operations:
    1. Starts a database transaction.
    2. Inserts new user record into the `users` table and retrieves the generated UUID.
    3. Inserts an event containing the payload into the `outbox_events` table with `status = 'pending'`, referenced to the user's UUID.
    4. Commits the transaction. Rollback is triggered automatically if any step fails.
  * `GetPendingEvents()`: Fetches up to 10 outbox events where `status = 'pending'`.
  * `MarkEventProcessed()`: Updates status to `'processed'` for a given event ID.

### `internal/worker/worker.go`
* **Function:** `Poll(ctx, conn, rdb)`
* **Responsibility:**
  * Runs a continuous polling loop.
  * Queries up to 10 pending events from PostgreSQL.
  * For each event, uses `rdb.RPush()` to push the event payload to a Redis list key called `events_queue`.
  * On a successful Redis write, executes `database.MarkEventProcessed()` to update PostgreSQL status.
  * Introduces a 5-second sleep interval if no events are found or if an error is encountered.

---

## 5. Main Data & Event Flow

```mermaid
sequenceDiagram
    autonumber
    actor Client as Client/Main
    participant Postgres as PostgreSQL DB
    participant Worker as Background Worker
    participant Redis as Redis Queue

    Note over Client,Postgres: Ingestion & Atomicity Stage (Tx)
    Client->>Postgres: Begin Transaction
    Client->>Postgres: Insert into users table
    Postgres-->>Client: Return generated user UUID
    Client->>Postgres: Insert into outbox_events table (status: 'pending')
    Client->>Postgres: Commit Transaction
    Postgres-->>Client: Success!

    Note over Worker,Redis: Polling & Propagation Stage
    loop Every 5s or when events found
        Worker->>Postgres: SELECT pending events (LIMIT 10)
        Postgres-->>Worker: List of pending events
        loop For each pending event
            Worker->>Redis: RPush event payload to "events_queue"
            Redis-->>Worker: Pushed successfully
            Worker->>Postgres: UPDATE status to 'processed'
            Postgres-->>Worker: Status updated
        end
    end
```

---

## 6. Infrastructure & Running Setup

### External Dependencies
* **PostgreSQL (v16):** Serves as the single source of truth for both business data (`users`) and event logs (`outbox_events`). This is crucial because it allows the application to utilize PostgreSQL's ACID transactional features.
* **Redis (v8):** Serves as the message broker. The background worker publishes serialized JSON payloads to the Redis List structure at key `events_queue`.

### Configuration Configuration & Environment
The application expects a `.env` file at the root directory containing:
```env
DATABASE_URL=postgres://postgres:postgres@localhost:5432/outsync?sslmode=disable
GEMINI_API_KEY=your_gemini_api_key_here
```

### Running Locally
To run the project locally, execute the following commands:

1. **Spin up Infrastructure:**
   ```bash
   docker-compose up -d
   ```
2. **Execute Go Entrypoint:**
   ```bash
   go run cmd/outsync/main.go
   ```

---

## 7. Critical Architecture Discrepancies & Gaps

During our deep analysis, we identified several structural discrepancies, unused configurations, and unimplemented packages that need correction or are part of the future development tasks listed in `task.md`:

1. **Readme vs Codebase Mismatch:**
   * The root `README.md` details a **Python (FastAPI + Tortoise ORM + Aerich + Kafka)** implementation. This is completely outdated, as the current repository is fully rewritten in **Go (1.25+) using pgx, go-redis, and native SQL**.
2. **Unwired Background Worker:**
   * The background worker poller (`internal/worker/worker.go`) is fully written but is **never invoked or imported** in the main execution path `cmd/outsync/main.go`. Running `main.go` only applies migrations, inserts one test user/event record, and stops.
3. **Missing Redis Config Setup:**
   * The worker requires a connection to a `redis.Client`, but the configuration loader (`internal/config/config.go`) does not currently parse a `REDIS_URL` or `REDIS_ADDR` env variable, even though `task.md` lists `REDIS_URL` as a step.
4. **Empty API & HTTP handlers:**
   * The `internal/api/models.go` file is empty. The HTTP API routing server (`api/handler.go`) mentioned in `task.md` Phase 5 has not yet been implemented.
