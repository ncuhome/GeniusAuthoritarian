import { useEffect, useRef } from "react";

type EventListener<K extends keyof DocumentEventMap> = (
  ev: DocumentEventMap[K]
) => any;
export const useEventListener = <K extends keyof DocumentEventMap>(
  event: K | null,
  callback: EventListener<K>,
  options?: boolean | AddEventListenerOptions
) => {
  const savedCallback = useRef<EventListener<K>>();
  useEffect(() => {
    savedCallback.current = callback;
  }, [callback]);
  useEffect(() => {
    if (event != null) {
      const handler: EventListener<K> = (ev) => {
        savedCallback.current!(ev);
      };
      document.addEventListener(event, handler, options);
      return () => {
        document.removeEventListener(event, handler);
      };
    }
  }, [event]);
};
