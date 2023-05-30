import useSWR, { SWRConfiguration } from "swr";
import { apiV1 } from "@api/base";
import { useLoadingToast } from "@hooks";
import { useEffect } from "react";

export const useApiV1 = <T>(url: string, config?: SWRConfiguration<T>) =>
  useSWR<T>(url, (url) => apiV1.get(url).then((res) => res.data.data), config);

export const useApiV1WithLoading = <T>(
  url: string,
  config?: SWRConfiguration<T>
) => {
  const swr = useApiV1<T>(url, config);
  const [showToast, closeToast] = useLoadingToast();
  useEffect(() => {
    if (swr.error?.msg) showToast(swr.error.msg);
    else closeToast();
  }, [swr.error]);
  return swr;
};
