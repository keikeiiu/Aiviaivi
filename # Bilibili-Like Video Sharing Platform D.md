# Bilibili-Like Video Sharing Platform Design Document

## 1. Overview
This document outlines the technical architecture and feature set for a Bilibili-inspired video-sharing platform. The system focuses on high-concurrency video streaming, real-time community interaction (Danmaku), and personalized content discovery.

## 2. Core Features

### 2.1 User Experience
- **Video Playback**: Adaptive bitrate streaming (HLS/DASH) for smooth playback across devices.
- **Danmaku (Bullet Comments)**: Real-time, overlay scrolling comments synchronized with video timestamps.
- **Content Discovery**: Personalized feed, categories, search, and trending topics.
- **Community**: Comments, likes, shares, favorites, and user profiles.

### 2.2 Creator Tools
- **Creator Studio**: Interface for uploading, managing, and analyzing videos.
- **Analytics**: View counts, watch time, and engagement metrics.

## 3. System Architecture

### 3.1 Frontend (Client-Side)
- **Web App**: Built with React.js or Vue.js for a responsive single-page application (SPA).
- **Mobile App**: Built with Flutter or React Native for cross-platform iOS/Android support.
- **Video Player**: Custom player supporting HLS/DASH protocols with a Danmaku overlay engine (using HTML5 Canvas or SVG).
- **Key Components**:
  - Infinite scroll feed for content consumption.
  - Real-time WebSocket client for Danmaku and notifications.

### 3.2 Backend (Microservices)
The backend is divided into specialized microservices:

1. **API Gateway**
   - **Role**: Entry point for all client requests. Handles routing, rate limiting, and authentication (JWT).
   - **Tech**: Kong, NGINX, or Envoy.

2. **User Service**
   - **Role**: Manages user registration, login, profiles, and subscription relationships.
   - **Tech**: Node.js or Go.

3. **Video Service**
   - **Upload Manager**: Handles multipart file uploads from creators.
   - **Transcoding Engine**: Converts raw uploads into multiple resolutions (1080p, 720p, 480p) and formats (HLS/DASH).
   - **Metadata Service**: Stores video titles, descriptions, tags, and thumbnails.
   - **Tech**: Go (for upload handling), Python/FFmpeg (for transcoding).

4. **Danmaku Service**
   - **Role**: Manages real-time bullet comments.
   - **Flow**: Receives comments via WebSocket, validates/sanitizes input, persists data, and broadcasts to relevant video channels.
   - **Tech**: Go or Node.js (high concurrency), Redis (pub/sub), MongoDB.

5. **Recommendation Service**
   - **Role**: Generates personalized video feeds based on user behavior (watch history, likes, clicks).
   - **Tech**: Python (TensorFlow/PyTorch for ML models), Collaborative Filtering algorithms.

6. **Search Service**
   - **Role**: Full-text search for videos, users, and tags.
   - **Tech**: Elasticsearch.

7. **Interaction Service**
   - **Role**: Manages likes, dislikes, shares, and comments (non-Danmaku).
   - **Tech**: Node.js/Go.

### 3.3 Infrastructure & Storage

- **Object Storage**: AWS S3, Alibaba Cloud OSS, or MinIO (self-hosted) for storing video files, thumbnails, and assets.
- **CDN (Content Delivery Network)**: Cloudflare or Alibaba Cloud CDN to cache video segments globally and reduce latency.
- **Database**:
  - **PostgreSQL**: Relational data (users, video metadata, transactions).
  - **MongoDB**: Unstructured data (Danmaku logs, detailed comments).
  - **Redis**: Caching (hot videos, sessions), leaderboards, and WebSocket pub/sub.
- **Message Queue**: Kafka or RabbitMQ for async tasks (transcoding jobs, analytics aggregation, notifications).
- **Containerization**: Docker and Kubernetes for orchestration and scaling.

## 4. Data Flow Examples

### 4.1 Video Upload Flow
1. Creator uploads raw video file via Frontend.
2. API Gateway authenticates user and routes to **Video Service**.
3. **Video Service** uploads raw file to **Object Storage**.
4. **Video Service** publishes a message to **Kafka** with video metadata.
5. **Transcoding Engine** consumes Kafka message, fetches video, and processes it into multiple resolutions.
6. Transcoded segments are uploaded back to **Object Storage**.
7. **Video Service** updates metadata (available resolutions) in **PostgreSQL**.

