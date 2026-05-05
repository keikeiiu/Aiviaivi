import { useState } from "react";
import { View, TextInput, FlatList, Text, StyleSheet } from "react-native";
import { search as searchApi } from "../services/api";
import VideoCard from "../components/VideoCard";

export default function SearchScreen() {
  const [query, setQuery] = useState("");
  const [results, setResults] = useState<any[]>([]);
  const [loading, setLoading] = useState(false);

  const handleSearch = async (q: string) => {
    setQuery(q);
    if (q.trim().length < 2) {
      setResults([]);
      return;
    }
    setLoading(true);
    try {
      const { data } = await searchApi(q.trim());
      setResults(data.data || []);
    } catch {} finally {
      setLoading(false);
    }
  };

  return (
    <View style={styles.container}>
      <TextInput
        style={styles.input}
        value={query}
        onChangeText={handleSearch}
        placeholder="Search videos..."
        placeholderTextColor="#666"
        autoFocus
        returnKeyType="search"
      />
      <FlatList
        data={results}
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
        contentContainerStyle={styles.list}
        ListEmptyComponent={
          query.length >= 2 && !loading ? (
            <Text style={styles.empty}>No results for "{query}"</Text>
          ) : null
        }
      />
    </View>
  );
}

const styles = StyleSheet.create({
  container: { flex: 1, backgroundColor: "#0f0f0f" },
  input: { margin: 12, backgroundColor: "#1a1a1a", borderRadius: 12, paddingHorizontal: 16, paddingVertical: 12, color: "#fff", fontSize: 16 },
  list: { paddingHorizontal: 12 },
  empty: { color: "#666", textAlign: "center", marginTop: 40, fontSize: 14 },
});
