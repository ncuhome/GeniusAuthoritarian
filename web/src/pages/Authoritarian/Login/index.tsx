import { FC } from "react";
import { useNavigate } from "react-router";
import { createUseQuery } from "@hooks/useQuery";
import useMount from "@hooks/useMount";
import { ThrowError } from "@util/nav";

import LoginLoading from "@components/auth/LoginLoading";

import { apiV1 } from "@api/base";

import useUser from "@store/useUser";
import toast from "react-hot-toast";

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
      nav("/user/", { replace: true });
    } catch (err) {
      if (err instanceof Error) toast.error(err.message);
    }
  }

  useMount(() => {
    if (token == "") {
      ThrowError(nav, "登录失败", "参数缺失", "");
      return;
    }

    handleLogin();
  });

  return <LoginLoading />;
};
export default Login;