### 4.2 Video Playback & Danmaku Flow
1. User requests video details.
2. **API Gateway** checks **Redis** cache. If miss, queries **PostgreSQL**.
3. User starts playback. Frontend connects to **CDN** for video segments (HLS/DASH).
4. User sends a Danmaku comment via WebSocket.
5. **Danmaku Service** receives comment, validates it, saves to **MongoDB**.
6. **Danmaku Service** broadcasts the comment to all other clients subscribed to this video ID via **Redis Pub/Sub**.

## 5. Technology Stack Summary

| Component          | Recommended Technology       |
|--------------------|------------------------------|
| **Frontend**       | React.js / Vue.js, Flutter  |
| **Backend**        | Go (High Concurrency), Node.js, Python (ML) |
| **Database**       | PostgreSQL, MongoDB, Redis   |
| **Search**         | Elasticsearch                |
| **Object Storage** | AWS S3 / MinIO               |
| **Streaming**      | HLS / DASH, Nginx-RTMP       |
| **Infrastructure** | Docker, Kubernetes, Kafka    |

## 6. Next Steps (Learning/Prototype Edition — ~10 concurrent users)

> For a small-scale prototype, the microservices + Kafka + Kubernetes stack is overkill. Below is a simplified, practical plan using a monolith backend + 2 databases + Docker Compose.

---

### 6.1 Simplified Architecture

```
┌─────────────────────────────────────────────┐
│  Frontend (Mobile-first)                     │
│  React Native 0.83 + Expo SDK 55             │
│  expo-router v4 · Zustand · Axios            │
│  react-native-video · Custom Canvas Danmaku │
└──────────────────┬──────────────────────────┘
                   │ HTTP/WS
┌──────────────────▼──────────────────────────┐
│  Backend (Monolith)                          │
│  Go (net/http + chi/gin router)              │
│  JWT auth · REST API · File upload           │
│  FFmpeg subprocess for transcoding           │
│  Serve HLS segments directly or via nginx    │
└──────────────────┬──────────────────────────┘
                   │
      ┌────────────┼────────────┐
      ▼            ▼            ▼
┌──────────┐ ┌──────────┐ ┌──────────┐
│PostgreSQL│ │  Local   │ │  Redis   │
│  (data)  │ │  Disk    │ │(optional)│
└──────────┘ └──────────┘ └──────────┘

Deploy: Docker Compose (3 containers: app + postgres + redis)
```

| Layer | Choice | Rationale |
|-------|--------|-----------|
| **Backend** | Single Go server (monolith) | One binary, handles thousands of concurrent connections, no microservice overhead |
| **Database** | PostgreSQL only | Relational data for everything — even danmaku; no need for MongoDB at this scale |
| **Cache** | Redis (optional, defer to P2) | Rate limiting + session cache; skip for P0 |
| **File Storage** | Local `/uploads/` directory first, then MinIO | Start simple, add object storage in P3 |
| **Search** | PostgreSQL `tsvector` / `ILIKE` | Built-in full-text search; add Elasticsearch only if needed |
| **Deploy** | Docker Compose (3 containers) | No Kubernetes required |

---

### 6.2 Database Schema (PostgreSQL DDL)

```
┌──────────┐     ┌──────────┐     ┌──────────────┐
│  users   │────→│  videos  │────→│video_quality │
└──────────┘     └──────────┘     └──────────────┘
     │                │
     │    ┌───────────┼───────────┐───────────┐
     ▼    ▼           ▼           ▼           ▼
┌──────┐ ┌────────┐ ┌────────┐ ┌──────────┐ ┌──────────┐
│follow│ │danmaku │ │comments│ │  likes   │ │favorites │
└──────┘ └────────┘ └────────┘ └──────────┘ └──────────┘
```

