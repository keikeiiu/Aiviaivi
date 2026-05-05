import { useState, useEffect } from "react";
import { View, Text, FlatList, TouchableOpacity, StyleSheet, ActivityIndicator } from "react-native";
import { useLocalSearchParams } from "expo-router";
import { users as usersApi, social as socialApi } from "../../services/api";
import { useAuthStore } from "../../store/authStore";
import VideoCard from "../../components/VideoCard";
import { formatCount } from "../../utils/format";

export default function ProfileScreen() {
  const { id } = useLocalSearchParams<{ id: string }>();
  const currentUser = useAuthStore((s) => s.user);
  const [profile, setProfile] = useState<any>(null);
  const [videos, setVideos] = useState<any[]>([]);
  const [loading, setLoading] = useState(true);
  const [isFollowing, setIsFollowing] = useState(false);

  useEffect(() => {
    if (!id) return;
    (async () => {
      try {
        const [profileRes, videosRes] = await Promise.all([
          usersApi.profile(id),
          usersApi.videos(id),
        ]);
        setProfile(profileRes.data.data);
        setIsFollowing(profileRes.data.data.is_following);
        setVideos(videosRes.data.data || []);
      } catch {} finally {
        setLoading(false);
      }
    })();
  }, [id]);

  const handleSubscribe = async () => {
    if (!id) return;
    try {
      if (isFollowing) {
        await socialApi.unsubscribe(id);
      } else {
        await socialApi.subscribe(id);
      }
      setIsFollowing(!isFollowing);
    } catch {}
  };

  if (loading) {
    return (
      <View style={styles.center}>
        <ActivityIndicator color="#00a1d6" size="large" />
      </View>
    );
  }

  if (!profile) {
    return (
      <View style={styles.center}>
        <Text style={styles.errorText}>User not found</Text>
      </View>
    );
  }

  const isOwnProfile = currentUser?.id === id;

  return (
    <View style={styles.container}>
      <View style={styles.header}>
        <Text style={styles.username}>@{profile.username}</Text>
        {profile.bio ? <Text style={styles.bio}>{profile.bio}</Text> : null}
        <View style={styles.stats}>
          <View style={styles.statItem}>
            <Text style={styles.statValue}>{formatCount(profile.follower_count || 0)}</Text>
            <Text style={styles.statLabel}>Followers</Text>
          </View>
          <View style={styles.statItem}>
            <Text style={styles.statValue}>{formatCount(profile.following_count || 0)}</Text>
            <Text style={styles.statLabel}>Following</Text>
          </View>
        </View>
        {!isOwnProfile && (
          <TouchableOpacity
            style={[styles.subBtn, isFollowing && styles.subBtnActive]}
            onPress={handleSubscribe}
          >
            <Text style={[styles.subBtnText, isFollowing && styles.subBtnTextActive]}>
              {isFollowing ? "Following" : "Subscribe"}
            </Text>
          </TouchableOpacity>
        )}
      </View>

      <FlatList
        data={videos}
        keyExtractor={(item) => item.id}
        renderItem={({ item }) => (
          <VideoCard
            id={item.id}
            title={item.title}
            cover_url={item.cover_url}
            duration={item.duration}
            view_count={item.view_count}
            user={{ username: profile.username, avatar_url: profile.avatar_url }}
            created_at={item.created_at}
          />
        )}
        contentContainerStyle={styles.list}
        ListEmptyComponent={<Text style={styles.empty}>No videos yet</Text>}
      />
    </View>
  );
}

const styles = StyleSheet.create({
  container: { flex: 1, backgroundColor: "#0f0f0f" },
  center: { flex: 1, justifyContent: "center", alignItems: "center", backgroundColor: "#0f0f0f" },
  errorText: { color: "#999", fontSize: 16 },
  header: { padding: 16, alignItems: "center", borderBottomWidth: 1, borderBottomColor: "#1a1a1a" },
  username: { color: "#fff", fontSize: 22, fontWeight: "700", marginBottom: 4 },
  bio: { color: "#999", fontSize: 14, textAlign: "center", marginBottom: 12 },
  stats: { flexDirection: "row", gap: 32, marginBottom: 16 },
  statItem: { alignItems: "center" },
  statValue: { color: "#fff", fontSize: 18, fontWeight: "700" },
  statLabel: { color: "#666", fontSize: 12, marginTop: 2 },
  subBtn: { paddingHorizontal: 32, paddingVertical: 10, borderRadius: 20, backgroundColor: "#00a1d6" },
  subBtnActive: { backgroundColor: "#333" },
  subBtnText: { color: "#fff", fontSize: 14, fontWeight: "600" },
  subBtnTextActive: { color: "#999" },
  list: { paddingHorizontal: 12, paddingTop: 16 },
  empty: { color: "#666", textAlign: "center", marginTop: 60, fontSize: 16 },
});
