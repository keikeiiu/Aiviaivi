import { useState, useRef } from "react";
import {
  View, Text, Image, TouchableOpacity, StyleSheet,
  PanResponder, Animated, Dimensions,
} from "react-native";
import { useRouter } from "expo-router";
import { useVideoStore } from "../store/videoStore";
import { usePlayerStore } from "../store/playerStore";

interface Props {
  videoId: string;
  title: string;
  coverUrl: string;
}

const { width: SCREEN_WIDTH, height: SCREEN_HEIGHT } = Dimensions.get("window");
const MINI_HEIGHT = 64;

export default function MiniPlayer({ videoId, title, coverUrl }: Props) {
  const router = useRouter();
  const isPlaying = useVideoStore((s) => s.isPlaying);
  const setPlaying = useVideoStore((s) => s.setPlaying);
  const reset = useVideoStore((s) => s.reset);
  const [dismissed, setDismissed] = useState(false);

  const pan = useRef(new Animated.ValueXY()).current;
  const opacity = useRef(new Animated.Value(1)).current;

  const panResponder = useRef(
    PanResponder.create({
      onMoveShouldSetPanResponder: (_, gs) => Math.abs(gs.dy) > 5,
      onPanResponderMove: (_, gs) => {
        if (gs.dy > 0) pan.setValue({ x: 0, y: gs.dy });
      },
      onPanResponderRelease: (_, gs) => {
        if (gs.dy > 80) {
          // Dismiss
          Animated.parallel([
            Animated.timing(pan, { toValue: { x: 0, y: 200 }, duration: 200, useNativeDriver: true }),
            Animated.timing(opacity, { toValue: 0, duration: 200, useNativeDriver: true }),
          ]).start(() => {
            setDismissed(true);
            reset();
          });
        } else {
          Animated.spring(pan, { toValue: { x: 0, y: 0 }, useNativeDriver: true }).start();
        }
      },
    })
  ).current;

  if (dismissed || !videoId) return null;

  return (
    <Animated.View
      style={[
        styles.container,
        { transform: [{ translateY: pan.y }], opacity },
      ]}
      {...panResponder.panHandlers}
    >
      <TouchableOpacity
        style={styles.touchable}
        onPress={() => router.push(`/video/${videoId}`)}
        activeOpacity={0.9}
      >
        <Image
          source={{ uri: coverUrl || "https://placehold.co/80x45/333/666?text=Video" }}
          style={styles.thumb}
        />
        <View style={styles.info}>
          <Text style={styles.title} numberOfLines={1}>{title}</Text>
          <Text style={styles.hint}>Tap to open · Swipe down to close</Text>
        </View>
        <TouchableOpacity style={styles.playBtn} onPress={() => setPlaying(!isPlaying)}>
          <Text style={styles.playIcon}>{isPlaying ? "⏸" : "▶"}</Text>
        </TouchableOpacity>
      </TouchableOpacity>
    </Animated.View>
  );
}

const styles = StyleSheet.create({
  container: {
    position: "absolute",
    bottom: 0,
    left: 0,
    right: 0,
    height: MINI_HEIGHT,
    backgroundColor: "#1a1a1a",
    borderTopWidth: 1,
    borderTopColor: "#333",
    zIndex: 100,
  },
  touchable: {
    flex: 1,
    flexDirection: "row",
    alignItems: "center",
    paddingHorizontal: 12,
    gap: 10,
  },
  thumb: {
    width: 80,
    height: 45,
    borderRadius: 4,
    backgroundColor: "#2a2a2a",
  },
  info: {
    flex: 1,
  },
  title: {
    color: "#fff",
    fontSize: 13,
    fontWeight: "600",
  },
  hint: {
    color: "#666",
    fontSize: 11,
    marginTop: 2,
  },
  playBtn: {
    width: 36,
    height: 36,
    borderRadius: 18,
    backgroundColor: "#333",
    justifyContent: "center",
    alignItems: "center",
  },
  playIcon: {
    color: "#fff",
    fontSize: 14,
  },
});
