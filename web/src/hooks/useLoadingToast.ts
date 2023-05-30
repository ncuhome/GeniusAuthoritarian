import { useEffect, useRef } from "react";
import { useMount } from "./useMount";
import toast, { ToastOptions } from "react-hot-toast";
import { SWRConfiguration, SWRResponse } from "swr";

type fnShowToast = (msg: string, options?: ToastOptions) => void;

export function useLoadingToast(): [
  fnShowToast,
  (msg?: string, success?: boolean) => void
];
export function useLoadingToast(): [fnShowToast, () => void];
export function useLoadingToast() {
  const id = useRef<string | null>(null);

  const showToast = (msg: string, options?: ToastOptions) => {
    if (id.current) toast.loading(msg, { id: id.current });
    else id.current = toast.loading(msg, options);
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

export const createSwrWithLoading =
  (useApi: <T>(url: string, config?: SWRConfiguration<T>) => SWRResponse<T>) =>
  <T>(url: string, config?: SWRConfiguration<T>) => {
    const swr = useApi<T>(url, config);
    const [showToast, closeToast] = useLoadingToast();
    useEffect(() => {
      if (swr.error?.msg) showToast(swr.error.msg);
      else closeToast();
    }, [swr.error]);
    return swr;
  };
