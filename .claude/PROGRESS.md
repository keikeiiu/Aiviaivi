# AiliVili Build Progress

## Status: **Delivered** — P0–P3 Complete, Tested, Live Verified

**125 files | ~7,700 lines | 13 migrations | 40 endpoints | 11 app screens | 8 components | 5 hooks | 67 + 3 Playwright tests | 16 bugs fixed | 1 knowledge graph (445 nodes, 751 edges) | 3/3 decisions resolved**

---

## Plan Compliance: 52/52 Items — 100% Complete

| Phase | Items | Delivered | Status |
|-------|-------|-----------|--------|
| P0 MVP | 15 | 15 | ✅ Complete |
| P1 Core | 11 | 11 | ✅ Complete |
| P2 Real-time | 12 | 12 | ✅ Complete |
| P3 Polish | 14 | 14 | ✅ Complete |
| **Total** | **52** | **52** | **100%** |

---

## Recent Updates (2026-05-06 last session)

### Bugs Fixed This Session (3)
15. **Upload silently failed — FormData sent as JSON**: Axios default `transformRequest` serialized FormData with `JSON.stringify()`, converting it to `{}`. Added `transformRequest: [(data) => data]` to prevent serialization.
16. **Upload Content-Type missing boundary**: Axios default `Content-Type: application/json` overrode browser's `multipart/form-data; boundary=...`. Set explicit `Content-Type: multipart/form-data` per-request.
17. **HLS files not served locally**: No nginx in local dev → Go backend couldn't serve HLS segments. Added `/uploads/*` file server route in httpapi.go, changed `HLS_BASE_URL` to port 8080.

### Playwright E2E Tests Added
- `frontend/e2e/upload.spec.ts` — 3 tests: API connectivity, registration, upload with FormData
- `frontend/playwright.config.ts` — Chromium headless, 30s timeout
- All tests pass against live backend + frontend

### Video Playback Verified
- HLS manifest: HTTP 200 at `http://localhost:8080/uploads/hls/<id>/720p/index.m3u8`
- expo-av VideoPlayer loads and plays HLS streams
- Tested with 48-second real video (7.8MB upload → 13 segments)

---

## Full Bug List (16 Total)

1. `cover_url`/`avatar_url` NULL → `NOT NULL DEFAULT ''`
2. `ANY($1)` → `ANY($1::uuid[])` UUID array cast
3. Thumbnail `-ss 5s` exceeded short video duration → `-ss 1s` before `-i`
4. `AuthRefresh` masked DB errors as "not found" → check `ErrUserNotFound`
5. WebSocket route killed by 30s timeout → route outside `Timeout` middleware
6. REST danmaku counter not incrementing → added `metrics.IncDanmaku()`
7. Duplicate `Playlist` type → removed from social.go
8. Unused imports (`io`, `os`, `time`) → cleaned
9. `UserIDFromContext(nil)` panic → use `context.Background()`
10. Expo dependency graph unresolvable → `create-expo-app` bootstrap
11. `async-storage` v2/v3 API mismatch → v2 API (`multiSet`/`multiGet`/`multiRemove`)
12. WebSocket upgrade blocked by `statusRecorder` → added `Hijack()`/`Flush()`/`Push()`/`Unwrap()`
13. `time.Now().Unix()` → identical JWTs within same second → `jti` claim with `crypto/rand`
14. `react-native-worklets` peer dep missing → installed `react-native-worklets@0.5.1`
15. **Upload: FormData serialized as JSON** → added `transformRequest: [(data) => data]`
16. **Upload: Content-Type missing boundary** → set `Content-Type: multipart/form-data` per-request
17. **HLS not served in local dev** → added `/uploads/*` file server, changed HLS_BASE_URL to 8080

---

## Test Results

| Layer | Count | Status |
|-------|-------|--------|
| Go unit tests | 47 | ✅ All pass |
| Go integration tests (real DB) | 20 | ✅ All pass |
| MinIO integration test (real server) | 1 | ✅ Pass |
| Playwright E2E tests | 3 | ✅ All pass |
| TypeScript strict | — | ✅ 0 errors |
| go vet | — | ✅ 0 issues |
| Live endpoint tests (curl) | 35 | ✅ All pass |
| WebSocket danmaku | — | ✅ Round-trip verified |
| E2E flow (register→upload→transcode→danmaku→comment→like) | — | ✅ Verified |
| Token rotation + reuse rejection | — | ✅ Verified |
| Video playback (HLS) | — | ✅ Verified locally |

---

## All Decisions Resolved

| # | Decision | Resolution |
|---|----------|------------|
| 1 | Refresh Tokens | DB-backed with rotation, jti uniqueness |
| 2 | Storage Backend | MinIO SDK integrated, local default |
| 3 | Deploy Target | fly.io — `fly.toml` + `scripts/deploy.sh` ready |

---

## How to Run

```bash
# Full stack (Docker)
docker compose up --build

# Local dev
DATABASE_URL="postgres://..." JWT_SECRET="..." go run ./backend/cmd/server
cd frontend && npx expo start --web

# Tests
cd backend && go test ./...
DATABASE_URL="..." go test -tags=integration ./internal/httpapi/ -v
npx playwright test --config frontend/playwright.config.ts
```
