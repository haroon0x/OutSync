# OutSync - Go Implementation Guide (Redis Edition)

This task list is designed to guide you through building OutSync step-by-step. Since you are learning Go, each task is broken down into small, manageable chunks.

## Phase 1: Project Setup & "Hello World"
**Goal:** Get a running Go program and understand the folder structure.

- [x] **Initialize the Project** <!-- id: 0 -->
    - [x] Open your terminal in `e:\OutSync`.
    - [x] Run `go mod init outsync`. (This creates `go.mod`, which tracks your dependencies).
- [x] **Create the Entry Point** <!-- id: 1 -->
    - [x] Create a file `cmd/outsync/main.go`.
    - [x] Write a `package main` and a `func main()`.
    - [x] Make it print "Hello, OutSync!" using `fmt.Println`.
    - [x] Run it: `go run cmd/outsync/main.go`.
- [x] **Create the Folder Structure** <!-- id: 2 -->
    - [x] Create folders: `internal/config`, `internal/database`, `internal/api`, `internal/worker`.
    - [x] *Why?* In Go, `internal` is a special folder that cannot be imported by other projects. It keeps your code private.

## Phase 2: Configuration & Environment
**Goal:** Learn how to read settings (like DB passwords) from the environment.

- [x] **Install `godotenv`** <!-- id: 3 -->
    - [x] Run `go get github.com/joho/godotenv`.
- [x] **Create `.env` file** <!-- id: 4 -->
    - [x] Create a file named `.env` in the root.
    - [x] Add: `DATABASE_URL=postgres://outsync_user:outsync_password@localhost:5432/outsync_db`
    - [x] Add: `REDIS_URL=localhost:6379`
- [x] **Write Config Code** <!-- id: 5 -->
    - [x] Create `internal/config/config.go`.
    - [x] Define a struct: `type Config struct { DatabaseURL string; RedisURL string }`.
    - [x] Write a function `func Load() (*Config, error)`.
    - [x] Use `godotenv.Load()` to read the file, then `os.Getenv()` to fill your struct.
    - [x] **Test it:** Update `main.go` to call `config.Load()` and print the config.

## Phase 3: Infrastructure (Docker)
**Goal:** Spin up the database and Redis so we have something to connect to.

- [x] **Create `docker-compose.yml`** <!-- id: 6 -->
    - [x] Create a service for Postgres.
    - [x] Create a service for Redis (image: `redis:alpine`).
    - [x] Run `docker-compose up -d`.
    - [x] Verify they are running with `docker ps`.

## Phase 4: Database Layer (The Foundation)
**Goal:** Connect to Postgres and create tables.

- [x] **Install `pgx` (Postgres Driver)** <!-- id: 7 -->
    - [x] Run `go get github.com/jackc/pgx/v5`.
- [x] **Connect to DB** <!-- id: 8 -->
    - [x] Create `internal/database/db.go`.
    - [x] Write a function `func NewConnection(url string) (*pgx.Conn, error)`.
    - [x] Use `pgx.Connect(context.Background(), url)`.
- [x] **Create Tables (Migration)** <!-- id: 9 -->
    - [x] Create `internal/database/schema.sql`.
    - [x] Write SQL to create table `users` (id UUID, email TEXT).
    - [x] Write SQL to create table `outbox_events` (id UUID, payload JSONB, status TEXT, created_at TIMESTAMP).
    - [x] *Bonus:* Write a Go function in `db.go` that reads this file and executes it using `conn.Exec`. Call this from `main.go`.

## Phase 5: The "In" Side (API)
**Goal:** Accept a request and save to DB + Outbox in ONE transaction.

- [ ] **Define the Models** <!-- id: 10 -->
    - [ ] Create `internal/api/models.go`.
    - [ ] Define `type User struct { ... }`.
    - [ ] Define `type CreateUserRequest struct { Email string }`.
- [x] **Implement the Transaction Logic** <!-- id: 11 -->
    - [x] Create `internal/database/queries.go`.
    - [x] Write a function `CreateUserWithEvent(ctx, conn, email string) error`.
    - [x] **Step 1:** Start transaction: `tx, _ := conn.Begin(ctx)`.
    - [x] **Step 2:** Insert User: `tx.Exec(..., "INSERT INTO users ...")`.
    - [x] **Step 3:** Insert Event: `tx.Exec(..., "INSERT INTO outbox_events ...")`.
    - [x] **Step 4:** Commit: `tx.Commit(ctx)`.
- [ ] **Build the HTTP Handler** <!-- id: 12 -->
    - [ ] Create `internal/api/handler.go`.
    - [ ] Write `func HandleCreateUser(w http.ResponseWriter, r *http.Request)`.
    - [ ] Parse JSON body using `json.NewDecoder(r.Body).Decode(...)`.
    - [ ] Call your `CreateUserWithEvent` function.
- [ ] **Start the Server** <!-- id: 13 -->
    - [ ] In `main.go`, use `http.HandleFunc("/users", api.HandleCreateUser)`.
    - [ ] Start listening: `http.ListenAndServe(":8080", nil)`.

## Phase 6: The "Out" Side (Worker)
**Goal:** Read from Outbox and send to Redis.

- [x] **Install Redis Driver** <!-- id: 14 -->
    - [x] Run `go get github.com/redis/go-redis/v9`.
- [/] **Implement the Poller** <!-- id: 15 -->
    - [ ] Create `internal/worker/poller.go`.
    - [ ] Write a loop: `for { ... }`.
    - [ ] Inside loop: `SELECT * FROM outbox_events WHERE status='pending' LIMIT 10`.
    - [ ] *Tip:* Use `time.Sleep` at the end of the loop so you don't hammer the DB.
- [ ] **Implement the Publisher** <!-- id: 16 -->
    - [ ] Inside the loop, for each event:
    - [ ] Create a Redis client.
    - [ ] `client.RPush(ctx, "events_queue", eventPayload)`.
- [ ] **Mark as Processed** <!-- id: 17 -->
    - [ ] If Redis write is successful: `UPDATE outbox_events SET status='processed' WHERE id=$1`.

## Phase 7: Run & Verify
- [ ] **Run the App** <!-- id: 18 -->
    - [ ] `go run cmd/outsync/main.go`.
- [ ] **Test** <!-- id: 19 -->
    - [ ] Send POST request: `curl -X POST -d '{"email":"test@test.com"}' http://localhost:8080/users`.
    - [ ] Check DB: `SELECT * FROM outbox_events;` (Should be 'processed').
    - [ ] Check Redis: `LRANGE events_queue 0 -1` (Should see the event).
