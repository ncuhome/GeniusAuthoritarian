import { create } from "zustand";

type AdminDataState = {
  login?: User.LoginRecordAdminView[];

  setState: <T extends keyof AdminDataState>(
    key: T,
  ) => (value: AdminDataState[T]) => void;
};

export const useAdminData = create<AdminDataState>()((set) => ({
  setState: (key) => (value) => set({ [key]: value }),
}));
export default useAdminData;
