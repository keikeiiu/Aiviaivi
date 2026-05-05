import { useRef, useEffect } from "react";
import { View, StyleSheet, TouchableOpacity, Text } from "react-native";
import { Video, ResizeMode, AVPlaybackStatus } from "expo-av";
import { useVideoStore } from "../store/videoStore";
import { usePlayerStore } from "../store/playerStore";
import { HLS_BASE_URL } from "../utils/constants";

interface Props {
  manifestUrl: string;
  onProgress?: (time: number) => void;
  onDuration?: (duration: number) => void;
}

export default function VideoPlayer({ manifestUrl, onProgress, onDuration }: Props) {
  const videoRef = useRef<Video>(null);
  const isPlaying = useVideoStore((s) => s.isPlaying);
  const setPlaying = useVideoStore((s) => s.setPlaying);
  const volume = usePlayerStore((s) => s.volume);
  const isMuted = usePlayerStore((s) => s.isMuted);

  const fullUrl = manifestUrl.startsWith("http")
    ? manifestUrl
    : `${HLS_BASE_URL}${manifestUrl}`;

  useEffect(() => {
    if (videoRef.current) {
      isPlaying ? videoRef.current.playAsync() : videoRef.current.pauseAsync();
    }
  }, [isPlaying]);

  const onPlaybackStatusUpdate = (status: AVPlaybackStatus) => {
    if (status.isLoaded) {
      if (status.didJustFinish) {
        setPlaying(false);
      }
      if (status.positionMillis !== undefined) {
        onProgress?.(status.positionMillis / 1000);
      }
      if (status.durationMillis !== undefined) {
        onDuration?.(status.durationMillis / 1000);
      }
    }
  };

  const togglePlay = () => setPlaying(!isPlaying);

  return (
    <TouchableOpacity style={styles.container} onPress={togglePlay} activeOpacity={1}>
      <Video
        ref={videoRef}
        source={{ uri: fullUrl }}
        style={styles.video}
        resizeMode={ResizeMode.CONTAIN}
        shouldPlay
        isMuted={isMuted}
        volume={volume}
        onPlaybackStatusUpdate={onPlaybackStatusUpdate}
        progressUpdateIntervalMillis={250}
      />
      {!isPlaying && (
        <View style={styles.playOverlay}>
          <Text style={styles.playIcon}>▶</Text>
        </View>
      )}
    </TouchableOpacity>
  );
}

const styles = StyleSheet.create({
  container: { width: "100%", aspectRatio: 16 / 9, backgroundColor: "#000" },
  video: { flex: 1 },
  playOverlay: { ...StyleSheet.absoluteFillObject, justifyContent: "center", alignItems: "center", backgroundColor: "rgba(0,0,0,0.3)" },
  playIcon: { color: "#fff", fontSize: 48, opacity: 0.9 },
});
