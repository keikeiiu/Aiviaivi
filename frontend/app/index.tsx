import { useEffect } from "react";
import { View, FlatList, Text, TouchableOpacity, StyleSheet } from "react-native";
import { useRouter } from "expo-router";
import { useFeed } from "../hooks/useFeed";
import { useAuth } from "../hooks/useAuth";
import VideoCard from "../components/VideoCard";

export default function HomeScreen() {
  const router = useRouter();
  const { isAuthenticated } = useAuth();
  const { items, loading, hasMore, refresh, loadMore } = useFeed("latest");

  useEffect(() => { refresh(); }, []);

  return (
    <View style={styles.container}>
      <View style={styles.header}>
        <TouchableOpacity onPress={() => router.push("/search")} style={styles.searchBar}>
          <Text style={styles.searchText}>Search videos...</Text>
        </TouchableOpacity>
        {isAuthenticated && (
          <TouchableOpacity onPress={() => router.push("/upload")} style={styles.uploadBtn}>
            <Text style={styles.uploadBtnText}>+</Text>
          </TouchableOpacity>
        )}
      </View>

      <FlatList
        data={items}
        keyExtractor={(item) => item.id}
        renderItem={({ item }) => (
          <VideoCard
            id={item.id}
            title={item.title}
            cover_url={item.cover_url}
            duration={item.duration}
            view_count={item.view_count}
            user={item.user}
            created_at={item.created_at}
          />
        )}
        refreshing={loading}
        onRefresh={refresh}
        onEndReached={loadMore}
        onEndReachedThreshold={0.5}
        contentContainerStyle={styles.list}
        ListEmptyComponent={
          !loading ? <Text style={styles.empty}>No videos yet</Text> : null
        }
      />
    </View>
  );
}

const styles = StyleSheet.create({
  container: { flex: 1, backgroundColor: "#0f0f0f" },
  header: { flexDirection: "row", alignItems: "center", paddingHorizontal: 12, paddingVertical: 8, gap: 8 },
  searchBar: { flex: 1, backgroundColor: "#1a1a1a", borderRadius: 20, paddingHorizontal: 16, paddingVertical: 10 },
  searchText: { color: "#666", fontSize: 14 },
  uploadBtn: { width: 40, height: 40, borderRadius: 20, backgroundColor: "#00a1d6", justifyContent: "center", alignItems: "center" },
  uploadBtnText: { color: "#fff", fontSize: 24, fontWeight: "300" },
  list: { paddingHorizontal: 12, paddingTop: 8 },
  empty: { color: "#666", textAlign: "center", marginTop: 100, fontSize: 16 },
});
