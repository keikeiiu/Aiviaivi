import { Stack } from "expo-router";
import { StatusBar } from "expo-status-bar";
import { useAuth } from "../hooks/useAuth";
import { View, Text, ActivityIndicator } from "react-native";

export default function RootLayout() {
  const { isLoading } = useAuth();

  if (isLoading) {
    return (
      <View style={{ flex: 1, justifyContent: "center", alignItems: "center", backgroundColor: "#0f0f0f" }}>
        <ActivityIndicator color="#00a1d6" size="large" />
      </View>
    );
  }

  return (
    <>
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
        <Stack.Screen name="profile/[id]" options={{ title: "Profile" }} />
      </Stack>
    </>
  );
}
