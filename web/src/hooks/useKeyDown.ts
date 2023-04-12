import { useEffect, useRef } from "react";
import { useEventListener } from ".";

export const useKeyDown = (key: string | null, callback: () => void) => {
  useEventListener(key ? "keydown" : null, (e: KeyboardEvent) =>
    e.key === key ? callback() : null
  );
};
