# Undecided / Needs Decision

> 2 decisions remain. #1 was resolved and implemented (DB-backed refresh tokens).

| # | Category | Question | Why it matters |
|---|----------|----------|----------------|
| 1 | ~~Refresh Tokens~~ | ~~Stateless vs DB-backed~~ | **Resolved**: DB-backed implemented with rotation, jti uniqueness, token reuse detection. |
| 2 | Storage Backend | Keep local storage or implement MinIO SDK? | Local `/uploads/` verified working. MinIO needed at scale. |
| 3 | Deploy Target | Docker Compose → fly.io, railway, AWS? | `docker compose up` + `npx expo start --web` verified. |
