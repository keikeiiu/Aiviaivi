import { useEffect, useRef } from "react";
import { View, Animated, StyleSheet } from "react-native";

interface Props {
  width?: number | string;
  height?: number;
  borderRadius?: number;
}

export function Skeleton({ width = "100%", height = 20, borderRadius = 4 }: Props) {
  const anim = useRef(new Animated.Value(0.3)).current;

  useEffect(() => {
    const loop = Animated.loop(
      Animated.sequence([
        Animated.timing(anim, { toValue: 1, duration: 800, useNativeDriver: true }),
        Animated.timing(anim, { toValue: 0.3, duration: 800, useNativeDriver: true }),
      ])
    );
    loop.start();
    return () => loop.stop();
  }, [anim]);

  return (
    <Animated.View
      style={[
        styles.base,
        {
          width: width as any,
          height,
          borderRadius,
          opacity: anim,
        },
      ]}
    />
  );
}

export function VideoCardSkeleton() {
  return (
    <View style={styles.card}>
      <Skeleton height={200} borderRadius={8} />
      <View style={styles.cardInfo}>
        <Skeleton width="80%" height={16} />
        <View style={styles.cardMeta}>
          <Skeleton width="30%" height={12} />
          <Skeleton width="20%" height={12} />
        </View>
      </View>
    </View>
  );
}

const styles = StyleSheet.create({
  base: { backgroundColor: "#2a2a2a" },
  card: { marginBottom: 16, paddingHorizontal: 12 },
  cardInfo: { paddingTop: 8, gap: 6 },
  cardMeta: { flexDirection: "row", gap: 8, marginTop: 4 },
});
