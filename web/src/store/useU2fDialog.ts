import { create } from "zustand";
import { shallow } from "zustand/shallow";
import { apiV1User } from "@api/v1/user/base";

interface U2fDialog extends User.U2F.Status {
  open: boolean;
  tabValue: User.U2F.Methods;
  tip?: string;

  resolve?: (value: User.U2F.Result | PromiseLike<User.U2F.Result>) => void;
  reject?: (reason?: any) => void;

  u2f?: User.U2F.Result;

  openDialog: (tip?: string) => Promise<User.U2F.Result>;
  closeDialog: () => void;
  setStatus: (status: User.U2F.Status) => void;
  setToken: (token: User.U2F.Result) => void;
  setTabValue: (value: User.U2F.Methods) => void;
  resetTabValue: () => void;
}

export const useU2fDialog = create<U2fDialog>()((set, getState) => ({
  prefer: "",
  phone: false,
  mfa: false,
  passkey: false,

  open: false,
  tabValue: "",

  openDialog: (tip?: string) => {
    return new Promise<User.U2F.Result>((resolve, reject) => {
      set({ open: true, tip, resolve, reject });
    });
  },
  closeDialog: () => set({ open: false }),
  setStatus: (target) => {
    if (!shallow<User.U2F.Status>(target, getState())) {
      set(target);
      getState().resetTabValue();
    }
  },
  setToken: (u2f) => set({ u2f }),
  setTabValue: (tabValue) => set({ tabValue }),
  resetTabValue: () => {
    const states = getState();
    let tabValue: User.U2F.Methods | undefined = undefined;
    if (states.prefer != "" && states[states.prefer]) tabValue = states.prefer;
    else if (states.passkey) tabValue = "passkey";
    else if (states.mfa) tabValue = "mfa";
    else if (states.phone) tabValue = "phone";
    if (tabValue !== undefined) set({ tabValue });
  },
}));

export default useU2fDialog;
