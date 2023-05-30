import { FC } from "react";
import { createUseQuery, useMount } from "@hooks";
import { useNavigate } from "react-router-dom";
import toast from "react-hot-toast";
import feishuLogo from "@/assets/img/login/feishu.png";
import dingLogo from "@/assets/img/login/ding.png";
import { ThrowError } from "@util/nav";

import { Stack, Box, Typography, List, Paper, Skeleton } from "@mui/material";
import { LoginItem } from "./components";

import { ErrNetwork, apiV1 } from "@api/base";
import { useApiV1WithLoading } from "@api/v1/hook";

import { useUser } from "@store";

export const LoginForm: FC = () => {
  const nav = useNavigate();
  const useQuery = createUseQuery();

  const [appCode] = useQuery("appCode", "");

  const token = useUser((state) => state.token);

  const { data: appInfo } = useApiV1WithLoading<App.LoginInfo>(
    `public/app/?appCode=${appCode}`,
    {
      onError(err) {
        if (err.msg !== ErrNetwork) {
          ThrowError(nav, "登录对象异常", err.msg);
        }
      },
    }
  );

  async function onGoLogin(thirdParty: string) {
    try {
      const {
        data: {
          data: { url },
        },
      } = await apiV1.get(`public/login/${thirdParty}/link/${appCode}`);
      window.open(url, "_self");
    } catch ({ msg }) {
      if (msg) toast.error(msg as string);
    }
  }

  const onGoFeishuLogin = () => onGoLogin("feishu");
  const onGoDingTalkLogin = () => onGoLogin("dingTalk");

  useMount(() => {
    if (token && !appCode) nav("/user");
    switch (true) {
      case navigator.userAgent.indexOf("Feishu") !== -1:
        onGoFeishuLogin();
        break;
    }
  });

  return (
    <Box
      sx={{
        width: "25rem",
        maxWidth: "100%",
        overflowY: "auto",
        padding: "2rem 3rem",
        borderRadius: "0.4rem",
      }}
      component={Paper}
      elevation={5}
    >
      <Stack
        sx={{
          minWidth: "100%",
          textAlign: "center",
        }}
        justifyContent={"center"}
      >
        <Stack
          alignItems={"center"}
          sx={{
            marginBottom: "2rem",
          }}
        >
          {appInfo ? (
            <>
              <Stack
                flexDirection={"row"}
                sx={{
                  "&>h6": {
                    mx: "0.2rem",
                  },
                }}
              >
                <Typography
                  variant={"h6"}
                  fontWeight={"bold"}
                  letterSpacing={"1px"}
                >
                  {appInfo.name}
                </Typography>
                <Typography variant={"h6"}>{` 请求登录你的帐号`}</Typography>
              </Stack>
              <Typography
                variant={"subtitle1"}
                sx={{ opacity: 0.6 }}
              >{`aka. ${appInfo.host}`}</Typography>
            </>
          ) : (
            <>
              <Skeleton
                variant={"text"}
                sx={{
                  fontSize: "1.25rem",
                  lineHeight: 1.6,
                  fontWeight: 400,
                  width: "75%",
                }}
              />
              <Skeleton
                variant={"text"}
                sx={{
                  fontSize: "1rem",
                  lineHeight: 1.75,
                  fontWeight: 400,
                  width: "30%",
                }}
              />
            </>
          )}
        </Stack>

        <Typography variant={"body2"} sx={{ opacity: 0.6, textAlign: "left" }}>
          请选择你的登录方式:
        </Typography>
        <List>
          <LoginItem
            logo={feishuLogo}
            text={"飞书"}
            onClick={onGoFeishuLogin}
          />
          <LoginItem
            logo={dingLogo}
            text={"钉钉"}
            onClick={onGoDingTalkLogin}
            disableDivider
          />
        </List>
      </Stack>
    </Box>
  );
};
export default LoginForm;
