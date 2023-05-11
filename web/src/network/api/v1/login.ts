import { apiV1 } from "@api/base";

export async function GetLoginUrl(
  thirdParty: string,
  appCode: string = ""
): Promise<string> {
  const {
    data: {
      data: { url },
    },
  } = await apiV1.get(`public/login/${thirdParty}/link/${appCode}`);
  return url;
}

export type UserLoginResult = {
  token: string;
  groups: string[];
};

export async function UserLogin(token: string): Promise<UserLoginResult> {
  const {
    data: { data },
  } = await apiV1.post("public/login/", {
    token,
  });
  return data;
}

export async function Login(
  thirdParty: string,
  code: string,
  appCode: string
): Promise<string> {
  const {
    data: {
      data: { token, callback: callbackUrl },
    },
  } = await apiV1.post(`public/login/${thirdParty}/${appCode}`, {
    code,
  });
  return callbackUrl;
}
