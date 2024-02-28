import { SWRConfiguration } from "swr";
import { useUserApiV1 } from "@api/v1/user/hook";

export const useLoginData = (
  range: Admin.DataFetchRange,
  config?: SWRConfiguration<Admin.LoginRecordAdminView[]>,
) => {
  return useUserApiV1<Admin.LoginRecordAdminView[]>(
    `admin/data/login?${new URLSearchParams({
      range: range,
    }).toString()}`,
    config,
  );
};
export default useLoginData;
