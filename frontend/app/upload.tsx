import { useState } from "react";
import { View, Text, TextInput, TouchableOpacity, StyleSheet, Alert } from "react-native";
import { useRouter } from "expo-router";
import * as DocumentPicker from "expo-document-picker";
import { videos as videosApi } from "../services/api";

export default function UploadScreen() {
  const router = useRouter();
  const [title, setTitle] = useState("");
  const [description, setDescription] = useState("");
  const [tags, setTags] = useState("");
  const [file, setFile] = useState<{ uri: string; name: string } | null>(null);
  const [uploading, setUploading] = useState(false);

  const pickFile = async () => {
    try {
      const result = await DocumentPicker.getDocumentAsync({
        type: "video/*",
        copyToCacheDirectory: true,
      });
      if (!result.canceled && result.assets?.[0]) {
        setFile({ uri: result.assets[0].uri, name: result.assets[0].name });
      }
    } catch {}
  };

  const handleUpload = async () => {
    if (!title.trim() || !file) {
      Alert.alert("Error", "Title and video file are required");
      return;
    }
    setUploading(true);
    try {
      const formData = new FormData();
      formData.append("title", title.trim());
      formData.append("description", description.trim());
      formData.append("tags", tags.trim());
      formData.append("file", {
        uri: file.uri,
        name: file.name,
        type: "video/mp4",
      } as any);

      await videosApi.upload(formData);
      Alert.alert("Success", "Video uploaded! Processing will begin shortly.", [
        { text: "OK", onPress: () => router.back() },
      ]);
    } catch {
      Alert.alert("Error", "Upload failed");
    } finally {
      setUploading(false);
    }
  };

  return (
    <View style={styles.container}>
      <TouchableOpacity style={styles.filePicker} onPress={pickFile}>
        <Text style={styles.filePickerText}>
          {file ? `📹 ${file.name}` : "Tap to select video file"}
        </Text>
      </TouchableOpacity>

      <TextInput
        style={styles.input}
        value={title}
        onChangeText={setTitle}
        placeholder="Video title *"
        placeholderTextColor="#666"
        maxLength={200}
      />
      <TextInput
        style={[styles.input, styles.textArea]}
        value={description}
        onChangeText={setDescription}
        placeholder="Description"
        placeholderTextColor="#666"
        multiline
        numberOfLines={4}
      />
      <TextInput
        style={styles.input}
        value={tags}
        onChangeText={setTags}
        placeholder="Tags (comma-separated)"
        placeholderTextColor="#666"
      />

      <TouchableOpacity
        style={[styles.uploadBtn, (!title.trim() || !file || uploading) && styles.uploadBtnDisabled]}
        onPress={handleUpload}
        disabled={!title.trim() || !file || uploading}
      >
        <Text style={styles.uploadBtnText}>
          {uploading ? "Uploading..." : "Upload Video"}
        </Text>
      </TouchableOpacity>
    </View>
  );
}

const styles = StyleSheet.create({
  container: { flex: 1, backgroundColor: "#0f0f0f", padding: 16, gap: 16 },
  filePicker: { backgroundColor: "#1a1a1a", borderRadius: 12, padding: 24, alignItems: "center", borderWidth: 1, borderColor: "#333", borderStyle: "dashed" },
  filePickerText: { color: "#00a1d6", fontSize: 14 },
  input: { backgroundColor: "#1a1a1a", borderRadius: 12, paddingHorizontal: 16, paddingVertical: 12, color: "#fff", fontSize: 16 },
  textArea: { height: 100, textAlignVertical: "top" },
  uploadBtn: { backgroundColor: "#00a1d6", borderRadius: 12, paddingVertical: 16, alignItems: "center" },
  uploadBtnDisabled: { opacity: 0.4 },
  uploadBtnText: { color: "#fff", fontSize: 16, fontWeight: "700" },
});
