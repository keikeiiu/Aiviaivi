# Graph Report - /Users/keiiu/Documents/Github/AiliVili  (2026-05-06)

## Corpus Check
- 141 files · ~317,632 words
- Verdict: corpus is large enough that graph structure adds value.

## Summary
- 445 nodes · 751 edges · 46 communities (41 shown, 5 thin omitted)
- Extraction: 60% EXTRACTED · 40% INFERRED · 0% AMBIGUOUS · INFERRED: 297 edges (avg confidence: 0.81)
- Token cost: 0 input · 0 output

## Community Hubs (Navigation)
- [[_COMMUNITY_API Handler Methods|API Handler Methods]]
- [[_COMMUNITY_Project Documentation & Decisions|Project Documentation & Decisions]]
- [[_COMMUNITY_Core Server Infrastructure|Core Server Infrastructure]]
- [[_COMMUNITY_WebSocket & Danmaku Engine|WebSocket & Danmaku Engine]]
- [[_COMMUNITY_Video Pipeline & Transcoding|Video Pipeline & Transcoding]]
- [[_COMMUNITY_HTTP Response & Comments|HTTP Response & Comments]]
- [[_COMMUNITY_JWT Authentication|JWT Authentication]]
- [[_COMMUNITY_Integration Tests & Auth UI|Integration Tests & Auth UI]]
- [[_COMMUNITY_HTTP Middleware & Routing|HTTP Middleware & Routing]]
- [[_COMMUNITY_Configuration & Storage|Configuration & Storage]]
- [[_COMMUNITY_Frontend UI Components|Frontend UI Components]]
- [[_COMMUNITY_Social Features (LikesFollows)|Social Features (Likes/Follows)]]
- [[_COMMUNITY_Frontend App Shell & Navigation|Frontend App Shell & Navigation]]
- [[_COMMUNITY_Comments Data Model|Comments Data Model]]
- [[_COMMUNITY_Creator Analytics|Creator Analytics]]
- [[_COMMUNITY_Auth Handler Types|Auth Handler Types]]
- [[_COMMUNITY_Playlist Handler Types|Playlist Handler Types]]
- [[_COMMUNITY_Category Data Model|Category Data Model]]
- [[_COMMUNITY_Danmaku Handler Types|Danmaku Handler Types]]
- [[_COMMUNITY_Watch Handler Types|Watch Handler Types]]
- [[_COMMUNITY_File Storage Interface|File Storage Interface]]
- [[_COMMUNITY_Environment Configuration|Environment Configuration]]
- [[_COMMUNITY_CORS Configuration|CORS Configuration]]

## God Nodes (most connected - your core abstractions)
1. `Error()` - 46 edges
2. `Handler` - 45 edges
3. `OK()` - 38 edges
4. `UserIDFromContext()` - 31 edges
5. `setupIntegrationDB()` - 24 edges
6. `newTestRouter()` - 22 edges
7. `WithPagination()` - 11 edges
8. `ParseToken()` - 11 edges
9. `Hub` - 11 edges
10. `NewToken()` - 10 edges

## Surprising Connections (you probably didn't know these)
- `LocalStore Disk Implementation` --shares_data_with--> `nginx HLS Static File Serving`  [AMBIGUOUS]
  PRODUCTION.md → README.md
- `Redis Cache` --conceptually_related_to--> `Redis Pub/Sub Cross-Instance Communication`  [INFERRED]
  README.md → .claude/PROGRESS.md
- `Adaptive Bitrate HLS Transcoding (4 qualities)` --shares_data_with--> `nginx HLS Static File Serving`  [INFERRED]
  .claude/PROGRESS.md → README.md
- `MinioStore Object Storage Implementation` --conceptually_related_to--> `nginx HLS Static File Serving`  [AMBIGUOUS]
  PRODUCTION.md → README.md
- `TestUserIDFromContextNotSet()` --calls--> `UserIDFromContext()`  [INFERRED]
  backend/internal/middleware/auth_test.go → backend/internal/middleware/auth.go

## Communities (46 total, 5 thin omitted)

### Community 0 - "API Handler Methods"
Cohesion: 0.12
Nodes (17): Handler, UserIDFromContext(), RequireRole(), Playlist, AddVideoToPlaylist(), CreatePlaylist(), DeletePlaylist(), GetPlaylistByID() (+9 more)

### Community 1 - "Project Documentation & Decisions"
Cohesion: 0.07
Nodes (38): Undecided: Deployment Target (fly.io vs railway vs AWS), Resolved: DB-Backed Refresh Tokens with Rotation, Resolved: MinIO SDK Storage Integration, HLS Configuration (HLS_BASE_URL), JWT Configuration (JWT_SECRET, JWT_EXPIRES_MINUTES), MinIO Configuration (MINIO_ENDPOINT, keys, bucket, SSL), Redis Config (REDIS_URL), Storage Backend Config (STORAGE=local or minio) (+30 more)

### Community 2 - "Core Server Infrastructure"
Cohesion: 0.07
Nodes (13): Open(), ApplyMigrations(), ensureMigrationsTable(), hasMigration(), main(), NewLocalStore(), LocalStore, Hub (+5 more)

### Community 3 - "WebSocket & Danmaku Engine"
Cohesion: 0.08
Nodes (16): Deps, New(), DecWS(), IncDanmaku(), IncTranscoded(), IncUpload(), IncWS(), Danmaku (+8 more)

