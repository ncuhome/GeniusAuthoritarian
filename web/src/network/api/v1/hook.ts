import useSWR, { SWRConfiguration } from "swr";
import { createSwrWithLoading } from "@hooks";
import { apiV1 } from "@api/base";

export const useApiV1 = <T>(url: string, config?: SWRConfiguration<T>) =>
  useSWR<T>(url, (url) => apiV1.get(url).then((res) => res.data.data), config);

export const useApiV1WithLoading = createSwrWithLoading(useApiV1);
