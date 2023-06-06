import { FC } from "react";
import { useNavigate } from "react-router-dom";
import { createUseQuery, useMount } from "@hooks";
import { ThrowError } from "@util/nav";

import { OnLogin } from "@components";

import { apiV1 } from "@api/base";

import { useUser } from "@store";

// 用户登录用户中心处理
export const Login: FC = () => {
  const nav = useNavigate();
  const useQuery = createUseQuery();
  const [token] = useQuery("token", "");
  const setAuth = useUser((state) => state.setAuth);

  async function handleLogin() {
    try {
      const {
        data: { data: res },
      } = await apiV1.post<{
        data: User.Login.Result;
      }>("public/login/", {
        token,
      });
      setAuth(res.token, res.groups);
      nav("/user/");
    } catch ({ msg }) {
      if (msg) ThrowError(nav, "登录失败", msg as string);
    }
  }

  useMount(() => {
    if (token == "") {
      ThrowError(nav, "登录失败", "参数缺失");
      return;
    }

    handleLogin();
  });

  return <OnLogin />;
};
