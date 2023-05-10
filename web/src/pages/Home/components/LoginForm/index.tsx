import { FC, useState } from "react";
import { createUseQuery, useMount, useInterval } from "@hooks";
import { useNavigate } from "react-router-dom";
import toast from "react-hot-toast";
import feishuLogo from "@/assets/img/login/feishu.png";
import dingLogo from "@/assets/img/login/ding.png";
import { ThrowError } from "@util/nav";

import { Stack, Box, Typography, List, Paper, Skeleton } from "@mui/material";
import { LoginItem } from "./components";

import { ErrNetwork } from "@api/base";
import { GetLoginUrl } from "@api/v1/login";
import { GetAppInfo, AppInfo } from "@api/v1/app";

import { useUser } from "@store";

export const LoginForm: FC = () => {
  const nav = useNavigate();
  const useQuery = createUseQuery();

  const [appCode] = useQuery("appCode", "");

  const token = useUser((state) => state.token);

  const [appInfo, setAppInfo] = useState<AppInfo | null>(null);
  const [onRequestAppInfo, setOnRequestAppInfo] = useState(true);

  async function goLogin(thirdParty: string) {
    try {
      const url = await GetLoginUrl(thirdParty, appCode);
      window.open(url, "_self");
    } catch ({ msg }) {
      if (msg) toast.error(msg as string);
    }
  }

  const goFeishuLogin = () => goLogin("feishu");
  const goDingTalkLogin = () => goLogin("dingTalk");

  async function loadAppInfo() {
    setOnRequestAppInfo(true);
    try {
      const data = await GetAppInfo(appCode);
      setAppInfo(data);
    } catch ({ msg }) {
      if (msg) {
        if (msg === ErrNetwork) {
          toast.error(msg);
        } else {
          ThrowError(nav, "登录对象异常", msg as string);
        }
      }
    }
    setOnRequestAppInfo(false);
  }

  useInterval(loadAppInfo, !appInfo && !onRequestAppInfo ? 2000 : null);

  useMount(() => {
    if (token && !appCode) nav("/user");
    switch (true) {
      case navigator.userAgent.indexOf("Feishu") !== -1:
        goFeishuLogin();
        break;
    }

    loadAppInfo();
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
              <Typography variant={"h6"}>
                <Typography
                  variant={"h6"}
                  display={"inline"}
                  fontWeight={"bold"}
                  letterSpacing={"1px"}
                >
                  {appInfo.name}
                </Typography>
                {` 请求登录你的帐号`}
              </Typography>
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
          <LoginItem logo={feishuLogo} text={"飞书"} onClick={goFeishuLogin} />
          <LoginItem
            logo={dingLogo}
            text={"钉钉"}
            onClick={goDingTalkLogin}
            disableDivider
          />
        </List>
      </Stack>
    </Box>
  );
};
export default LoginForm;
