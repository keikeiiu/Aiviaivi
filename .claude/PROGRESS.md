# AiliVili Build Progress

## Status: Backend Fully Verified — P0–P3 Live Tested

**99 source files | ~5,800 lines | 12 migrations | 40 endpoints | 6 app screens | 47 unit tests | 35 live tests**

---

## Live Test Results (2026-05-06)

### Infrastructure
- PostgreSQL 16.13 — installed, running, 12 tables verified
- Redis 8.6 — installed, running, pub/sub + rate limiting active
- FFmpeg 8.1 — installed, HLS transcode verified (4 qualities + thumbnail)
- Backend running at localhost:8080 with all 12 migrations applied

### Full Pipeline Verified
1. **Upload → Transcode → Published**: 3s test video → 4 HLS qualities (1080p/720p/480p/360p) + thumbnail → status "published" ✅
2. **HLS segments**: `.m3u8` manifests + `.ts` segment files generated per quality ✅
3. **Thumbnail**: 8.5KB JPEG generated at video 1s mark ✅
4. **File upload**: Raw files saved to `uploads/raw/`, HLS to `uploads/hls/{id}/` ✅

### Endpoints Tested (35 with curl)

| # | Endpoint | Result | # | Endpoint | Result |
|---|----------|--------|---|----------|--------|
| 1 | POST /auth/register | ✅ | 19 | POST /videos/{id}/favorite | ✅ |
| 2 | POST /auth/login | ✅ | 20 | POST /videos/{id}/watch | ✅ |
| 3 | POST /auth/refresh | ✅ | 21 | GET /users/me/history | ✅ |
| 4 | GET /health | ✅ | 22 | GET /feed/trending | ✅ |
| 5 | GET /users/me | ✅ | 23 | GET /search?q=test | ✅ |
| 6 | GET /users/{id} | ✅ | 24 | GET /playlists | ✅ |
| 7 | PUT /users/{id} | ✅ | 25 | POST /playlists | ✅ |
| 8 | GET /users/{id}/videos | ✅ | 26 | POST /playlists/{id}/videos | ✅ |
| 9 | GET /users/{id}/favorites | ✅ | 27 | GET /playlists/{id} | ✅ |
| 10 | POST /users/{id}/subscribe | ✅ | 28 | GET /analytics/overview | ✅ |
| 11 | GET /categories | ✅ | 29 | GET /analytics/videos | ✅ |
| 12 | GET /videos | ✅ | 30 | GET /metrics | ✅ |
| 13 | GET /videos/{id} | ✅ | 31 | Redis rate limiting | ✅ |
| 14 | POST /videos/upload | ✅ | 32 | Redis pub/sub | ✅ |
| 15 | POST /videos/{id}/danmaku | ✅ | 33 | Danmaku counter | ✅ |
| 16 | GET /videos/{id}/danmaku | ✅ | 34 | HLS segments on disk | ✅ |
| 17 | POST /videos/{id}/comments | ✅ | 35 | Thumbnail on disk | ✅ |
| 18 | POST /videos/{id}/like | ✅ | | | |

### Database Contents
```
categories        |   12 (seed)
schema_migrations |   12 (applied)
users             |    1
videos            |    5
danmaku           |    5
comments          |    1
likes             |    1
favorites         |    1
playlists         |    1
playlist_videos   |    1
watch_history     |    1
```

### Bugs Found & Fixed During Live Testing
- **Thumbnail always empty**: seek time 5s > video duration → changed to 1s, moved `-ss` before `-i`
- **REST danmaku counter not incrementing**: added `metrics.IncDanmaku()` to REST handler
- **Expo web blocked**: manual package.json has version mismatch with expo-router → needs `npx create-expo-app` bootstrap

## Unit Tests — 47/47 Pass
```
ok  	ailivili/internal/auth       ✅
ok  	ailivili/internal/config     ✅
ok  	ailivili/internal/httpapi    ✅
ok  	ailivili/internal/middleware  ✅
ok  	ailivili/internal/response   ✅
```