### Community 4 - "Video Pipeline & Transcoding"
Cohesion: 0.07
Nodes (20): Config, PrometheusMetrics(), statusRecorder, Video, CreateVideo(), CreateVideoQuality(), GetVideoQualities(), IncrementViewCount() (+12 more)

### Community 5 - "HTTP Response & Comments"
Cohesion: 0.12
Nodes (18): parseInt(), parseIntQuery(), commentCreateReq, pageOrDefault(), sizeOrDefault(), ListVideos(), GetWatchHistory(), RecordWatch() (+10 more)

### Community 6 - "JWT Authentication"
Cohesion: 0.12
Nodes (19): NewRefreshToken(), NewToken(), ParseToken(), randomID(), TestNewRefreshToken(), TestNewToken(), TestNewTokenExpired(), TestParseTokenEmpty() (+11 more)

### Community 7 - "Integration Tests & Auth UI"
Cohesion: 0.21
Nodes (26): handleLogin(), authRequest(), login(), newTestRouter(), registerAndLogin(), setupIntegrationDB(), TestIntegrationAnalytics(), TestIntegrationAuthValidation() (+18 more)

### Community 8 - "HTTP Middleware & Routing"
Cohesion: 0.12
Nodes (16): Deps, New(), parseBearer(), RequireAuth(), TestParseBearer(), TestRequireAuthInvalidToken(), TestRequireAuthNoHeader(), TestRequireAuthValidToken() (+8 more)

### Community 9 - "Configuration & Storage"
Cohesion: 0.15
Nodes (12): getenv(), Load(), parsePort(), setenv(), TestLoadCustomPort(), TestLoadDefaults(), TestLoadInvalidPort(), TestLoadMissingDB() (+4 more)

### Community 10 - "Frontend UI Components"
Cohesion: 0.12
Nodes (6): useDanmaku(), useVideo(), formatCount(), formatDuration(), formatTimeAgo(), pad()

### Community 11 - "Social Features (Likes/Follows)"
Cohesion: 0.13
Nodes (12): FavoriteVideo(), FollowUser(), GetFollowerCount(), GetFollowingCount(), IsFollowing(), LikeVideo(), ListFavorites(), SearchUsers() (+4 more)

### Community 12 - "Frontend App Shell & Navigation"
Cohesion: 0.17
Nodes (4): useAuth(), useFeed(), search(), trending()

### Community 15 - "Comments Data Model"
Cohesion: 0.33
Nodes (6): Comment, CreateComment(), DeleteComment(), getReplies(), ListComments(), CommentListParams

### Community 16 - "Creator Analytics"
Cohesion: 0.29
Nodes (6): GetCreatorOverview(), GetCreatorVideoStats(), CreatorOverview, DailyMetric, VideoStats, VideoStatsDetail

### Community 17 - "Auth Handler Types"
Cohesion: 0.5
Nodes (3): loginRequest, refreshRequest, registerRequest

### Community 18 - "Playlist Handler Types"
Cohesion: 0.5
Nodes (3): addVideoReq, playlistCreateReq, playlistUpdateReq

### Community 19 - "Category Data Model"
Cohesion: 0.5
Nodes (3): Category, GetCategoryByID(), ListCategories()

## Ambiguous Edges - Review These
- `nginx HLS Static File Serving` → `MinioStore Object Storage Implementation`  [AMBIGUOUS]
  PRODUCTION.md · relation: conceptually_related_to
- `nginx HLS Static File Serving` → `LocalStore Disk Implementation`  [AMBIGUOUS]
  PRODUCTION.md · relation: shares_data_with

## Knowledge Gaps
- **54 isolated node(s):** `Deps`, `commentCreateReq`, `registerRequest`, `loginRequest`, `refreshRequest` (+49 more)
  These have ≤1 connection - possible missing edges or undocumented components.
- **5 thin communities (<3 nodes) omitted from report** — run `graphify query` to explore isolated nodes.

## Suggested Questions
_Questions this graph is uniquely positioned to answer:_

- **What is the exact relationship between `nginx HLS Static File Serving` and `MinioStore Object Storage Implementation`?**
  _Edge tagged AMBIGUOUS (relation: conceptually_related_to) - confidence is low._
- **What is the exact relationship between `nginx HLS Static File Serving` and `LocalStore Disk Implementation`?**
  _Edge tagged AMBIGUOUS (relation: shares_data_with) - confidence is low._
- **Why does `Error()` connect `API Handler Methods` to `HTTP Middleware & Routing`, `Video Pipeline & Transcoding`, `HTTP Response & Comments`, `JWT Authentication`?**
  _High betweenness centrality (0.150) - this node is a cross-community bridge._
- **Why does `ApplyMigrations()` connect `Core Server Infrastructure` to `Integration Tests & Auth UI`?**
  _High betweenness centrality (0.149) - this node is a cross-community bridge._
- **Why does `Handler` connect `API Handler Methods` to `WebSocket & Danmaku Engine`, `Video Pipeline & Transcoding`, `HTTP Response & Comments`, `JWT Authentication`?**
  _High betweenness centrality (0.120) - this node is a cross-community bridge._
- **Are the 44 inferred relationships involving `Error()` (e.g. with `.CommentList()` and `.CommentCreate()`) actually correct?**
  _`Error()` has 44 INFERRED edges - model-reasoned connections that need verification._
- **Are the 36 inferred relationships involving `OK()` (e.g. with `.CommentCreate()` and `.CommentDelete()`) actually correct?**
  _`OK()` has 36 INFERRED edges - model-reasoned connections that need verification._