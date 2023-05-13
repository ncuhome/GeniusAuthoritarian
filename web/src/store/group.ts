import { create } from "zustand";
import { Group } from "@api/v1/user/group";

type GroupState = {
  groups?: Group[];

  setState: <T extends keyof GroupState>(
    key: T
  ) => (value: GroupState[T]) => void;
};

export const useGroup = create<GroupState>()((set) => ({
  setState: (key) => (value) => set({ [key]: value }),
}));
