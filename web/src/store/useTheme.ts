import { create } from "zustand";
import { persist, createJSONStorage } from "zustand/middleware";

interface ThemeState {
  dark?: boolean;

  setDarkMode: (dark: boolean) => void;
}

export const useTheme = create<ThemeState>()(
  persist(
    (set, get) => ({
      dark: window.matchMedia("(prefers-color-scheme: dark)").matches ?? true,

      setDarkMode: (dark: boolean) => {
        set({ dark });
      },
    }),
    {
      name: "theme",
      storage: createJSONStorage(() => localStorage),
    },
  ),
);
export default useTheme;
