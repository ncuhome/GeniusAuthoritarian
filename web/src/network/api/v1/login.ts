import { apiV1 } from "@api/base";

export async function LoginUrl(
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
