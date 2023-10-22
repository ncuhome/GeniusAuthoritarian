import { create } from "zustand";
import { persist, createJSONStorage } from "zustand/middleware";

interface ThemeState {
  dark?: boolean;

  setState: <T extends keyof ThemeState>(
    key: T
  ) => (value: ThemeState[T]) => void;
}

export const useTheme = create<ThemeState>()(
  persist(
    (set, get) => ({
      dark: window.matchMedia("(prefers-color-scheme: dark)").matches ?? true,

      setState: (key) => (value) => set({ [key]: value }),
    }),
    {
      name: "theme",
      storage: createJSONStorage(() => localStorage),
    }
  )
);
export default useTheme;
