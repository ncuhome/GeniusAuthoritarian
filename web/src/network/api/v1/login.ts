import { apiV1 } from "@api/base";

export async function UserLogin(token: string): Promise<string> {
  const {
    data: {
      data: { token: authToken },
    },
  } = await apiV1.post("public/login/", {
    token,
  });
  return authToken;
}

export async function GetFeishuLoginUrl(appCode: string): Promise<string> {
  const {
    data: {
      data: { url },
    },
  } = await apiV1.get("public/login/feishu/link", {
    params: {
      appCode,
    },
  });
  return url;
}

export async function FeishuLogin(
  code: string,
  callback: string
): Promise<string> {
  const {
    data: {
      data: { token, callback: callbackUrl },
    },
  } = await apiV1.post("public/login/feishu/", {
    code,
    callback,
  });
  return callbackUrl;
}

export async function GetDingTalkLoginUrl(appCode: string): Promise<string> {
  const {
    data: {
      data: { url },
    },
  } = await apiV1.get("public/login/dingTalk/link", {
    params: {
      appCode,
    },
  });
  return url;
}

export async function DingTalkLogin(
  code: string,
  callback: string
): Promise<string> {
  const {
    data: {
      data: { token, callback: callbackUrl },
    },
  } = await apiV1.post("public/login/dingTalk/", {
    code,
    callback,
  });
  return callbackUrl;
}
