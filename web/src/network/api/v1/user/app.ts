import { apiV1User } from "@api/v1/user/base";

export async function GetOwnedAppList(): Promise<App.Detailed[]> {
  const {
    data: { data },
  } = await apiV1User.get("app/");
  return data;
}

export async function ApplyApp(
  name: string,
  callback: string,
  permitAll: boolean,
  permitGroups?: number[]
): Promise<App.New> {
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
