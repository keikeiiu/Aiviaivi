# AiliVili

Prototype backend for a Bilibili-like video platform.

## What’s Included

- Go API (chi router)
- PostgreSQL
- JWT auth (Bearer token)
- Auto-applied SQL migrations on startup

## Run (Docker)

```bash
docker compose up --build
```

API: `http://localhost:8080/api/v1`

## Try It (curl)

Health:

```bash
curl -s http://localhost:8080/api/v1/health
```

Register:

```bash
curl -s -X POST http://localhost:8080/api/v1/auth/register \
  -H 'Content-Type: application/json' \
  -d '{"username":"u1","email":"u1@test.com","password":"pass"}'
```

Login:

```bash
curl -s -X POST http://localhost:8080/api/v1/auth/login \
  -H 'Content-Type: application/json' \
  -d '{"email":"u1@test.com","password":"pass"}'
```

Me:

```bash
TOKEN="$(curl -s -X POST http://localhost:8080/api/v1/auth/login \
  -H 'Content-Type: application/json' \
  -d '{"email":"u1@test.com","password":"pass"}' | python3 -c 'import json,sys; print(json.load(sys.stdin)["data"]["token"])')"

curl -s http://localhost:8080/api/v1/users/me \
  -H "Authorization: Bearer ${TOKEN}"
```

## Config

The server reads:

- `PORT` (default: 8080)
- `DATABASE_URL` (required)
- `JWT_SECRET` (required)
- `JWT_EXPIRES_MINUTES` (default: 60)
- `MIGRATIONS_DIR` (default: `migrations`)
