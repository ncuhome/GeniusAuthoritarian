import { useEventListener } from "@hooks/useEventListener";

export const useKeyDown = (key: string | null, callback: () => void) => {
  useEventListener(key ? "keydown" : null, (e: KeyboardEvent) => {
    if (e.key === key) {
      callback();
      e.preventDefault();
    }
  });
};
export default useKeyDown;
