# Undecided / Needs Decision

| # | Category | Question | Context |
|---|----------|----------|---------|
| 1 | Expo Bootstrap | Frontend source code is complete + TypeScript clean. But Expo needs `npx create-expo-app .` for correct dependency resolution. See README for bootstrap steps. | Manually managing Expo SDK 55 deps hits cascading peer dep issues (react-native-worklets, expo-router/internal/routing). `create-expo-app` fixes this. |
| 2 | Refresh Tokens | Stateless JWT (current) vs DB-backed (revocable per plan)? | Both work. Stateless is simpler. |
| 3 | HLS Base URL | manifest_urls are relative paths. Make configurable via env? | Works for local dev. Production needs env var. |
| 4 | MinIO SDK | Storage abstraction exists (interface + local + MinIO stub). Add real `minio-go`? | Local storage verified working. MinIO for production scale. |
| 5 | Production Deploy | Docker Compose ready. Cloud deploy configs (fly.io, railway, etc.)? | `docker compose up` works. No cloud configs yet. |
