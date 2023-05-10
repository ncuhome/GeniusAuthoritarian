import { apiV1 } from "@api/base";

export type AppInfo = {
  name: string;
  host: string
  createdAt: number;
};

export async function GetAppInfo(appCode: string): Promise<AppInfo> {
  const {
    data: { data },
  } = await apiV1.get("public/app/", {
    params: {
      appCode,
    },
  });
  return data;
}