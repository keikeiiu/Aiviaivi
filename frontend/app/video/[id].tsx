import { useState, useCallback } from "react";
import { View, Text, ScrollView, TouchableOpacity, StyleSheet, Dimensions, ActivityIndicator } from "react-native";
import { useLocalSearchParams, useRouter } from "expo-router";
import { useVideo } from "../../hooks/useVideo";
import { useDanmaku } from "../../hooks/useDanmaku";
import { useEffect } from "react";
import { useVideoStore } from "../../store/videoStore";
import { usePlayerStore } from "../../store/playerStore";
import { useAuthStore } from "../../store/authStore";
import VideoPlayer from "../../components/VideoPlayer";
import DanmakuCanvas from "../../components/DanmakuCanvas";
import DanmakuInput from "../../components/DanmakuInput";
import CommentSection from "../../components/CommentSection";
import { social as socialApi } from "../../services/api";
import { formatCount } from "../../utils/format";

const { width: SCREEN_WIDTH } = Dimensions.get("window");

export default function VideoScreen() {
  const { id } = useLocalSearchParams<{ id: string }>();
  const router = useRouter();
  const { video, loading, error } = useVideo(id || null);
  const progress = useVideoStore((s) => s.progress);
  const setVideo = useVideoStore((s) => s.setVideo);
  const { items: danmakuItems, viewCount, useWebSocket: wsActive, sendDanmaku } = useDanmaku(id || null, progress);

  // Update mini-player info when video loads
  useEffect(() => {
    if (video && id) {
      setVideo(id, video.title, video.cover_url);
    }
  }, [video, id, setVideo]);
  const [liked, setLiked] = useState(false);
  const [favorited, setFavorited] = useState(false);
  const isAuth = useAuthStore((s) => s.isAuthenticated);

  const handleDanmakuSend = useCallback(
    (content: string) => sendDanmaku(content, progress),
    [sendDanmaku, progress]
  );

  const handleLike = async () => {
    if (!id || !isAuth) return;
    try {
      if (liked) {
        await socialApi.unlike(id);
        setLiked(false);
      } else {
        await socialApi.like(id);
        setLiked(true);
      }
    } catch {}
  };

  const handleFavorite = async () => {
    if (!id || !isAuth) return;
    try {
      if (favorited) {
        await socialApi.unfavorite(id);
        setFavorited(false);
      } else {
        await socialApi.favorite(id);
        setFavorited(true);
      }
    } catch {}
  };

  if (loading) {
    return (
      <View style={styles.center}>
        <ActivityIndicator color="#00a1d6" size="large" />
      </View>
    );
  }

  if (error || !video) {
    return (
      <View style={styles.center}>
        <Text style={styles.errorText}>{error || "Video not found"}</Text>
        <TouchableOpacity onPress={() => router.back()} style={styles.backBtn}>
          <Text style={styles.backBtnText}>Go Back</Text>
        </TouchableOpacity>
      </View>
    );
  }

  const manifestUrl = video.qualities?.[0]?.manifest_url || "";

  return (
    <View style={styles.container}>
      <ScrollView bounces={false}>
        <View style={styles.playerContainer}>
          <VideoPlayer manifestUrl={manifestUrl} />
          <DanmakuCanvas items={danmakuItems} currentTime={progress} />
          <DanmakuInput onSubmit={handleDanmakuSend} />
        </View>

        <View style={styles.info}>
          <Text style={styles.title}>{video.title}</Text>
          <View style={styles.stats}>
            <Text style={styles.statText}>{formatCount(viewCount || video.view_count)} views{wsActive ? " · 🟢 Live" : ""}</Text>
            <Text style={styles.statText}> · {formatCount(video.like_count)} likes</Text>
            <Text style={styles.statText}> · {formatCount(video.comment_count)} comments</Text>
          </View>

          <TouchableOpacity
            style={styles.userRow}
            onPress={() => router.push(`/profile/${video.user.id}`)}
          >
            <Text style={styles.username}>@{video.user.username}</Text>
          </TouchableOpacity>

          {video.description ? (
            <Text style={styles.description}>{video.description}</Text>
          ) : null}

          <View style={styles.actions}>
            <TouchableOpacity style={styles.actionBtn} onPress={handleLike}>
              <Text style={[styles.actionIcon, liked && { color: "#ff4d4f" }]}>
                {liked ? "❤️" : "🤍"}
              </Text>
              <Text style={styles.actionLabel}>Like</Text>
            </TouchableOpacity>
            <TouchableOpacity style={styles.actionBtn} onPress={handleFavorite}>
              <Text style={[styles.actionIcon, favorited && { color: "#ffc107" }]}>
                {favorited ? "⭐" : "☆"}
              </Text>
              <Text style={styles.actionLabel}>Favorite</Text>
            </TouchableOpacity>
          </View>

          <View style={styles.qualities}>
            {video.qualities?.map((q) => (
              <View key={q.id} style={styles.qualityTag}>
                <Text style={styles.qualityText}>{q.quality}</Text>
              </View>
            ))}
          </View>
        </View>

        <CommentSection videoId={id!} />
      </ScrollView>

      <TouchableOpacity style={styles.closeBtn} onPress={() => router.back()}>
        <Text style={styles.closeBtnText}>✕</Text>
      </TouchableOpacity>
    </View>
  );
}

const styles = StyleSheet.create({
  container: { flex: 1, backgroundColor: "#0f0f0f" },
  center: { flex: 1, justifyContent: "center", alignItems: "center", backgroundColor: "#0f0f0f" },
  errorText: { color: "#999", fontSize: 16, marginBottom: 16 },
  backBtn: { paddingHorizontal: 20, paddingVertical: 10, backgroundColor: "#00a1d6", borderRadius: 8 },
  backBtnText: { color: "#fff", fontWeight: "600" },
  playerContainer: { width: SCREEN_WIDTH, aspectRatio: 16 / 9, backgroundColor: "#000", position: "relative" },
  info: { padding: 16 },
  title: { color: "#fff", fontSize: 18, fontWeight: "700", marginBottom: 8 },
  stats: { flexDirection: "row", marginBottom: 12 },
  statText: { color: "#999", fontSize: 13 },
  userRow: { marginBottom: 12 },
  username: { color: "#00a1d6", fontSize: 14, fontWeight: "600" },
  description: { color: "#ccc", fontSize: 14, lineHeight: 20, marginBottom: 16 },
  actions: { flexDirection: "row", gap: 24, marginBottom: 16 },
  actionBtn: { flexDirection: "row", alignItems: "center", gap: 6 },
  actionIcon: { fontSize: 20 },
  actionLabel: { color: "#999", fontSize: 14 },
  qualities: { flexDirection: "row", gap: 8, flexWrap: "wrap" },
  qualityTag: { backgroundColor: "#1a1a1a", paddingHorizontal: 10, paddingVertical: 4, borderRadius: 4 },
  qualityText: { color: "#00a1d6", fontSize: 12, fontWeight: "600" },
  closeBtn: { position: "absolute", top: 50, left: 16, width: 36, height: 36, borderRadius: 18, backgroundColor: "rgba(0,0,0,0.6)", justifyContent: "center", alignItems: "center" },
  closeBtnText: { color: "#fff", fontSize: 16 },
});
