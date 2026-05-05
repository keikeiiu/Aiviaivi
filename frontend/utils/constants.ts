export const API_BASE_URL = __DEV__
  ? "http://localhost:8080/api/v1"
  : "https://api.ailivili.com/api/v1";

export const WS_BASE_URL = __DEV__
  ? "ws://localhost:8080/api/v1"
  : "wss://api.ailivili.com/api/v1";

export const HLS_BASE_URL = __DEV__
  ? "http://localhost:8081"
  : "https://cdn.ailivili.com";

export const PAGE_SIZE = 20;
export const DANMAKU_POLL_INTERVAL_MS = 2000;
