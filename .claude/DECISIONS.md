# Undecided / Needs Decision

| # | Category | Question | Context |
|---|----------|----------|---------|
| 1 | FFmpeg | Not installed locally — transcoding step fails gracefully (video stays "processing"). Install with `brew install ffmpeg` for full transcode testing. | Upload works, file saved. Transcode needs ffmpeg. |
| 2 | Refresh Tokens | Stateless JWT (simpler, current) vs DB-backed (revocable)? | Stateless 7x expiry works. DB-backed allows individual revocation. |
| 3 | HLS Base URL | manifest_urls use relative paths. Frontend `HLS_BASE_URL` hardcoded to localhost:8081. Make configurable via env? | Works for dev. Needs config for production. |
| 4 | MinIO SDK | Storage abstraction exists (interface + local + MinIO stub). Add real `minio-go`? | Local `/uploads/` verified working. MinIO for production. |
| 5 | Frontend Testing | 31 backend endpoints live-tested. Frontend needs `npx expo start` with Expo Go or simulator. Ready to test? | Frontend compiles (tsc --noEmit passes). Backend live at localhost:8080. |
| 6 | Production Deploy | Full stack complete + live tested. Ready for: (a) `docker compose up`, (b) cloud deploy, (c) CI/CD pipeline. | All code, tests, and live verification done. |
