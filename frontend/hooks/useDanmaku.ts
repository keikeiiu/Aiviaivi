import { useState, useEffect, useRef, useCallback } from "react";
import { danmaku as danmakuApi } from "../services/api";
import { DANMAKU_POLL_INTERVAL_MS, WS_BASE_URL } from "../utils/constants";
import { useAuthStore } from "../store/authStore";

interface DanmakuItem {
  id: number;
  content: string;
  video_time: number;
  color: string;
  font_size: string;
  mode: string;
  username: string;
}

export function useDanmaku(videoId: string | null, currentTime: number) {
  const [items, setItems] = useState<DanmakuItem[]>([]);
  const [viewCount, setViewCount] = useState(0);
  const [useWebSocket, setUseWebSocket] = useState(false);
  const wsRef = useRef<WebSocket | null>(null);
  const token = useAuthStore((s) => s.token);

  // WebSocket connection
  useEffect(() => {
    if (!videoId || !useWebSocket) return;

    const wsUrl = `${WS_BASE_URL}/videos/${videoId}/danmaku/ws${token ? `?token=${token}` : ""}`;
    const ws = new WebSocket(wsUrl);
    wsRef.current = ws;

    ws.onmessage = (event) => {
      try {
        const msg = JSON.parse(event.data);
        if (msg.type === "danmaku") {
          setItems((prev) => [...prev.slice(-200), msg as DanmakuItem]);
        } else if (msg.type === "view_count") {
          setViewCount(msg.count || 0);
        }
      } catch {}
    };

    ws.onerror = () => {
      setUseWebSocket(false); // Fall back to polling
    };

    return () => {
      ws.close();
      wsRef.current = null;
    };
  }, [videoId, useWebSocket, token]);

  // REST polling fallback
  useEffect(() => {
    if (!videoId || useWebSocket) return;

    const poll = async () => {
      try {
        const { data } = await danmakuApi.list(videoId, currentTime - 2, currentTime + 5);
        const newItems = data.data || [];
        setItems((prev) => {
          const existingIds = new Set(prev.map((i) => i.id));
          const merged = [...prev];
          for (const item of newItems) {
            if (!existingIds.has(item.id)) merged.push(item);
          }
          return merged.slice(-300);
        });
      } catch {}
    };

    poll();
    const interval = setInterval(poll, DANMAKU_POLL_INTERVAL_MS);
    return () => clearInterval(interval);
  }, [videoId, currentTime, useWebSocket]);

  const sendDanmaku = useCallback(
    async (content: string, videoTime: number, color?: string, mode?: string) => {
      if (!videoId) return;
      try {
        await danmakuApi.send(videoId, {
          content,
          video_time: videoTime,
          color: color || "#FFFFFF",
          mode: mode || "scroll",
        });
      } catch {}
    },
    [videoId]
  );

  return { items, viewCount, sendDanmaku, useWebSocket, setUseWebSocket };
}
