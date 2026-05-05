import { View, Text, Image, TouchableOpacity, StyleSheet } from "react-native";
import { useRouter } from "expo-router";
import { formatCount, formatDuration, formatTimeAgo } from "../utils/format";

interface Props {
  id: string;
  title: string;
  cover_url: string;
  duration: number;
  view_count: number;
  user: { username: string; avatar_url: string };
  created_at: string;
}

export default function VideoCard({ id, title, cover_url, duration, view_count, user, created_at }: Props) {
  const router = useRouter();

  return (
    <TouchableOpacity style={styles.card} onPress={() => router.push(`/video/${id}`)} activeOpacity={0.8}>
      <View style={styles.cover}>
        <Image source={{ uri: cover_url || "https://placehold.co/400x225/333/666?text=Video" }} style={styles.image} />
        <View style={styles.durationBadge}>
          <Text style={styles.durationText}>{formatDuration(duration)}</Text>
        </View>
      </View>
      <View style={styles.info}>
        <Text style={styles.title} numberOfLines={2}>{title}</Text>
        <View style={styles.meta}>
          <Text style={styles.metaText}>{user.username}</Text>
          <Text style={styles.metaDot}> · </Text>
          <Text style={styles.metaText}>{formatCount(view_count)} views</Text>
          <Text style={styles.metaDot}> · </Text>
          <Text style={styles.metaText}>{formatTimeAgo(created_at)}</Text>
        </View>
      </View>
    </TouchableOpacity>
  );
}

const styles = StyleSheet.create({
  card: { marginBottom: 16 },
  cover: { position: "relative", aspectRatio: 16 / 9, borderRadius: 8, overflow: "hidden", backgroundColor: "#1a1a1a" },
  image: { width: "100%", height: "100%" },
  durationBadge: { position: "absolute", bottom: 6, right: 6, backgroundColor: "rgba(0,0,0,0.8)", paddingHorizontal: 6, paddingVertical: 2, borderRadius: 4 },
  durationText: { color: "#fff", fontSize: 12, fontWeight: "600" },
  info: { paddingTop: 8, paddingHorizontal: 4 },
  title: { color: "#fff", fontSize: 14, fontWeight: "600", lineHeight: 20 },
  meta: { flexDirection: "row", alignItems: "center", marginTop: 4 },
  metaText: { color: "#999", fontSize: 12 },
  metaDot: { color: "#555", fontSize: 12 },
});
