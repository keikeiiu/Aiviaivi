import { View, Text, StyleSheet } from "react-native";
import { usePlayerStore } from "../store/playerStore";

interface DanmakuItem {
  id: number;
  content: string;
  video_time: number;
  color: string;
  font_size: string;
  mode: string;
}

interface Props {
  items: DanmakuItem[];
  currentTime: number;
}

const FONT_SIZES: Record<string, number> = { small: 14, medium: 18, large: 24 };
const SCROLL_SECONDS = 6;

export default function DanmakuCanvas({ items, currentTime }: Props) {
  const enabled = usePlayerStore((s) => s.danmakuEnabled);
  const opacity = usePlayerStore((s) => s.danmakuOpacity);

  if (!enabled || items.length === 0) return null;

  const visible = items.filter(
    (d) => Math.abs(d.video_time - currentTime) < SCROLL_SECONDS
  );

  if (visible.length === 0) return null;

  return (
    <View style={[StyleSheet.absoluteFill, { opacity }]} pointerEvents="none">
      {visible.map((item, index) => {
        const fontSize = FONT_SIZES[item.font_size] || 18;
        const elapsed = currentTime - item.video_time;
        const leftPct = ((SCROLL_SECONDS - elapsed) / SCROLL_SECONDS) * 100;

        return (
          <View
            key={item.id}
            style={{
              position: "absolute",
              top: `${(index * 7 + 5) % 85}%`,
              left: `${Math.max(0, Math.min(100, leftPct))}%`,
              transform: item.mode === "top" ? [{ translateY: -30 }] : item.mode === "bottom" ? [{ translateY: 30 }] : [],
            }}
          >
            <Text
              style={{
                color: item.color,
                fontSize,
                fontWeight: "600",
                textShadowColor: "rgba(0,0,0,0.8)",
                textShadowOffset: { width: 1, height: 1 },
                textShadowRadius: 2,
              }}
            >
              {item.content}
            </Text>
          </View>
        );
      })}
    </View>
  );
}
