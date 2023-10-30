import { create } from "zustand";
import { apiV1User } from "@api/v1/user/base";

interface U2fDialog extends User.U2F.Status {
  open: boolean;
  tip?: string;

  resolve?: (value: User.U2F.Result | PromiseLike<User.U2F.Result>) => void;
  reject?: (reason?: any) => void;

  u2f?: User.U2F.Result;

  openDialog: (tip?: string) => Promise<User.U2F.Result>;
  closeDialog: () => void;
  setStatus: (status: User.U2F.Status) => void;
  setToken: (token: User.U2F.Result) => void;
  setPrefer: (method: User.U2F.Methods) => Promise<void>;
}

export const useU2fDialog = create<U2fDialog>()((set, getState) => ({
  prefer: "",
  phone: false,
  mfa: false,
  passkey: false,

  open: false,

  openDialog: (tip?: string) => {
    return new Promise<User.U2F.Result>((resolve, reject) => {
      set({ open: true, tip, resolve, reject });
    });
  },
  closeDialog: () => set({ open: false }),
  setStatus: (status) => set(status),
  setToken: (u2f) => set({ u2f }),
  setPrefer: async (method: User.U2F.Methods) => {
    if (method === getState().prefer) return;
    await apiV1User.put("u2f/prefer", {
      prefer: method,
    });
    set({ prefer: method });
  },
}));

export default useU2fDialog;
