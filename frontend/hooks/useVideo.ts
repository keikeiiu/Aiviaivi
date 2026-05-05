import { useState, useEffect, useCallback } from "react";
import { videos as videosApi } from "../services/api";

interface VideoDetail {
  id: string;
  title: string;
  description: string;
  cover_url: string;
  duration: number;
  status: string;
  tags: string[];
  view_count: number;
  like_count: number;
  comment_count: number;
  share_count: number;
  user: { id: string; username: string; avatar_url: string };
  qualities: { id: string; quality: string; manifest_url: string }[];
  category: { id: number; name: string } | null;
}

export function useVideo(videoId: string | null) {
  const [video, setVideo] = useState<VideoDetail | null>(null);
  const [loading, setLoading] = useState(false);
  const [error, setError] = useState<string | null>(null);

  const fetch = useCallback(async () => {
    if (!videoId) return;
    setLoading(true);
    setError(null);
    try {
      const { data } = await videosApi.detail(videoId);
      setVideo(data.data);
    } catch {
      setError("Failed to load video");
    } finally {
      setLoading(false);
    }
  }, [videoId]);

  useEffect(() => {
    fetch();
  }, [fetch]);

  return { video, loading, error, refetch: fetch };
}
