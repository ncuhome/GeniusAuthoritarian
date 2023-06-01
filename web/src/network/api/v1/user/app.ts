import { apiV1User } from "@api/v1/user/base";

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
