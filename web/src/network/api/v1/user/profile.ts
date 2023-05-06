import {apiV1User} from "./base";

export type UserProfile = {
  user: {
    id: number;
    name: string;
    phone: string;
  };
  loginRecord: Array<{
    id: number;
    createdAt: number;
    target: string;
    ip: string;
  }>;
};

export async function GetUserProfile(): Promise<UserProfile> {
  const {
    data: {data},
  } = await apiV1User.get("profile/");
  return data;
}
