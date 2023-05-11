import { create } from "zustand";
import { UserProfile } from "@api/v1/user/profile";

type DialogProps = {
  title: string;
  content?: string;
  callback: (ok: boolean) => void;
};

interface UserState {
  token: string | null;
  groups: string[];

  openDialog: boolean;
  dialog: DialogProps;

  profile: UserProfile | null;

  setAuth: (token: string | null, groups?: string[]) => void;
  setDialog: (props: DialogProps) => void;

  setState: <T extends keyof UserState>(
    key: T
  ) => (value: UserState[T]) => void;
}

export const useUser = create<UserState>()((set) => ({
  token: localStorage.getItem("token"),
  groups: localStorage.getItem("groups")?.split(",") || [],

  openDialog: false,
  dialog: { title: "", callback: () => null },

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
  setDialog: (props) => {
    const callback = props.callback;
    props.callback = (ok) => {
      callback(ok);
      set({ openDialog: false });
    };
    set({
      dialog: props,
      openDialog: true,
    });
  },

  setState: (key) => (value) => set({ [key]: value }),
}));
