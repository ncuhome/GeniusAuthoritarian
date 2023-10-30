import { create } from "zustand";

interface U2fDialog extends User.U2F.Status {
  open: boolean;

  resolve?: (value: User.U2F.Result | PromiseLike<User.U2F.Result>) => void;
  reject?: (reason?: any) => void;

  u2f?: User.U2F.Result;

  openDialog: () => Promise<User.U2F.Result>;
  setStatus: (status: User.U2F.Status) => void;
  setPrefer: (prefer: string) => void;
  setToken: (token: User.U2F.Result) => void;
}

export const useU2fDialog = create<U2fDialog>()((set) => ({
  prefer: "",
  phone: false,
  mfa: false,
  passkey: false,

  open: false,

  openDialog: () => {
    return new Promise<User.U2F.Result>((resolve, reject) => {
      set({ open: true, resolve, reject });
    });
  },
  setStatus: (status) => set(status),
  setPrefer: (prefer: string) => set({ prefer }),
  setToken: (u2f) => set({ u2f }),
}));

export default useU2fDialog;
