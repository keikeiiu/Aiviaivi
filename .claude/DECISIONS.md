# Undecided / Needs Decision

> Only 3 architectural decisions remain — all are production concerns, not blockers for the prototype.

| # | Category | Question | Why it matters |
|---|----------|----------|----------------|
| 1 | Refresh Tokens | Stateless JWT (current) vs DB-backed for revocation? | DB-backed lets you force-logout users. Needed before launch. |
| 2 | MinIO SDK | Add real `minio-go` implementation vs keep local storage? | Local `/uploads/` works fine for dev. MinIO needed at scale. |
| 3 | Deploy Target | Docker Compose (current) → fly.io, railway, AWS? | `docker compose up` + `npx expo start --web` both verified working. |
