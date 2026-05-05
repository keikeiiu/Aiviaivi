import { create } from "zustand";

interface VideoState {
  currentVideoId: string | null;
  isPlaying: boolean;
  progress: number;
  duration: number;

  setVideo: (id: string) => void;
  setPlaying: (playing: boolean) => void;
  setProgress: (progress: number) => void;
  setDuration: (duration: number) => void;
  reset: () => void;
}

export const useVideoStore = create<VideoState>((set) => ({
  currentVideoId: null,
  isPlaying: false,
  progress: 0,
  duration: 0,

  setVideo: (id) => set({ currentVideoId: id, progress: 0, duration: 0, isPlaying: true }),
  setPlaying: (playing) => set({ isPlaying: playing }),
  setProgress: (progress) => set({ progress }),
  setDuration: (duration) => set({ duration }),
  reset: () => set({ currentVideoId: null, isPlaying: false, progress: 0, duration: 0 }),
}));
