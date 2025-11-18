# Twitter Feed Backend

A scalable, real-time Twitter-like feed backend system built with Go, implementing HTTP streaming, message queue-based backpressure, and a distributed database cluster.

## Table of Contents

- [Features](#features)
- [Architecture](#architecture)
- [Code Structure](#code-structure)
- [How to Run](#how-to-run)
- [API Endpoints](#api-endpoints)
- [Technologies](#technologies)
- [Configuration](#configuration)

## Features

✅ **RESTful API** for message creation  
✅ **HTTP Server-Sent Events (SSE)** for real-time feed streaming  
✅ **Kafka-based backpressure** for reliable message processing  
✅ **CockroachDB cluster** (2-node) for distributed data storage  
✅ **Message generation bot** with configurable speed  
✅ **One-command startup** with Docker Compose  
✅ **Graceful shutdown** with proper resource cleanup  
✅ **Database migrations** with version control  

## Architecture

### Components

#### API Service
- **Handlers**: HTTP request processing
- **Messaging**: Kafka producer for event publishing
- **Repository**: Database operations and migrations
- **Worker**: Kafka consumer for message processing
- **Broadcaster**: Central relay for streaming messages to connected clients

#### Bot Service
- **Generator**: Random message content generation
- **Scheduler**: Configurable interval-based task execution
- **Client**: HTTP client for API communication

#### Infrastructure
- **Kafka**: Message queue with 2 topics for backpressure
- **CockroachDB**: 2-node distributed SQL database cluster
- **Docker Compose**: Container orchestration

### Data Flow

1. **Message Creation Flow:**
   - Bot/Client → POST `/api/messages`  Message Handler
   - Message Handler → Kafka Producer (Topic: `events-to-process`)
   - Returns `202 Accepted` immediately

2. **Message Processing Flow:**
   - Kafka (Topic: `events-to-process`) → Worker Consumer
   - Worker Consumer → CockroachDB
   - Worker Consumer → Kafka Producer (Topic: `events-processed`)

3. **Real-time Streaming Flow:**
   - Client → GET `/api/feed` → Feed Handler
   - Feed Handler → Send historical messages from DB
   - Kafka (Topic: `events-processed`) → Subscriber Consumer
   - Subscriber Consumer → Broadcaster → SSE Stream → Client

## Code Structure

```
.
├── api/                              # Main API service
│   ├── cmd/
│   │   └── api/
│   │       └── main.go               # API entry point
│   ├── internal/
│   │   ├── app/
│   │   │   ├── app.go                # Application lifecycle management
│   │   │   └── contract.go           # Interface definitions (dependency inversion)
│   │   ├── handler/
│   │   │   ├── router.go             # HTTP route definitions
│   │   │   ├── message.go            # POST /api/messages handler
│   │   │   ├── feed.go               # GET /api/feed handler (SSE)
│   │   │   ├── health.go             # Health check endpoint
│   │   │   ├── broadcaster.go        # Central relay for streaming messages to connected clients
│   │   │   ├── subscriber.go         # Kafka consumer for processed events
│   │   │   └── types.go              # Handler interfaces
│   │   ├── messaging/
│   │   │   ├── message.go            # Event wrapper with metadata
│   │   │   └── producer.go           # Kafka producer implementation
│   │   ├── repository/
│   │   │   ├── connection.go         # Database connection pool
│   │   │   ├── message.go            # Message entity (private fields)
│   │   │   ├── repository.go         # Database operations
│   │   │   └── migrate.go            # Migration runner
│   │   └── worker/
│   │       ├── worker.go             # Worker orchestration
│   │       ├── consumer.go           # Kafka consumer
│   │       ├── processor.go          # Message processing logic
│   │       └── types.go              # Worker interfaces
│   ├── migrations/
│   │   ├── 000001_create_messages_table.up.sql
│   │   └── 000001_create_messages_table.down.sql
│   ├── Dockerfile
│   ├── go.mod
│   └── go.sum
│
├── bot/                              # Message generation bot
│   ├── cmd/
│   │   └── bot/
│   │       └── main.go               # Bot entry point
│   ├── internal/
│   │   ├── bot/
│   │   │   ├── bot.go                # Bot orchestration
│   │   │   ├── factory.go            # Message factory
│   │   │   └── types.go              # Bot interfaces
│   │   ├── client/
│   │   │   ├── client.go             # HTTP client implementation
│   │   │   └── sendable.go           # Sendable interface
│   │   ├── generator/
│   │   │   └── generator.go          # Random message generator
│   │   └── scheduler/
│   │       └── scheduler.go          # Interval-based task scheduler
│   ├── Dockerfile
│   └── go.mod
│
├── scripts/
│   ├── init-cluster.sh               # CockroachDB cluster initialization
│   └── create-kafka-topics.sh        # Kafka topics creation
│
├── docker-compose.yaml               # Container orchestration
├── start.sh                          # One-command startup script
└── README.md
```

## How to Run

### Prerequisites

- **Docker** (with Docker Compose)
- **Bash** (Git Bash on Windows, or WSL)

**No other installations required!**

### Start the System

Run the following command from the project root:

```bash
bash start.sh
```

Or directly:

```bash
docker-compose up -d --build
```

### What Happens

1. **CockroachDB Cluster** (roach1, roach2) starts and forms a 2-node cluster
2. **Cluster Initialization** runs automatically via `roach-init` container
3. **Kafka** starts in KRaft mode (no Zookeeper required)
4. **Kafka Topics** (`events-to-process`, `events-processed`) are created
5. **API Service** starts, runs migrations, and connects to dependencies
6. **Bot Service** starts and begins generating messages every 10 seconds
7. **Worker** consumes messages from Kafka and persists to CockroachDB
8. **Subscriber** consumes processed events and broadcasts to SSE clients

### Verify System is Running

```bash
# Check container status
docker-compose ps

# Check API health
curl http://localhost:8090/api/health

# View API logs
docker-compose logs -f api

# View bot logs
docker-compose logs -f bot
```

### Access Services

| Service | URL | Purpose |
|---------|-----|---------|
| API | http://localhost:8090 | Main API endpoint |
| CockroachDB Node 1 UI | http://localhost:8080 | Database admin UI |
| CockroachDB Node 2 UI | http://localhost:8081 | Database admin UI |
| Kafka | localhost:9092 | Message broker |

### Stop the System

```bash
docker-compose down

# To remove volumes (clean state)
docker-compose down -v
```

## API Endpoints

### 1. Health Check

```http
GET /api/health
```

**Response:** `200 OK`
```json
{
  "status": "ok"
}
```

---

### 2. Add Message

```http
POST /api/messages
Content-Type: application/json
```

**Request Body:**
```json
{
  "user_id": "user123",
  "content": "Hello, Twitter!"
}
```

**Response:** `202 Accepted`

> **Note:** Returns immediately with `202 Accepted` status. Message is queued in Kafka and processed asynchronously (backpressure pattern).

**Example:**
```bash
curl -X POST http://localhost:8090/api/messages \
  -H "Content-Type: application/json" \
  -d '{"user_id": "john_doe", "content": "This is my first tweet!"}'
```

---

### 3. Get Feed (SSE Streaming)

```http
GET /api/feed
```

**Response:** Server-Sent Events (SSE) stream

**Behavior:**
1. Immediately sends all historical messages from database
2. Keeps connection open and streams new messages as they arrive
3. Messages are broadcast in real-time to all connected clients

**Response Format:**
```
data: {"id":"uuid","user_id":"user1","content":"Hello","created_at":"2024-..."}

data: {"id":"uuid","user_id":"user2","content":"World","created_at":"2024-..."}
```

## Technologies

| Technology | Version | Purpose |
|------------|---------|---------|
| **Go** | 1.23+ | Primary programming language |
| **Kafka** | 8.1.0 (KRaft) | Message queue for backpressure |
| **CockroachDB** | v25.3.4 | Distributed SQL database (2-node cluster) |
| **Docker** | 20.10+ | Containerization |
| **Docker Compose** | 3.8+ | Multi-container orchestration |

### Go Libraries

- `github.com/segmentio/kafka-go` - Kafka client
- `github.com/jackc/pgx/v5` - PostgreSQL/CockroachDB driver
- `github.com/golang-migrate/migrate/v4` - Database migrations
- `github.com/google/uuid` - UUID generation

## Configuration

### Environment Variables

#### API Service

| Variable | Default | Description |
|----------|---------|-------------|
| `PORT` | `8090` | API HTTP port |
| `KAFKA_HOST` | `kafka` | Kafka broker hostname |
| `KAFKA_PORT` | `29092` | Kafka broker port |
| `DB_HOST` | `roach1` | CockroachDB hostname |
| `DB_PORT` | `26257` | CockroachDB SQL port |

#### Bot Service

| Variable | Default | Description |
|----------|---------|-------------|
| `API_HOST` | `api` | API service hostname |
| `API_PORT` | `8090` | API service port |

### Bot Configuration

Edit `bot/cmd/bot/main.go` to adjust:

```go
const (
    UserCount     = 3                  // Number of simulated users
    SleepInterval = 10 * time.Second   // Time between messages
)
```

### Kafka Topics

- `events-to-process` - Incoming messages from API
- `events-processed` - Messages persisted to database

### CockroachDB Cluster

- **Node 1**: `roach1:26257` (SQL), `roach1:8080` (Admin UI)
- **Node 2**: `roach2:26258` (SQL), `roach2:8081` (Admin UI)
- **Replication Factor**: Automatic (distributed across nodes)

## Testing the System

### 1. Test Message Creation

```bash
curl -X POST http://localhost:8090/api/messages \
  -H "Content-Type: application/json" \
  -d '{"user_id": "alice", "content": "Testing the feed!"}'
```

### 2. Test Feed Streaming

Open multiple terminal windows and run:

```bash
curl -N http://localhost:8090/api/feed
```

You should see:
1. Historical messages immediately
2. New messages streaming in real-time as the bot generates them
3. All terminals receive the same messages simultaneously (fan-out)

### 3. Test Database Persistence

```bash
# Connect to CockroachDB
docker-compose exec roach1 cockroach sql --host=roach1:26257 --insecure

# Query messages
SELECT * FROM messages ORDER BY created_at DESC LIMIT 10;
```

### 4. Test Backpressure

```bash
# Send 100 messages rapidly
for i in {1..100}; do
  curl -X POST http://localhost:8090/api/messages \
    -H "Content-Type: application/json" \
    -d '{"user_id": "user$i", "content": "Message $i"}' &
done

# All should return 202 Accepted immediately
# Worker processes them asynchronously
```

### 5. Monitor Kafka Topics

```bash
# View topic messages
docker-compose exec kafka kafka-console-consumer \
  --bootstrap-server localhost:9092 \
  --topic events-processed \
  --from-beginning
```

## Troubleshooting

### API Won't Start

```bash
# Check if CockroachDB cluster is initialized
docker-compose logs roach-init

# Check if Kafka topics are created
docker-compose logs kafka-init

# Restart API
docker-compose restart api
```

### No Messages in Feed

```bash
# Check worker is consuming
docker-compose logs -f api 2>&1 | grep -i --color "Processor received"

# Check Kafka topics exist
docker-compose exec kafka kafka-topics --list --bootstrap-server localhost:9092
```

### Database Connection Issues

```bash
# Check CockroachDB nodes are healthy
docker-compose ps

# View database logs
docker-compose logs roach1
docker-compose logs roach2
```

**Built with ❤️ using Go, Kafka, and CockroachDB**

