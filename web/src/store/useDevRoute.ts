import { create } from "zustand";

type DevRouteState = {
  index: number;

  setIndex: (index: number) => void;
};

export const useDevRoute = create<DevRouteState>()((set) => ({
  index: 0,

  setIndex: (index: number) => {
    set({ index });
  },
}));
export default useDevRoute;
