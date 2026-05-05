import { useState, useCallback } from "react";
import { Alert, Platform } from "react-native";
import { Paths, File } from "expo-file-system";
import * as Sharing from "expo-sharing";

interface DownloadState {
  downloading: boolean;
  progress: number;
  error: string | null;
}

export function useDownload() {
  const [state, setState] = useState<DownloadState>({
    downloading: false,
    progress: 0,
    error: null,
  });

  const download = useCallback(async (url: string, filename: string) => {
    setState({ downloading: true, progress: 0, error: null });

    try {
      const destination = new File(Paths.document, filename);

      // Download file to device
      const result = await File.downloadFileAsync(url, destination);
      setState({ downloading: false, progress: 1, error: null });

      // Offer to share
      if (Platform.OS !== "web" && (await Sharing.isAvailableAsync())) {
        Alert.alert("Download Complete", filename, [
          { text: "OK", style: "cancel" },
          { text: "Share", onPress: () => Sharing.shareAsync(result.uri) },
        ]);
      } else {
        Alert.alert("Download Complete", filename);
      }
    } catch (e: any) {
      setState({ downloading: false, progress: 0, error: e?.message || "Download failed" });
      Alert.alert("Download Failed", e?.message || "Unknown error");
    }
  }, []);

  return { ...state, download };
}
