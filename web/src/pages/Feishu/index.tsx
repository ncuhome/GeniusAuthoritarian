import { FC } from "react";
import { useMount, createUseQuery } from "@hooks";
import { useNavigate } from "react-router-dom";
import { ThrowError } from "@util/nav";

import { OnLogin } from "@components";

import { LoginUrl } from "@api/v1/login";

export const Feishu: FC = () => {
  const nav = useNavigate();
  const useQuery = createUseQuery();
  const [code] = useQuery("code", "");
  const [appCode] = useQuery("state", "");

  async function login() {
    try {
      const callbackUrl = await LoginUrl("feishu", code, appCode);
      window.open(callbackUrl, "_self");
    } catch ({ msg }) {
      if (msg) ThrowError(nav, "登录失败", msg as string);
    }
  }

  useMount(() => {
    if (!code) {
      ThrowError(nav, "登录失败", "参数缺失");
      return;
    }
    login();
  });

  return <OnLogin />;
};
export default Feishu;
