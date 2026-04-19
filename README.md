# Nexo

Nexo is a self-hosted alternative to Notion. Organize documents, build flexible databases with multiple views (table, board, calendar, gallery), and automate workflows — all under your control.

## Getting Started

### Prerequisites

- Go 1.22+
- SQLite (default) or PostgreSQL

### Installation

```bash
go build -o nexo ./cmd
```

### Running

```bash
# Minimal — SQLite, required vars only
SESSION_SECRET_KEY="your-secret-at-least-32-chars-long" ./nexo server

# With config file
./nexo server --config config
```

---

## Configuration

All options can be set via environment variable, CLI flag, or `config.yaml`.
Priority: **env var > CLI flag > config file > default**.

### Required

| Env var | Description |
|---------|-------------|
| `SESSION_SECRET_KEY` | JWT signing key. **Minimum 32 characters.** The server refuses to start if this is shorter or unset. Generate with: `openssl rand -base64 48` |

### Server

| Env var | CLI flag | Default | Description |
|---------|----------|---------|-------------|
| `HTTP_PORT` | `--http.port` | `8080` | Listening port |
| `HTTP_LOGS` | `--http.logs` | `false` | Enable HTTP access logs |
| `HTTP_CORS_ALLOW_ORIGINS` | `--http.cors_allow_origins` | `*` | Comma-separated list of allowed CORS origins. **Set this in production** (e.g. `https://app.example.com`). Use `*` only for local dev. |

### Database

| Env var | CLI flag | Default | Description |
|---------|----------|---------|-------------|
| `DATABASE_DIALECT` | `--database.dialect` | `sqlite` | `sqlite` or `postgres` |
| `DATABASE_DSN` | `--database.dsn` | `./database.sqlite` | SQLite file path or PostgreSQL DSN |

### Session / JWT

| Env var | CLI flag | Default | Description |
|---------|----------|---------|-------------|
| `SESSION_SECRET_KEY` | `--session.secret_key` | *(none)* | **Required.** ≥ 32 chars |
| `SESSION_EXPIRATION_MINUTES` | `--session.expiration_minutes` | `43200` (30 days) | Token lifetime in minutes |
| `SESSION_ISSUER` | `--session.issuer` | `nexo` | JWT `iss` claim |

### Logger

| Env var | CLI flag | Default | Description |
|---------|----------|---------|-------------|
| `LOGGER_LEVEL` | `--logger.level` | `info` | `debug`, `info`, `warn`, `error` |
| `LOGGER_PRETTY` | `--logger.pretty` | `false` | Human-readable logs (dev only) |

---

## Config file (`config.yaml`)

```yaml
http:
  port: 8080
  logs: true
  cors_allow_origins: "https://app.example.com"

logger:
  level: info
  pretty: false

database:
  dialect: sqlite
  dsn: ./database.sqlite

session:
  # Required — do not commit real values to source control
  # Generate: openssl rand -base64 48
  secret_key: "CHANGE_ME_AT_LEAST_32_CHARACTERS_LONG"
  expiration_minutes: 43200
  issuer: nexo
```

See `config-example.yaml` for a minimal working example.

---

## Docker

```bash
docker build -t nexo .
docker run -p 8080:8080 \
  -e SESSION_SECRET_KEY="$(openssl rand -base64 48)" \
  -e HTTP_CORS_ALLOW_ORIGINS="https://app.example.com" \
  -e DATABASE_DSN="/data/nexo.sqlite" \
  -v nexo_data:/data \
  nexo
```

Or with `docker-compose.yaml`:

```bash
# Copy and fill in secrets
cp .env.example .env
docker compose up -d
```

---

## WebSocket collaboration

The collaboration endpoint at `/ws/collab/<roomId>` requires a valid JWT passed as the `token` query parameter. Every connection is authorized against the resource identified by the room ID:

| Room prefix | Resource checked |
|-------------|-----------------|
| `document:{id}` | Document permissions |
| `drawing:{id}` | Drawing → space permissions |
| `row:{dbId}:{rowId}` | Database → space permissions |

Connections with an invalid token, unknown room format, or insufficient permissions are rejected.

---

## License

MIT
