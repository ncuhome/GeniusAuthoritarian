import { create } from "zustand";
import { UserProfile } from "@api/v1/user/profile";

interface UserState {
  token: string | null;
  profile: UserProfile | null;

  setToken: (token: string | null) => void;
  setState: <T extends keyof UserState>(
    key: T
  ) => (value: UserState[T]) => void;
}

export const useUser = create<UserState>()((set) => ({
  token: localStorage.getItem("token"),
  profile: null,

  setToken: (token) => {
    if (token) localStorage.setItem("token", token);
    else localStorage.removeItem("token");
    set({ token });
  },
  setState: (key) => (value) => set({ [key]: value }),
}));
