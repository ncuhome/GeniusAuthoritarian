import { apiV1 } from "@api/base";

export async function LoginUrl(
  thirdParty: string,
  code: string,
  appCode: string
): Promise<User.ThirdPartyLoginResult> {
  const {
    data: { data },
  } = await apiV1.post(`public/login/${thirdParty}/${appCode}`, {
    code,
  });
  return data;
}
