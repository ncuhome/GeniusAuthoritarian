import { apiV1User } from "@api/v1/user/base";
import { Group } from "@api/v1/user/group";

export type App = {
  id: number;
  name: string;
  appCode: string;
  permitAllGroup: boolean;
  groups: Group[];
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
  permitGroups?: number[]
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

export async function DeleteApp(id: number): Promise<void> {
  await apiV1User.delete("app/", {
    params: {
      id,
    },
  });
}
