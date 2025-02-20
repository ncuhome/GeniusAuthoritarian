import { ReactNode } from "react";
import { create } from "zustand";

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

  profile: User.Profile | null;
  apps: App.Detailed[] | null;
  accessibleApps: App.Accessible | undefined;

  setAuth: (token: string | null, groups?: string[]) => void;
  setDialog: (props: DialogProps) => Promise<boolean>;
  setDeviceOffline: (id: number) => void;
  setProfile: (profile: User.Profile) => void;
  setApps: (apps: App.Detailed[]) => void;
  setAccessibleApps: (apps: App.Accessible) => void;
}

export const useUser = create<UserState>()((set, get) => ({
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
  setDeviceOffline: (id) => {
    const profile = get().profile;
    if (!profile) return;
    set({
      profile: {
        ...profile,
        loginRecord: {
          ...profile.loginRecord,
          online: profile.loginRecord.online.filter((item) => item.id != id),
        },
      },
    });
  },
  setProfile: (profile) => {
    set({ profile });
  },
  setApps: (apps) => {
    set({ apps });
  },
  setAccessibleApps: (accessibleApps) => {
    set({ accessibleApps });
  },
}));
export default useUser;
