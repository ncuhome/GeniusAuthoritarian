import { useRef } from "react";
import toast, { ToastOptions } from "react-hot-toast";

export function useLoadingToast() {
  const id = useRef<string | null>(null);

  const showToast = (msg: string, options?: ToastOptions) => {
    if (id.current) {
      toast.loading(msg, { id: id.current });
    } else {
      id.current = toast.loading(msg, options);
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
