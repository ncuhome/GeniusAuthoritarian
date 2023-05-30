import useSWR, { SWRConfiguration } from "swr";
import { createSwrWithLoading } from "@hooks";
import { apiV1User } from "./base";

export const useUserApiV1 = <T>(url: string, config?: SWRConfiguration<T>) =>
  useSWR<T>(
    url,
    (url) => apiV1User.get(url).then((res) => res.data.data),
    config
  );

export const useUserApiV1WithLoading = createSwrWithLoading(useUserApiV1);
