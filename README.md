# AiliVili

Prototype backend + frontend for a Bilibili-like video platform.

## Quick Start

### Backend (Go + PostgreSQL + Redis + FFmpeg)

```bash
# Option A: Docker Compose (recommended)
docker compose up --build

# Option B: Local (requires PostgreSQL, Redis, FFmpeg)
brew install postgresql@16 redis ffmpeg
brew services start postgresql@16 redis
createdb ailivili
DATABASE_URL="postgres://localhost:5432/ailivili?sslmode=disable" \
JWT_SECRET="dev-secret" REDIS_URL="localhost:6379" \
go run ./backend/cmd/server
```

API: `http://localhost:8080/api/v1`
Metrics: `http://localhost:8080/metrics`

### Frontend (React Native / Expo)

**Bootstrap (first time only):**
```bash
cd frontend
npx create-expo-app . --template blank-typescript --force
npx expo install expo-router expo-av expo-document-picker
npm install zustand axios @react-native-async-storage/async-storage
# All source files (app/, components/, hooks/, store/, services/, utils/) are pre-written
```

**Run:**
```bash
cd frontend
npx expo start --web     # or --ios, --android
```

## API — 40 Endpoints

See [.claude/PROGRESS.md](.claude/PROGRESS.md) for full endpoint map and test results.

### Test (curl)

```bash
# Health
curl -s http://localhost:8080/api/v1/health

# Register
curl -s -X POST http://localhost:8080/api/v1/auth/register \
  -H 'Content-Type: application/json' \
  -d '{"username":"u1","email":"u1@test.com","password":"pass"}'

# Login
curl -s -X POST http://localhost:8080/api/v1/auth/login \
  -H 'Content-Type: application/json' \
  -d '{"email":"u1@test.com","password":"pass"}'
```

### Go Tests
```bash
cd backend && go test ./...
```

## Architecture

```
backend/         Go monolith (chi router, PostgreSQL, Redis, FFmpeg)
frontend/        React Native / Expo (expo-router, Zustand, Axios)
nginx/           HLS static file serving
docker-compose.yml  postgres + redis + api + nginx
```

## Config

| Env Var | Default | Description |
|---------|---------|-------------|
| PORT | 8080 | Server port |
| DATABASE_URL | (required) | PostgreSQL DSN |
| JWT_SECRET | (required) | JWT signing key |
| JWT_EXPIRES_MINUTES | 60 | Access token lifetime |
| REDIS_URL | (optional) | Redis address for cache + rate limiting |
| STORAGE | local | Storage backend: "local" or "minio" |
| STORAGE_BASE_URL | (optional) | Public base URL for stored files |

## Deploy (fly.io)

```bash
# One-time setup
brew install flyctl        # macOS
# or: curl -L https://fly.io/install.sh | sh

fly auth signup             # Create account
fly launch                  # Auto-detects config

# Set secrets
fly secrets set DATABASE_URL="postgres://..." JWT_SECRET="$(openssl rand -hex 32)"

# Deploy
fly deploy

# Or use the automated script
./scripts/deploy.sh
```

API will be at `https://ailivili.fly.dev/api/v1/health`
