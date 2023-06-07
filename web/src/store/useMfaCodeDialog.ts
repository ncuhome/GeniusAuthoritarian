import { create } from "zustand";

type MfaCodeDialog = {
  callback: ((code: string) => Promise<void>) | null;

  setState: <T extends keyof MfaCodeDialog>(
    key: T
  ) => (value: MfaCodeDialog[T]) => void;
};

const useMfaCodeDialog = create<MfaCodeDialog>()((set) => ({
  callback: null,

  setState: (key) => (value) => set({ [key]: value }),
}));
export default useMfaCodeDialog
