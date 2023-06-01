import { create } from "zustand";

type GroupState = {
  groups?: User.Group[];

  setState: <T extends keyof GroupState>(
    key: T
  ) => (value: GroupState[T]) => void;
};

export const useGroup = create<GroupState>()((set) => ({
  setState: (key) => (value) => set({ [key]: value }),
}));
