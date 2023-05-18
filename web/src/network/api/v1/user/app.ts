import { apiV1User } from "@api/v1/user/base";
import { Group } from "@api/v1/user/group";

export type App = {
  id: number;
  name: string;
  callback: string;
  permitAllGroup: boolean;
};

export type AppOwned = {
  appCode: string;
} & App;

export type AppDetailed = {
  groups: Group[];
} & AppOwned;

export async function GetOwnedAppList(): Promise<AppDetailed[]> {
  const {
    data: { data },
  } = await apiV1User.get("app/");
  return data;
}

export type AccessibleApps = {
  permitAll: App[];
  accessible: {
    group: Group;
    app: App[];
  }[];
};

export async function GetAccessibleAppList(): Promise<AccessibleApps> {
  const {
    data: { data },
  } = await apiV1User.get("app/accessible");
  return data;
}

export type AppNew = {
  appSecret: string;
} & AppDetailed;

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

export async function ModifyApp(
  id: number,
  name: string,
  callback: string,
  permitAll: boolean,
  permitGroups?: number[]
): Promise<void> {
  await apiV1User.put("app/", {
    id,
    name,
    callback,
    permitAll,
    permitGroups,
  });
}

export async function GetLandingUrl(id: number): Promise<string> {
  const {
    data: {
      data: { url },
    },
  } = await apiV1User.get("app/landing", {
    params: {
      id,
    },
  });
  return url;
}
