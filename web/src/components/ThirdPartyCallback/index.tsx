import { FC } from "react";
import { useMount, createUseQuery } from "@hooks";
import { useNavigate } from "react-router-dom";
import { ThrowError } from "@util/nav";

import { OnLogin } from "@components";

import { LoginUrl } from "@api/v1/login";

interface Props {
  keyCode: string;
  keyAppCode: string;
  thirdParty: string;
}

export const ThirdPartyCallback: FC<Props> = ({
  keyCode,
  keyAppCode,
  thirdParty,
}) => {
  const nav = useNavigate();
  const useQuery = createUseQuery();
  const [code] = useQuery(keyCode, "");
  const [appCode] = useQuery(keyAppCode, "");

  async function login() {
    try {
      const data = await LoginUrl(thirdParty, code, appCode);
      if (!data.mfa) window.open(data.callback, "_self");
      else {
        //todo 双因素认证
      }
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
export default ThirdPartyCallback
