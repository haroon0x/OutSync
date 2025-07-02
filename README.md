# OutSync

Reliable, asynchronous event forwarding system built on the transactional outbox pattern.

OutSync is a modular backend system that ensures reliable delivery of database events to message queues like Kafka using the transactional outbox pattern. It gives you a production-ready async infrastructure to build event-driven systems, webhooks, triggers, and automation platforms.

Whether you're building a Zapier-like workflow engine, distributed services, or AI pipelines — OutSync is the base layer that ensures your events never get lost between your database and your queue.

----------------------------

KEY FEATURES

- Async-first architecture using FastAPI + Tortoise ORM
- Implements transactional outbox pattern for safe DB-to-queue delivery
- Outbox polling and dispatch logic built-in
- Connects to Kafka out of the box (easy to extend to Redis/NATS/etc.)
- Fully containerized with Docker & Docker Compose
- Built to be copied, forked, reused, extended in other projects

----------------------------

TECH STACK

Web Framework: FastAPI (async)  
ORM: Tortoise ORM  
Database: PostgreSQL  
Queue: Kafka  
Worker: Custom asyncio loop  
Migrations: Aerich  
Containers: Docker + Docker Compose

----------------------------

ARCHITECTURE

[ FastAPI App ]
      |
      |  ⬇
[ DB Transaction ]
  - Inserts business data
  - Appends event to `outbox_events` table
      |
      | (Async Worker Polls)
      ⬇
[ OutSync Worker ]
  - Reads unprocessed events
  - Publishes to Kafka
  - Marks as processed

----------------------------

USE CASES

- Webhook and trigger systems (Zapier-style)
- Async AI agents triggered by structured input
- Multi-service orchestration (event bus)
- Eventual consistency between microservices
- Audit logs / event sourcing backends

----------------------------

QUICK START

1. Clone the repo

git clone https://github.com/your-username/outsync.git
cd outsync

2. Start the stack

docker-compose up --build

3. Run migrations

docker exec -it outsync-api aerich upgrade

4. Trigger sample event (via FastAPI route)

Send a test trigger to the API — it’ll write to the DB and queue an event for Kafka.

5. Worker will auto-forward

The background worker polls the DB and forwards the event to Kafka (or any queue you plug in).

----------------------------

FUTURE EXTENSIONS

- Multi-queue support (Redis, NATS, etc.)
- Retry & backoff strategies
- Event schema validation (Pydantic)
- Webhook dispatcher
- Built-in metrics (Prometheus)

----------------------------

PROJECT STRUCTURE

outsync/
├── app/
│   ├── api/            --> FastAPI routes
│   ├── models/         --> Tortoise ORM models
│   ├── core/           --> Configs, DB setup
│   └── services/       --> Kafka publishing logic
├── worker/             --> Async outbox poller
├── docker-compose.yml
├── requirements.txt
└── README.md

----------------------------

PHILOSOPHY

OutSync doesn’t try to be a framework. It’s a starter core — a copy-pasteable, extendable, battle-tested blueprint for any backend project that relies on safe, async event delivery.

Don’t reinvent the wheel. Fork this repo, rename it, and plug in your own business logic.

----------------------------

