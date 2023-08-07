import { create } from "zustand";

type DevRouteState = {
  index: number;

  setState: <T extends keyof DevRouteState>(
    key: T
  ) => (value: DevRouteState[T]) => void;
};

export const useDevRoute = create<DevRouteState>()((set) => ({
  index: 0,
  setState: (key) => (value) => set({ [key]: value }),
}));
export default useDevRoute;
