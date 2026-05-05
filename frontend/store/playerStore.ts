import { create } from "zustand";

type Quality = "360p" | "480p" | "720p" | "1080p";

interface PlayerState {
  quality: Quality;
  volume: number;
  isMuted: boolean;
  isFullscreen: boolean;
  danmakuEnabled: boolean;
  danmakuOpacity: number;

  setQuality: (q: Quality) => void;
  setVolume: (v: number) => void;
  toggleMute: () => void;
  toggleFullscreen: () => void;
  toggleDanmaku: () => void;
  setDanmakuOpacity: (o: number) => void;
}

export const usePlayerStore = create<PlayerState>((set) => ({
  quality: "720p",
  volume: 1,
  isMuted: false,
  isFullscreen: false,
  danmakuEnabled: true,
  danmakuOpacity: 0.8,

  setQuality: (quality) => set({ quality }),
  setVolume: (volume) => set({ volume, isMuted: volume === 0 }),
  toggleMute: () => set((s) => ({ isMuted: !s.isMuted })),
  toggleFullscreen: () => set((s) => ({ isFullscreen: !s.isFullscreen })),
  toggleDanmaku: () => set((s) => ({ danmakuEnabled: !s.danmakuEnabled })),
  setDanmakuOpacity: (opacity) => set({ danmakuOpacity: opacity }),
}));
