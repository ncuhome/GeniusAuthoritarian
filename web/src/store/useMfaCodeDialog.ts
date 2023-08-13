import { create } from "zustand";

type MfaCodeDialog = {
  description?: string;
  callback: ((code: string | null) => void) | null;

  setDialog: (callback: (code: string | null) => void, desc?: string) => void;
  resetDialog: () => void;
};

const useMfaCodeDialog = create<MfaCodeDialog>()((set) => ({
  callback: null,

  setDialog: (callback, desc) => {
    set({
      description: desc,
      callback: callback,
    });
  },
  resetDialog: () => {
    set({ callback: null });
  },
}));
export default useMfaCodeDialog;
