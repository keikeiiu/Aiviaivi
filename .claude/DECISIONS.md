# Undecided / Needs Decision

> 1 decision remains. #1 (refresh tokens) and #2 (MinIO storage) are resolved.

| # | Category | Question | Why it matters |
|---|----------|----------|----------------|
| 1 | ~~Refresh Tokens~~ | ~~Stateless vs DB-backed~~ | **Resolved**: DB-backed with rotation, jti uniqueness, token reuse detection. |
| 2 | ~~Storage Backend~~ | ~~Local vs MinIO~~ | **Resolved**: MinIO SDK integrated. Set `STORAGE=minio` to switch. Integration tested against real MinIO server. |
| 3 | Deploy Target | Docker Compose → fly.io, railway, AWS? | `docker compose up` + `npx expo start --web` verified. |
