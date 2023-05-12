import { useRef } from "react";
import toast, { Renderable } from "react-hot-toast";

export function useToast() {
  const id = useRef<string | null>(null);

  const showToast = (msg: string, icon?: Renderable) => {
    if (id.current) {
      toast.loading(msg, { id: id.current });
    } else {
      id.current = toast.loading(msg, {
        icon: icon,
      });
    }
  };
  const closeToast = (msg: string, success = true) => {
    if (!id.current) return;
    if (success) toast.success(msg, { id: id.current });
    else toast.error(msg, { id: id.current });
    id.current = null;
  };

  return [showToast, closeToast];
}
