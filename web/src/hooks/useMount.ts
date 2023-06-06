import { useEffect, EffectCallback } from "react";

export const useMount = (callback: EffectCallback) => {
  useEffect(callback, []);
};
export default useMount;
