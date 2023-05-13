import { apiV1User } from "@api/v1/user/base";

export type Group = {
  id: number;
  name: string;
};

export async function ListGroups(): Promise<Group[]> {
  const {
    data: { data },
  } = await apiV1User.get("group/list");
  return data;
}
