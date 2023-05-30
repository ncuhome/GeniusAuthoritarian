import { ReactNode } from "react";
import { create } from "zustand";
import { UserProfile } from "@api/v1/user/profile";
import { AppDetailed } from "@api/v1/user/app";

type DialogProps = {
  title: string;
  content?: ReactNode;
};

interface UserState {
  token: string | null;
  groups: string[];

  openDialog: boolean;
  dialog: DialogProps;
  dialogResolver?: (ok: boolean) => void;

  profile: UserProfile | null;
  apps: AppDetailed[] | null;
  accessibleApps: App.Accessible | undefined;

  setAuth: (token: string | null, groups?: string[]) => void;
  setDialog: (props: DialogProps) => Promise<boolean>;

  setState: <T extends keyof UserState>(
    key: T
  ) => (value: UserState[T]) => void;
}

export const useUser = create<UserState>()((set) => ({
  token: localStorage.getItem("token"),
  groups: localStorage.getItem("groups")?.split(",") || [],

  openDialog: false,
  dialog: { title: "" },

  profile: null,
  apps: null,
  accessibleApps: undefined,

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
  setDialog: async (props): Promise<boolean> => {
    return new Promise((resolve) => {
      set({
        dialog: props,
        dialogResolver: (ok: boolean) => {
          resolve(ok);
          set({ openDialog: false, dialogResolver: undefined });
        },
        openDialog: true,
      });
    });
  },

  setState: (key) => (value) => set({ [key]: value }),
}));