```sql
-- Users
CREATE TABLE users (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    username VARCHAR(50) UNIQUE NOT NULL,
    email VARCHAR(255) UNIQUE NOT NULL,
    password_hash VARCHAR(255) NOT NULL,
    avatar_url VARCHAR(500),
    bio TEXT DEFAULT '',
    role VARCHAR(10) DEFAULT 'user' CHECK (role IN ('user','creator','admin')),
    created_at TIMESTAMPTZ DEFAULT now(),
    updated_at TIMESTAMPTZ DEFAULT now()
);

-- Categories
CREATE TABLE categories (
    id SERIAL PRIMARY KEY,
    name VARCHAR(50) NOT NULL,
    slug VARCHAR(50) UNIQUE NOT NULL
);

-- Videos
CREATE TABLE videos (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    title VARCHAR(200) NOT NULL,
    description TEXT DEFAULT '',
    cover_url VARCHAR(500),
    duration INTEGER DEFAULT 0,
    status VARCHAR(20) DEFAULT 'uploading'
        CHECK (status IN ('uploading','processing','published','private','deleted')),
    category_id INT REFERENCES categories(id),
    tags TEXT[] DEFAULT '{}',
    view_count BIGINT DEFAULT 0,
    like_count INT DEFAULT 0,
    comment_count INT DEFAULT 0,
    share_count INT DEFAULT 0,
    published_at TIMESTAMPTZ,
    created_at TIMESTAMPTZ DEFAULT now(),
    updated_at TIMESTAMPTZ DEFAULT now()
);

-- Video quality variants (HLS manifests)
CREATE TABLE video_qualities (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    video_id UUID NOT NULL REFERENCES videos(id) ON DELETE CASCADE,
    quality VARCHAR(10) NOT NULL,  -- '1080p','720p','480p','360p'
    manifest_url VARCHAR(500) NOT NULL,
    bitrate INT,
    file_size BIGINT,
    created_at TIMESTAMPTZ DEFAULT now()
);

-- Follows (creator subscription)
CREATE TABLE follows (
    follower_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    creator_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    created_at TIMESTAMPTZ DEFAULT now(),
    PRIMARY KEY (follower_id, creator_id)
);

-- Danmaku (bullet comments) — REST polling first, WebSocket later
CREATE TABLE danmaku (
    id BIGSERIAL PRIMARY KEY,
    video_id UUID NOT NULL REFERENCES videos(id) ON DELETE CASCADE,
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    content TEXT NOT NULL,
    video_time FLOAT NOT NULL,      -- seconds into video
    color VARCHAR(7) DEFAULT '#FFFFFF',
    font_size VARCHAR(10) DEFAULT 'medium'
        CHECK (font_size IN ('small','medium','large')),
    mode VARCHAR(10) DEFAULT 'scroll'
        CHECK (mode IN ('scroll','top','bottom')),
    created_at TIMESTAMPTZ DEFAULT now()
);
CREATE INDEX idx_danmaku_video_time ON danmaku(video_id, video_time);

-- Comments (thread-based, non-real-time)
CREATE TABLE comments (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    video_id UUID NOT NULL REFERENCES videos(id) ON DELETE CASCADE,
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    parent_id UUID REFERENCES comments(id) ON DELETE CASCADE,
    content TEXT NOT NULL,
    like_count INT DEFAULT 0,
    created_at TIMESTAMPTZ DEFAULT now(),
    updated_at TIMESTAMPTZ DEFAULT now()
);

-- Likes (unique per user per video)
CREATE TABLE likes (
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    video_id UUID NOT NULL REFERENCES videos(id) ON DELETE CASCADE,
    created_at TIMESTAMPTZ DEFAULT now(),
    PRIMARY KEY (user_id, video_id)
);

-- Favorites
CREATE TABLE favorites (
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    video_id UUID NOT NULL REFERENCES videos(id) ON DELETE CASCADE,
    created_at TIMESTAMPTZ DEFAULT now(),
    PRIMARY KEY (user_id, video_id)
);

-- Playlists
CREATE TABLE playlists (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    name VARCHAR(100) NOT NULL,
    description TEXT DEFAULT '',
    is_public BOOLEAN DEFAULT false,
    created_at TIMESTAMPTZ DEFAULT now(),
    updated_at TIMESTAMPTZ DEFAULT now()
);

CREATE TABLE playlist_videos (
    playlist_id UUID NOT NULL REFERENCES playlists(id) ON DELETE CASCADE,
    video_id UUID NOT NULL REFERENCES videos(id) ON DELETE CASCADE,
    position INT NOT NULL,
    added_at TIMESTAMPTZ DEFAULT now(),
    PRIMARY KEY (playlist_id, video_id)
);
```

---

### 6.3 API Specification

**Base URL:** `http://localhost:8080/api/v1`

**Response envelope (all endpoints):**
```json
{
  "code": 0,
  "message": "ok",
  "data": { ... },
  "pagination": { "page": 1, "size": 20, "total": 150 }
}
```

#### 6.3.1 Authentication

