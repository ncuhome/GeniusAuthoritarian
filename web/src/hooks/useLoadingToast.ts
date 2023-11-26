import { useRef } from "react";
import { useMount } from "./useMount";
import toast, { ToastOptions } from "react-hot-toast";

type fnShowToast = (msg: string, options?: ToastOptions) => void;

export function useLoadingToast(): [
  fnShowToast,
  (msg?: string, success?: boolean) => void
];
export function useLoadingToast(): [fnShowToast, () => void];
export function useLoadingToast() {
  const id = useRef<string | null>(null);

  const showToast = (msg: string, options?: ToastOptions) => {
    if (!id.current) id.current = toast.loading(msg, options);
    else toast.loading(msg, { id: id.current, ...options });
  };

  const closeToast = (msg?: string, success: boolean = true) => {
    if (!id.current) return;
    if (!msg) toast.dismiss(id.current);
    else if (success) toast.success(msg, { id: id.current });
    else toast.error(msg, { id: id.current });
    id.current = null;
  };

  useMount(() => {
    return () => closeToast();
  });

  return [showToast, closeToast];
}

export default useLoadingToast;
