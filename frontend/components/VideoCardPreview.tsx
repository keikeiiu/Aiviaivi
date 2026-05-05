import { useRef, useEffect } from "react";
import { View, Text, Image, TouchableOpacity, StyleSheet } from "react-native";
import { Video, ResizeMode } from "expo-av";
import { useRouter } from "expo-router";
import { formatCount, formatDuration } from "../utils/format";
import { HLS_BASE_URL } from "../utils/constants";

interface Props {
  id: string;
  title: string;
  cover_url: string;
  duration: number;
  view_count: number;
  user: { username: string; avatar_url: string };
  isFocused: boolean;
  manifestUrl?: string;
}

export default function VideoCardPreview({
  id, title, cover_url, duration, view_count, user, isFocused, manifestUrl,
}: Props) {
  const router = useRouter();
  const videoRef = useRef<Video>(null);

  useEffect(() => {
    if (videoRef.current) {
      if (isFocused && manifestUrl) {
        videoRef.current.playAsync();
      } else {
        videoRef.current.pauseAsync();
      }
    }
  }, [isFocused, manifestUrl]);

  const fullUrl = manifestUrl?.startsWith("http")
    ? manifestUrl
    : manifestUrl ? `${HLS_BASE_URL}${manifestUrl}` : null;

  return (
    <TouchableOpacity
      style={styles.card}
      onPress={() => router.push(`/video/${id}`)}
      activeOpacity={0.9}
    >
      <View style={styles.media}>
        {isFocused && fullUrl ? (
          <Video
            ref={videoRef}
            source={{ uri: fullUrl }}
            style={styles.video}
            resizeMode={ResizeMode.COVER}
            shouldPlay
            isMuted
            isLooping
            progressUpdateIntervalMillis={500}
          />
        ) : (
          <Image
            source={{ uri: cover_url || "https://placehold.co/400x225/333/666?text=Video" }}
            style={styles.image}
          />
        )}
        <View style={styles.durationBadge}>
          <Text style={styles.durationText}>{formatDuration(duration)}</Text>
        </View>
        {isFocused && !fullUrl && (
          <View style={styles.previewBadge}>
            <Text style={styles.previewText}>Preview</Text>
          </View>
        )}
      </View>
      <View style={styles.info}>
        <Text style={styles.title} numberOfLines={2}>{title}</Text>
        <View style={styles.meta}>
          <Text style={styles.metaText}>{user.username}</Text>
          <Text style={styles.metaDot}> · </Text>
          <Text style={styles.metaText}>{formatCount(view_count)} views</Text>
        </View>
      </View>
    </TouchableOpacity>
  );
}

const styles = StyleSheet.create({
  card: { marginBottom: 8 },
  media: { position: "relative", aspectRatio: 16 / 9, backgroundColor: "#1a1a1a" },
  video: { width: "100%", height: "100%" },
  image: { width: "100%", height: "100%" },
  durationBadge: { position: "absolute", bottom: 6, right: 6, backgroundColor: "rgba(0,0,0,0.8)", paddingHorizontal: 6, paddingVertical: 2, borderRadius: 4 },
  durationText: { color: "#fff", fontSize: 12, fontWeight: "600" },
  previewBadge: { position: "absolute", top: 6, left: 6, backgroundColor: "rgba(0,161,214,0.8)", paddingHorizontal: 8, paddingVertical: 3, borderRadius: 4 },
  previewText: { color: "#fff", fontSize: 11, fontWeight: "600" },
  info: { paddingTop: 8, paddingHorizontal: 4 },
  title: { color: "#fff", fontSize: 14, fontWeight: "600", lineHeight: 20 },
  meta: { flexDirection: "row", alignItems: "center", marginTop: 4 },
  metaText: { color: "#999", fontSize: 12 },
  metaDot: { color: "#555", fontSize: 12 },
});