| Method | Endpoint | Body / Params | Auth | Description |
|--------|----------|---------------|------|-------------|
| POST | `/auth/register` | `{ username, email, password }` | No | Create account |
| POST | `/auth/login` | `{ email, password }` | No | Returns `{ token, refresh_token }` |
| POST | `/auth/refresh` | `{ refresh_token }` | No | Returns new `{ token }` |

#### 6.3.2 Users

| Method | Endpoint | Query | Auth | Description |
|--------|----------|-------|------|-------------|
| GET | `/users/:id` | — | No | Public profile |
| PUT | `/users/:id` | — | Yes | Update profile (bio, avatar) |
| GET | `/users/:id/videos` | `?page=&size=` | No | User's uploaded videos |
| GET | `/users/:id/favorites` | `?page=&size=` | No | User's favorited videos |
| POST | `/users/:id/subscribe` | — | Yes | Follow this user |
| DELETE | `/users/:id/subscribe` | — | Yes | Unfollow |

#### 6.3.3 Videos

| Method | Endpoint | Query / Body | Auth | Description |
|--------|----------|--------------|------|-------------|
| GET | `/videos` | `?page=&size=&category=&sort=latest\|trending` | No | Feed (paginated) |
| GET | `/videos/:id` | — | No | Detail + available qualities |
| POST | `/videos/upload` | Multipart: `file` + `{ title, description, category_id, tags }` | Creator | Upload raw video |
| PUT | `/videos/:id` | `{ title?, description?, tags? }` | Owner | Update metadata |
| DELETE | `/videos/:id` | — | Owner | Soft delete |
| GET | `/videos/:id/related` | — | No | Related videos (same category) |
| GET | `/videos/:id/hls/:quality/index.m3u8` | — | No | Serve HLS manifest + segments |

#### 6.3.4 Danmaku (Bullet Comments)

| Method | Endpoint | Query / Body | Auth | Description |
|--------|----------|--------------|------|-------------|
| GET | `/videos/:id/danmaku` | `?t_start=&t_end=` | No | Fetch danmaku by video time range |
| POST | `/videos/:id/danmaku` | `{ content, video_time, color?, mode?, font_size? }` | Yes | Send a danmaku |

**P0 polling strategy:** Frontend polls `GET /videos/:id/danmaku?t_start=<currentTime-2>&t_end=<currentTime+5>` every 2 seconds. Switch to WebSocket in P2.

#### 6.3.5 Comments

| Method | Endpoint | Query / Body | Auth | Description |
|--------|----------|--------------|------|-------------|
| GET | `/videos/:id/comments` | `?page=&size=&sort=latest\|hot` | No | Paginated comments |
| POST | `/videos/:id/comments` | `{ content, parent_id? }` | Yes | Post comment / reply |
| DELETE | `/comments/:id` | — | Owner | Delete comment |

#### 6.3.6 Social / Interactions

| Method | Endpoint | Auth | Description |
|--------|----------|------|-------------|
| POST | `/videos/:id/like` | Yes | Like video |
| DELETE | `/videos/:id/like` | Yes | Unlike |
| POST | `/videos/:id/favorite` | Yes | Add to favorites |
| DELETE | `/videos/:id/favorite` | Yes | Remove from favorites |

#### 6.3.7 Search & Discovery

| Method | Endpoint | Query | Auth | Description |
|--------|----------|-------|------|-------------|
| GET | `/search` | `?q=&type=video\|user&page=&size=` | No | Full-text search |
| GET | `/categories` | — | No | List all categories |
| GET | `/feed/trending` | `?page=&size=` | No | Top videos by view count |

---

### 6.4 Tech Stack Summary (Simplified)

| Layer | Technology | Notes |
|-------|-----------|-------|
| **Frontend Framework** | React Native 0.83 + Expo SDK 55 | Cross-platform iOS/Android/Web |
| **Frontend Routing** | expo-router v4 | File-system based routing (model after JKVideo) |
| **Frontend State** | Zustand | Lightweight store (JKVideo-proven pattern) |
| **Frontend Network** | Axios + WebSocket | REST polling for danmaku in P0, WS in P2 |
| **Video Player** | react-native-video | DASH/HLS native decode |
| **Danmaku Renderer** | Custom Canvas overlay | Synchronized scrolling bullet comments |
| **Backend Language** | Go (1.22+) | `net/http` + `chi` router |
| **Backend Auth** | JWT (access + refresh tokens) | BCrypt password hashing |
| **Backend File Upload** | Multipart streaming | Chunked upload for large files |
| **Video Transcoding** | FFmpeg (subprocess) | Raw → HLS (1080p/720p/480p/360p) |
| **Database** | PostgreSQL 16 | All data including danmaku |
| **Cache** | Redis 7 (optional, P2) | Session cache, rate limiting, WS pub/sub |
| **File Storage** | Local disk → MinIO (P3) | `/uploads/` directory mounted as volume |
| **Search** | PostgreSQL `tsvector` | Built-in full-text search |
| **Containerization** | Docker Compose | 3 containers: `app` + `postgres` + `redis` |
| **HLS Serving** | Nginx (static files) | Serves `/uploads/hls/` as static content |
| **CI/CD** | GitHub Actions | Lint, test, build Docker image |

