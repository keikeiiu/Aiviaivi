import { useState, useEffect, useCallback } from "react";
import { View, Text, FlatList, TouchableOpacity, TextInput, StyleSheet, Alert } from "react-native";
import { useRouter } from "expo-router";
import { playlists as playlistsApi } from "../services/api";
import { useAuthStore } from "../store/authStore";

interface Playlist {
  id: string;
  name: string;
  description: string;
  is_public: boolean;
  video_count: number;
}

export default function PlaylistsScreen() {
  const router = useRouter();
  const isAuth = useAuthStore((s) => s.isAuthenticated);
  const [playlists, setPlaylists] = useState<Playlist[]>([]);
  const [loading, setLoading] = useState(true);
  const [showCreate, setShowCreate] = useState(false);
  const [newName, setNewName] = useState("");
  const [newDesc, setNewDesc] = useState("");
  const [isPublic, setIsPublic] = useState(false);

  const fetchPlaylists = useCallback(async () => {
    try {
      const { data } = await playlistsApi.list();
      setPlaylists(data.data || []);
    } catch {} finally {
      setLoading(false);
    }
  }, []);

  useEffect(() => {
    if (isAuth) fetchPlaylists();
  }, [isAuth, fetchPlaylists]);

  const handleCreate = async () => {
    if (!newName.trim()) return;
    try {
      await playlistsApi.create({ name: newName.trim(), description: newDesc.trim(), is_public: isPublic });
      setShowCreate(false);
      setNewName("");
      setNewDesc("");
      setIsPublic(false);
      fetchPlaylists();
    } catch {
      Alert.alert("Error", "Failed to create playlist");
    }
  };

  const handleDelete = async (id: string, name: string) => {
    Alert.alert("Delete", `Delete "${name}"?`, [
      { text: "Cancel", style: "cancel" },
      { text: "Delete", style: "destructive", onPress: async () => {
        try {
          await playlistsApi.delete(id);
          fetchPlaylists();
        } catch {}
      }},
    ]);
  };

  if (!isAuth) {
    return (
      <View style={styles.center}>
        <Text style={styles.empty}>Sign in to manage playlists</Text>
      </View>
    );
  }

  return (
    <View style={styles.container}>
      <TouchableOpacity style={styles.createBtn} onPress={() => setShowCreate(!showCreate)}>
        <Text style={styles.createBtnText}>{showCreate ? "Cancel" : "+ New Playlist"}</Text>
      </TouchableOpacity>

      {showCreate && (
        <View style={styles.createForm}>
          <TextInput style={styles.input} value={newName} onChangeText={setNewName} placeholder="Playlist name" placeholderTextColor="#666" maxLength={100} />
          <TextInput style={[styles.input, styles.textArea]} value={newDesc} onChangeText={setNewDesc} placeholder="Description (optional)" placeholderTextColor="#666" multiline />
          <TouchableOpacity style={styles.publicToggle} onPress={() => setIsPublic(!isPublic)}>
            <Text style={styles.publicText}>{isPublic ? "Public" : "Private"}</Text>
          </TouchableOpacity>
          <TouchableOpacity style={styles.submitBtn} onPress={handleCreate}>
            <Text style={styles.submitText}>Create</Text>
          </TouchableOpacity>
        </View>
      )}

      <FlatList
        data={playlists}
        keyExtractor={(item) => item.id}
        refreshing={loading}
        onRefresh={fetchPlaylists}
        renderItem={({ item }) => (
          <TouchableOpacity
            style={styles.card}
            onPress={() => Alert.alert(item.name, `${item.video_count} videos · ${item.is_public ? "Public" : "Private"}${item.description ? '\n\n' + item.description : ''}`)}
          >
            <View style={styles.cardInfo}>
              <Text style={styles.cardTitle}>{item.name}</Text>
              <Text style={styles.cardMeta}>
                {item.video_count} videos · {item.is_public ? "Public" : "Private"}
              </Text>
            </View>
            <TouchableOpacity onPress={() => handleDelete(item.id, item.name)}>
              <Text style={styles.deleteBtn}>Delete</Text>
            </TouchableOpacity>
          </TouchableOpacity>
        )}
        ListEmptyComponent={!loading ? <Text style={styles.empty}>No playlists yet</Text> : null}
      />
    </View>
  );
}

const styles = StyleSheet.create({
  container: { flex: 1, backgroundColor: "#0f0f0f" },
  center: { flex: 1, justifyContent: "center", alignItems: "center", backgroundColor: "#0f0f0f" },
  empty: { color: "#666", fontSize: 16, textAlign: "center", marginTop: 60 },
  createBtn: { margin: 16, padding: 14, backgroundColor: "#00a1d6", borderRadius: 12, alignItems: "center" },
  createBtnText: { color: "#fff", fontSize: 16, fontWeight: "700" },
  createForm: { marginHorizontal: 16, marginBottom: 16, gap: 12 },
  input: { backgroundColor: "#1a1a1a", borderRadius: 8, paddingHorizontal: 14, paddingVertical: 12, color: "#fff", fontSize: 14 },
  textArea: { height: 60, textAlignVertical: "top" },
  publicToggle: { padding: 12, backgroundColor: "#1a1a1a", borderRadius: 8, alignItems: "center" },
  publicText: { color: "#00a1d6", fontSize: 14, fontWeight: "600" },
  submitBtn: { padding: 14, backgroundColor: "#00a1d6", borderRadius: 8, alignItems: "center" },
  submitText: { color: "#fff", fontSize: 16, fontWeight: "700" },
  card: { flexDirection: "row", alignItems: "center", marginHorizontal: 16, marginBottom: 12, padding: 16, backgroundColor: "#1a1a1a", borderRadius: 10 },
  cardInfo: { flex: 1 },
  cardTitle: { color: "#fff", fontSize: 16, fontWeight: "600" },
  cardMeta: { color: "#999", fontSize: 12, marginTop: 4 },
  deleteBtn: { color: "#ff4d4f", fontSize: 13, padding: 8 },
});
