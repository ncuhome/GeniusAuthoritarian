import { apiV1User } from "@api/v1/user/base";

export type App = {
  id: number;
  name: string;
  appCode: string;
  permitAllGroup: boolean;
};

export async function GetOwnedAppList(): Promise<App[]> {
  const {
    data: { data },
  } = await apiV1User.get("app/");
  return data;
}

export type AppNew = {
  appSecret: string;
} & App;

export async function ApplyApp(
  name: string,
  callback: string,
  permitAll: boolean,
  permitGroups?: string[]
): Promise<AppNew> {
  const {
    data: { data },
  } = await apiV1User.post("app/", {
    name,
    callback,
    permitAll,
    permitGroups,
  });
  return data;
}