---

### 6.5 Project Directory Structure

```
ailivili/
├── backend/                          # Go monolith
│   ├── cmd/
│   │   └── server/
│   │       └── main.go               # Entry point
│   ├── internal/
│   │   ├── config/
│   │   │   └── config.go             # Env config (DB URL, JWT secret, etc.)
│   │   ├── handler/
│   │   │   ├── auth.go               # POST /auth/register, /login, /refresh
│   │   │   ├── user.go               # GET/PUT /users/:id
│   │   │   ├── video.go              # CRUD /videos, upload
│   │   │   ├── danmaku.go            # GET/POST /videos/:id/danmaku
│   │   │   ├── comment.go            # GET/POST /videos/:id/comments
│   │   │   ├── social.go             # Like, favorite, subscribe
│   │   │   ├── search.go             # GET /search
│   │   │   └── feed.go               # GET /videos, /feed/trending, categories
│   │   ├── middleware/
│   │   │   ├── auth.go               # JWT validation middleware
│   │   │   ├── cors.go               # CORS headers
│   │   │   └── logger.go             # Request logging
│   │   ├── model/
│   │   │   ├── user.go               # User struct + DB queries
│   │   │   ├── video.go              # Video struct + DB queries
│   │   │   ├── danmaku.go            # Danmaku struct + DB queries
│   │   │   ├── comment.go            # Comment struct + DB queries
│   │   │   └── social.go             # Like, follow, favorite queries
│   │   ├── service/
│   │   │   ├── auth_service.go       # Register/login logic, JWT generation
│   │   │   ├── video_service.go      # Upload orchestration, transcoding trigger
│   │   │   └── danmaku_service.go    # Danmaku validation, polling optimization
│   │   └── transcoder/
│   │       └── ffmpeg.go             # FFmpeg wrapper: raw → HLS qualities
│   ├── migrations/
│   │   ├── 001_create_users.up.sql
│   │   ├── 001_create_users.down.sql
│   │   ├── 002_create_categories.up.sql
│   │   ├── 003_create_videos.up.sql
│   │   ├── 004_create_video_qualities.up.sql
│   │   ├── 005_create_follows.up.sql
│   │   ├── 006_create_danmaku.up.sql
│   │   ├── 007_create_comments.up.sql
│   │   ├── 008_create_likes.up.sql
│   │   ├── 009_create_favorites.up.sql
│   │   └── 010_create_playlists.up.sql
│   ├── uploads/                      # Local file storage (dev)
│   │   ├── raw/                      # Original uploaded files
│   │   ├── hls/                      # Transcoding output (per video UUID)
│   │   └── thumbs/                   # Generated thumbnails
│   ├── Dockerfile
│   ├── go.mod
│   └── go.sum
├── frontend/                         # React Native (modeled after JKVideo)
│   ├── app/                          # expo-router pages
│   │   ├── index.tsx                 # Home feed (trending + categories)
│   │   ├── video/[id].tsx            # Video player + danmaku overlay
│   │   ├── search.tsx                # Search page
│   │   ├── upload.tsx                # Video upload screen
│   │   ├── profile/[id].tsx          # User profile + their videos
│   │   └── settings.tsx              # App settings
│   ├── components/
│   │   ├── VideoCard.tsx             # Feed card (thumbnail, title, stats)
│   │   ├── VideoPlayer.tsx           # react-native-video wrapper with controls
│   │   ├── DanmakuCanvas.tsx         # Canvas overlay for bullet comments
│   │   ├── DanmakuInput.tsx          # Inline input for sending danmaku
│   │   ├── CommentSection.tsx        # Threaded comment list
│   │   └── MiniPlayer.tsx            # Bottom floating player (P2)
│   ├── hooks/
│   │   ├── useAuth.ts                # Login state, token management
│   │   ├── useVideo.ts               # Fetch video detail + qualities
│   │   ├── useDanmaku.ts             # Poll danmaku + send danmaku
│   │   └── useFeed.ts                # Infinite scroll feed
│   ├── services/
│   │   └── api.ts                    # Axios instance + all API calls
│   ├── store/
│   │   ├── authStore.ts              # Zustand: user session
│   │   ├── videoStore.ts             # Zustand: current playing video state
│   │   └── playerStore.ts            # Zustand: player preferences (quality, volume)
│   └── utils/
│       ├── format.ts                 # Duration, count formatters
│       └── constants.ts              # API base URL, quality options
├── nginx/                            # HLS static file serving
│   └── nginx.conf
├── docker-compose.yml                # app + postgres + redis + nginx
└── README.md
```

