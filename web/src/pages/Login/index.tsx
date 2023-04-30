import { FC, useCallback } from "react";
import { useNavigate } from "react-router-dom";
import { createUseQuery, useMount } from "@hooks";
import { ThrowError } from "@util/nav";

import { OnLogin } from "@components";

import { UserLogin } from "@api/v1/login";

export const Login: FC = () => {
  const nav = useNavigate();
  const useQuery = createUseQuery();
  const [token] = useQuery("token", "");

  const handleLogin = useCallback(async () => {
    try {
      const authToken = await UserLogin(token);
      localStorage.setItem("token", authToken);
    } catch ({ msg }) {
      if (msg) ThrowError(nav, "登录失败", msg as string);
    }
  }, [token]);

  useMount(() => {
    if (token == "") {
      ThrowError(nav, "登录失败", "参数缺失");
      return;
    }

    handleLogin();
  });

  return <OnLogin />;
};
