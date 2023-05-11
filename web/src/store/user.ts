import { create } from "zustand";
import { UserProfile } from "@api/v1/user/profile";

interface UserState {
  token: string | null;
  groups: string[];
  profile: UserProfile | null;

  setAuth: (token: string | null, groups?: string[]) => void;
  setState: <T extends keyof UserState>(
    key: T
  ) => (value: UserState[T]) => void;
}

export const useUser = create<UserState>()((set) => ({
  token: localStorage.getItem("token"),
  groups: localStorage.getItem("groups")?.split(",") || [],
  profile: null,

  setAuth: (token, groups) => {
    if (token) localStorage.setItem("token", token);
    else localStorage.removeItem("token");
    if (groups) localStorage.setItem("groups", groups.join(","));
    else {
      localStorage.removeItem("groups");
      groups = [];
    }
    set({ token, groups });
  },
  setState: (key) => (value) => set({ [key]: value }),
}));
