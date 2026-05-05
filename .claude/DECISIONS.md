# Undecided / Needs Decision

| # | Category | Question | Context |
|---|----------|----------|---------|
| 1 | Refresh Tokens | Stateless JWT (current) vs DB-backed (revocable per plan)? | Both approaches valid. Stateless is simpler. |
| 2 | HLS Base URL | manifest_urls are relative paths. Frontend `HLS_BASE_URL` is localhost:8081. Make configurable via env? | Works for dev. Production needs env var. |
| 3 | MinIO SDK | Storage abstraction exists (interface + local + MinIO stub). Add real `minio-go`? | Local storage verified. MinIO for production. |
| 4 | App Store Deploy | Full stack verified. Ready for: (a) Docker deploy, (b) iOS/Android build, (c) cloud deploy. | `docker compose up` + `npx expo start --web` both verified. |
