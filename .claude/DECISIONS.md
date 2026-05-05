# Decisions — All Resolved

> All 3 architectural decisions are resolved. 52/52 plan items delivered (100%).

| # | Category | Resolution |
|---|----------|------------|
| 1 | Refresh Tokens | DB-backed with rotation, jti uniqueness, token reuse detection |
| 2 | Storage Backend | MinIO SDK integrated (`STORAGE=minio`), local default |
| 3 | Deploy Target | **fly.io** — `fly.toml` + `scripts/deploy.sh` ready. One command: `./scripts/deploy.sh` |

No undecided items remain.
