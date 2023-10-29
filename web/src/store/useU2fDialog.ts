import { create } from "zustand";

interface U2fDialog extends User.U2F.Status {
  open: boolean;

  resolve?: (value: User.U2F.Result | PromiseLike<User.U2F.Result>) => void;
  reject?: (reason?: any) => void;

  openDialog: () => Promise<User.U2F.Result>;
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
}));

export default useU2fDialog;
