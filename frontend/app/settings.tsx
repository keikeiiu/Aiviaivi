import { View, Text, TouchableOpacity, StyleSheet, ScrollView } from "react-native";
import { usePlayerStore } from "../store/playerStore";
import { useAuthStore } from "../store/authStore";

const QUALITIES = ["360p", "480p", "720p", "1080p"] as const;

export default function SettingsScreen() {
  const quality = usePlayerStore((s) => s.quality);
  const setQuality = usePlayerStore((s) => s.setQuality);
  const danmakuEnabled = usePlayerStore((s) => s.danmakuEnabled);
  const toggleDanmaku = usePlayerStore((s) => s.toggleDanmaku);
  const danmakuOpacity = usePlayerStore((s) => s.danmakuOpacity);
  const setDanmakuOpacity = usePlayerStore((s) => s.setDanmakuOpacity);
  const logout = useAuthStore((s) => s.logout);
  const isAuth = useAuthStore((s) => s.isAuthenticated);

  return (
    <ScrollView style={styles.container}>
      <Text style={styles.section}>Playback</Text>
      <View style={styles.card}>
        <Text style={styles.label}>Default Quality</Text>
        <View style={styles.row}>
          {QUALITIES.map((q) => (
            <TouchableOpacity
              key={q}
              style={[styles.chip, quality === q && styles.chipActive]}
              onPress={() => setQuality(q)}
            >
              <Text style={[styles.chipText, quality === q && styles.chipTextActive]}>{q}</Text>
            </TouchableOpacity>
          ))}
        </View>
      </View>

      <Text style={styles.section}>Danmaku</Text>
      <View style={styles.card}>
        <View style={styles.switchRow}>
          <Text style={styles.label}>Enable Danmaku</Text>
          <TouchableOpacity
            style={[styles.toggle, danmakuEnabled && styles.toggleActive]}
            onPress={toggleDanmaku}
          >
            <Text style={[styles.toggleText, danmakuEnabled && styles.toggleTextActive]}>
              {danmakuEnabled ? "ON" : "OFF"}
            </Text>
          </TouchableOpacity>
        </View>
        <Text style={styles.label}>Opacity: {Math.round(danmakuOpacity * 100)}%</Text>
        <View style={styles.row}>
          {[0.3, 0.5, 0.7, 0.9].map((o) => (
            <TouchableOpacity
              key={o}
              style={[styles.chip, danmakuOpacity === o && styles.chipActive]}
              onPress={() => setDanmakuOpacity(o)}
            >
              <Text style={[styles.chipText, danmakuOpacity === o && styles.chipTextActive]}>
                {Math.round(o * 100)}%
              </Text>
            </TouchableOpacity>
          ))}
        </View>
      </View>

      {isAuth && (
        <>
          <Text style={styles.section}>Account</Text>
          <TouchableOpacity style={styles.logoutBtn} onPress={logout}>
            <Text style={styles.logoutText}>Log Out</Text>
          </TouchableOpacity>
        </>
      )}
    </ScrollView>
  );
}

const styles = StyleSheet.create({
  container: { flex: 1, backgroundColor: "#0f0f0f" },
  section: { color: "#999", fontSize: 12, fontWeight: "600", textTransform: "uppercase", marginTop: 24, marginBottom: 8, paddingHorizontal: 16 },
  card: { backgroundColor: "#1a1a1a", marginHorizontal: 16, borderRadius: 12, padding: 16 },
  label: { color: "#ccc", fontSize: 14, marginBottom: 12 },
  row: { flexDirection: "row", gap: 8, flexWrap: "wrap" },
  chip: { paddingHorizontal: 16, paddingVertical: 8, borderRadius: 8, backgroundColor: "#2a2a2a" },
  chipActive: { backgroundColor: "#00a1d6" },
  chipText: { color: "#999", fontSize: 13, fontWeight: "600" },
  chipTextActive: { color: "#fff" },
  switchRow: { flexDirection: "row", justifyContent: "space-between", alignItems: "center", marginBottom: 16 },
  toggle: { paddingHorizontal: 16, paddingVertical: 8, borderRadius: 8, backgroundColor: "#2a2a2a" },
  toggleActive: { backgroundColor: "#00a1d6" },
  toggleText: { color: "#999", fontSize: 13, fontWeight: "600" },
  toggleTextActive: { color: "#fff" },
  logoutBtn: { marginHorizontal: 16, marginTop: 16, padding: 14, backgroundColor: "#ff4d4f", borderRadius: 12, alignItems: "center" },
  logoutText: { color: "#fff", fontSize: 16, fontWeight: "700" },
});
