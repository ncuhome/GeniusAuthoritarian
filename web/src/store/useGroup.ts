import { create } from "zustand";

type GroupState = {
  groups?: User.Group[];

  setGroups: (groups: User.Group[]) => void;
};

export const useGroup = create<GroupState>()((set) => ({
  setGroups: (groups) => {
    set({ groups });
  },
}));
export default useGroup;
