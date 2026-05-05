import { create } from "zustand";

interface VideoState {
  currentVideoId: string | null;
  currentTitle: string;
  currentCover: string;
  isPlaying: boolean;
  progress: number;
  duration: number;

  setVideo: (id: string, title?: string, cover?: string) => void;
  setPlaying: (playing: boolean) => void;
  setProgress: (progress: number) => void;
  setDuration: (duration: number) => void;
  reset: () => void;
}

export const useVideoStore = create<VideoState>((set) => ({
  currentVideoId: null,
  currentTitle: "",
  currentCover: "",
  isPlaying: false,
  progress: 0,
  duration: 0,

  setVideo: (id, title = "", cover = "") =>
    set({ currentVideoId: id, currentTitle: title, currentCover: cover, progress: 0, duration: 0, isPlaying: true }),
  setPlaying: (playing) => set({ isPlaying: playing }),
  setProgress: (progress) => set({ progress }),
  setDuration: (duration) => set({ duration }),
  reset: () => set({ currentVideoId: null, currentTitle: "", currentCover: "", isPlaying: false, progress: 0, duration: 0 }),
}));