---

### 6.6 Phased Roadmap

#### P0 — MVP (Weeks 1–6)
**Goal:** A working video upload → transcode → playback → danmaku loop

| Area | Deliverables |
|------|-------------|
| **Backend** | User registration/login (JWT), video upload + FFmpeg HLS transcode, video list/detail API, danmaku POST/GET (REST polling every 2s), like/unlike, trending feed |
| **Frontend** | Auth screens (login/register), home feed (FlatList), video player + danmaku overlay (Canvas + polling), upload screen with progress bar |
| **Infrastructure** | Docker Compose (Go server + PostgreSQL), migration runner, nginx for HLS serving |

#### P1 — Core Features (Weeks 7–10)
**Goal:** Community features and basic search

| Area | Deliverables |
|------|-------------|
| **Backend** | Threaded comments API, full-text search (PostgreSQL `tsvector`), subscribe/unsubscribe, personalized feed (basic category filter), favorites |
| **Frontend** | Comment section (threaded replies), search page, user profiles, subscribe button, favorite button |
| **Infrastructure** | Add Redis container for session caching |

#### P2 — Real-time (Weeks 11–13)
**Goal:** Replace polling with WebSocket, add real-time features

| Area | Deliverables |
|------|-------------|
| **Backend** | WebSocket danmaku server (gorilla/websocket), Redis pub/sub for cross-instance broadcast, view count increment API, watch history tracking |
| **Frontend** | Switch danmaku from polling → WebSocket, live view counter, in-feed video auto-preview, mini-player overlay (like JKVideo) |
| **Infrastructure** | Rate limiting via Redis, WebSocket connection pooling |

#### P3 — Polish & Scale (Weeks 14–17)
**Goal:** Production readiness within prototype scope

| Area | Deliverables |
|------|-------------|
| **Backend** | MinIO for object storage, thumbnail generation (FFmpeg snapshot), adaptive bitrate selection, basic creator analytics (views, engagement) |
| **Frontend** | Playlist management UI, video download support, settings page (quality preference, dark mode), loading skeletons |
| **Infrastructure** | Prometheus + Grafana monitoring, automated DB backups, GitHub Actions CI/CD pipeline |

---

### 6.7 Key Design Decisions (Informed by JKVideo)

| Decision | Reasoning |
|----------|-----------|
| **Go monolith over microservices** | At 10 concurrent users, a single Go binary with PostgreSQL handles thousands of QPS; microservices add deployment complexity with no benefit |
| **PostgreSQL for danmaku (not MongoDB)** | Keeps infrastructure to 1 database; danmaku per video is bounded (~10k max), well within PostgreSQL's range with proper indexing |
| **REST polling before WebSocket** | 2-second polling interval is imperceptible for danmaku at prototype scale; WebSocket adds connection management complexity for little early payoff |
| **Local disk storage over S3/MinIO in P0** | `/uploads/` directory + nginx static serving is 0 configuration; MinIO can be swapped in transparently when needed |
| **Zustand over Redux** | JKVideo validates Zustand for bilibili-like apps — simpler API, less boilerplate, perfect for small teams |
| **expo-router file-based routing** | Same pattern as JKVideo (`app/video/[id].tsx`, `app/index.tsx`); familiar structure, auto deep-linking |
| **FFmpeg subprocess over microservice** | Go calls `ffmpeg` directly as a subprocess for transcoding; no need for a separate Python service at this scale |
| **JWT (access + refresh tokens)** | Standard, stateless auth; refresh tokens stored in PostgreSQL with expiry; no Redis dependency for sessions |
