# Undecided / Needs Decision

| # | Category | Question | Context |
|---|----------|----------|---------|
| 1 | Expo Bootstrap | Frontend code is complete but can't start due to version mismatch between manually-configured package.json and expo-router. Fix: `cd frontend && npx create-expo-app .` to regenerate proper deps, then re-apply source files. | Frontend code is correct (tsc passes). Package versions need Expo's auto-resolution. |
| 2 | Refresh Tokens | Stateless JWT (current) vs DB-backed (revocable per plan)? | Both approaches are valid. Current stateless works. |
| 3 | HLS Base URL | manifest_urls are relative paths. Frontend `HLS_BASE_URL` is localhost:8081. Make configurable? | Works for dev. Needs env var for production. |
| 4 | MinIO SDK | Storage abstraction exists (interface + local + MinIO stub). Add real `minio-go`? | Local storage verified working. MinIO for production scale. |
| 5 | Deploy Configs | Docker Compose ready. Add production deploy scripts (fly.io, railway, etc.)? | `docker compose up` works. No cloud deploy configs yet. |
