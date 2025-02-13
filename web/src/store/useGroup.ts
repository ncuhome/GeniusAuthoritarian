import { create } from "zustand";

type GroupState = {
  groups?: User.Group[];
};

export const useGroup = create<GroupState>()((set) => ({
  groups: undefined,
}));
export default useGroup;
