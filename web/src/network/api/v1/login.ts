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

export async function FeishuLogin(code: string, appCode: string): Promise<string> {
  const {
    data: {
      data: { token, callback: callbackUrl },
    },
  } = await apiV1.post("public/login/feishu/", {
    code,
    appCode,
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

export async function DingTalkLogin(code: string, appCode: string): Promise<string> {
  const {
    data: {
      data: { token, callback: callbackUrl },
    },
  } = await apiV1.post("public/login/dingTalk/", {
    code,
    appCode,
  });
  return callbackUrl;
}
