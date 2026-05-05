import { useState, useCallback } from "react";
import { videos as videosApi, trending as trendingApi } from "../services/api";

interface FeedItem {
  id: string;
  title: string;
  cover_url: string;
  duration: number;
  view_count: number;
  like_count: number;
  user: { id: string; username: string; avatar_url: string };
  created_at: string;
}

export function useFeed(type: "latest" | "trending" = "latest", categoryId?: number) {
  const [items, setItems] = useState<FeedItem[]>([]);
  const [page, setPage] = useState(1);
  const [hasMore, setHasMore] = useState(true);
  const [loading, setLoading] = useState(false);

  const fetch = useCallback(
    async (pageNum: number) => {
      if (loading) return;
      setLoading(true);
      try {
        const api = type === "trending" ? trendingApi : videosApi.list;
        const { data } = await api({
          page: pageNum,
          ...(categoryId ? { category: categoryId } : {}),
          ...(type === "trending" ? {} : { sort: type }),
        });
        const newItems = data.data || [];
        if (pageNum === 1) {
          setItems(newItems);
        } else {
          setItems((prev) => [...prev, ...newItems]);
        }
        setHasMore(newItems.length >= 20);
        setPage(pageNum);
      } catch {} finally {
        setLoading(false);
      }
    },
    [type, categoryId, loading]
  );

  const refresh = useCallback(() => fetch(1), [fetch]);
  const loadMore = useCallback(() => {
    if (hasMore && !loading) fetch(page + 1);
  }, [hasMore, loading, page, fetch]);

  return { items, loading, hasMore, refresh, loadMore };
}
