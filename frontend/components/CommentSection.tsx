import { useState, useEffect } from "react";
import { View, Text, TextInput, TouchableOpacity, FlatList, StyleSheet } from "react-native";
import { comments as commentsApi } from "../services/api";
import { useAuthStore } from "../store/authStore";
import { formatTimeAgo } from "../utils/format";

interface Comment {
  id: string;
  content: string;
  username: string;
  like_count: number;
  created_at: string;
  replies: Comment[];
}

interface Props {
  videoId: string;
}

function CommentItem({ comment, depth = 0 }: { comment: Comment; depth?: number }) {
  return (
    <View style={[styles.comment, { marginLeft: depth * 20 }]}>
      <View style={styles.commentHeader}>
        <Text style={styles.username}>{comment.username}</Text>
        <Text style={styles.time}>{formatTimeAgo(comment.created_at)}</Text>
      </View>
      <Text style={styles.content}>{comment.content}</Text>
      {comment.replies?.map((reply) => (
        <CommentItem key={reply.id} comment={reply} depth={depth + 1} />
      ))}
    </View>
  );
}

export default function CommentSection({ videoId }: Props) {
  const [comments, setComments] = useState<Comment[]>([]);
  const [newComment, setNewComment] = useState("");
  const isAuthenticated = useAuthStore((s) => s.isAuthenticated);

  useEffect(() => {
    (async () => {
      try {
        const { data } = await commentsApi.list(videoId);
        setComments(data.data || []);
      } catch {}
    })();
  }, [videoId]);

  const handlePost = async () => {
    const trimmed = newComment.trim();
    if (!trimmed) return;
    try {
      await commentsApi.create(videoId, trimmed);
      setNewComment("");
      // Refetch
      const { data } = await commentsApi.list(videoId);
      setComments(data.data || []);
    } catch {}
  };

  return (
    <View style={styles.container}>
      <Text style={styles.heading}>Comments</Text>
      {isAuthenticated && (
        <View style={styles.inputRow}>
          <TextInput
            style={styles.input}
            value={newComment}
            onChangeText={setNewComment}
            placeholder="Add a comment..."
            placeholderTextColor="#666"
            onSubmitEditing={handlePost}
          />
          <TouchableOpacity style={styles.postBtn} onPress={handlePost}>
            <Text style={styles.postText}>Post</Text>
          </TouchableOpacity>
        </View>
      )}
      <FlatList
        data={comments}
        keyExtractor={(item) => item.id}
        renderItem={({ item }) => <CommentItem comment={item} />}
        scrollEnabled={false}
      />
    </View>
  );
}

const styles = StyleSheet.create({
  container: { padding: 16 },
  heading: { color: "#fff", fontSize: 16, fontWeight: "700", marginBottom: 12 },
  inputRow: { flexDirection: "row", marginBottom: 16 },
  input: { flex: 1, color: "#fff", backgroundColor: "#1a1a1a", borderRadius: 8, paddingHorizontal: 12, paddingVertical: 10, fontSize: 14 },
  postBtn: { marginLeft: 8, paddingHorizontal: 16, justifyContent: "center", backgroundColor: "#00a1d6", borderRadius: 8 },
  postText: { color: "#fff", fontWeight: "600", fontSize: 14 },
  comment: { marginBottom: 12 },
  commentHeader: { flexDirection: "row", alignItems: "center", gap: 8 },
  username: { color: "#00a1d6", fontSize: 13, fontWeight: "600" },
  time: { color: "#666", fontSize: 11 },
  content: { color: "#ddd", fontSize: 14, marginTop: 4, lineHeight: 20 },
});
