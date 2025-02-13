import { create } from "zustand";

type DevRouteState = {
  index: number;
};

export const useDevRoute = create<DevRouteState>()((set) => ({
  index: 0,
}));
export default useDevRoute;
