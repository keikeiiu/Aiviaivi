import { Stack } from "expo-router";
import { StatusBar } from "expo-status-bar";
import { useAuth } from "../hooks/useAuth";
import { useVideoStore } from "../store/videoStore";
import MiniPlayer from "../components/MiniPlayer";
import { View, Text, ActivityIndicator } from "react-native";

export default function RootLayout() {
  const { isLoading } = useAuth();
  const currentVideoId = useVideoStore((s) => s.currentVideoId);
  const currentTitle = useVideoStore((s) => s.currentTitle);
  const currentCover = useVideoStore((s) => s.currentCover);

  if (isLoading) {
    return (
      <View style={{ flex: 1, justifyContent: "center", alignItems: "center", backgroundColor: "#0f0f0f" }}>
        <ActivityIndicator color="#00a1d6" size="large" />
      </View>
    );
  }

  return (
    <View style={{ flex: 1, backgroundColor: "#0f0f0f" }}>
      <StatusBar style="light" />
      <Stack
        screenOptions={{
          headerStyle: { backgroundColor: "#0f0f0f" },
          headerTintColor: "#fff",
          headerTitleStyle: { fontWeight: "600" },
          contentStyle: { backgroundColor: "#0f0f0f" },
          animation: "slide_from_right",
        }}
      >
        <Stack.Screen name="index" options={{ title: "AiliVili", headerTitleAlign: "center" }} />
        <Stack.Screen name="video/[id]" options={{ headerShown: false }} />
        <Stack.Screen name="search" options={{ title: "Search" }} />
        <Stack.Screen name="upload" options={{ title: "Upload Video" }} />
        <Stack.Screen name="playlists" options={{ title: "Playlists" }} />
        <Stack.Screen name="settings" options={{ title: "Settings" }} />
        <Stack.Screen name="profile/[id]" options={{ title: "Profile" }} />
        <Stack.Screen name="login" options={{ title: "Sign In", headerShown: false }} />
        <Stack.Screen name="register" options={{ title: "Sign Up", headerShown: false }} />
      </Stack>

      {currentVideoId ? (
        <MiniPlayer videoId={currentVideoId} title={currentTitle} coverUrl={currentCover} />
      ) : null}
    </View>
  );
}
