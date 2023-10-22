import { useEffect } from "react";
import { useLoadingToast } from "@hooks/useLoadingToast";
import useSWR, { SWRConfiguration } from "swr";
import { AxiosInstance } from "axios";

export const createFetchHook = (api: AxiosInstance) => {
  const fetcher = (url: string) => api.get(url).then((res) => res.data.data);
  return <T>(
    url: string | null,
    config?: SWRConfiguration<T> & {
      enableLoading?: boolean;
      immutable?: boolean;
    }
  ) => {
    if (config?.immutable) {
      config.revalidateIfStale = false;
      config.revalidateOnFocus = false;
      config.revalidateOnReconnect = false;
    }

    const swr = useSWR<T, ApiError<ApiResponse<T>>>(url, fetcher, config);

    if (config?.enableLoading) {
      const [showToast, closeToast] = useLoadingToast();
      useEffect(() => {
        if (swr.error?.msg) showToast(swr.error.msg);
        else closeToast();
        return closeToast;
      }, [swr.error]);
    }
    return swr;
  };
};
