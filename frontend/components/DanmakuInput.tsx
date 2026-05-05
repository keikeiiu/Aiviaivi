import { useState } from "react";
import { View, TextInput, TouchableOpacity, Text, StyleSheet } from "react-native";
import { useAuthStore } from "../store/authStore";

interface Props {
  onSubmit: (content: string) => void;
}

export default function DanmakuInput({ onSubmit }: Props) {
  const [text, setText] = useState("");
  const isAuthenticated = useAuthStore((s) => s.isAuthenticated);

  const handleSend = () => {
    const trimmed = text.trim();
    if (trimmed) {
      onSubmit(trimmed);
      setText("");
    }
  };

  if (!isAuthenticated) return null;

  return (
    <View style={styles.container}>
      <TextInput
        style={styles.input}
        value={text}
        onChangeText={setText}
        placeholder="Send a danmaku..."
        placeholderTextColor="#666"
        maxLength={100}
        onSubmitEditing={handleSend}
        returnKeyType="send"
      />
      <TouchableOpacity style={styles.sendBtn} onPress={handleSend}>
        <Text style={styles.sendText}>Send</Text>
      </TouchableOpacity>
    </View>
  );
}

const styles = StyleSheet.create({
  container: { flexDirection: "row", alignItems: "center", paddingHorizontal: 12, paddingVertical: 8, backgroundColor: "rgba(0,0,0,0.6)" },
  input: { flex: 1, color: "#fff", backgroundColor: "rgba(255,255,255,0.15)", borderRadius: 20, paddingHorizontal: 16, paddingVertical: 8, fontSize: 14 },
  sendBtn: { marginLeft: 8, paddingHorizontal: 16, paddingVertical: 8, backgroundColor: "#00a1d6", borderRadius: 20 },
  sendText: { color: "#fff", fontSize: 14, fontWeight: "600" },
});
