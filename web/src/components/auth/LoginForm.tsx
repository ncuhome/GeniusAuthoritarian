import { FC } from "react";
import { createUseQuery } from "@hooks/useQuery";
import useMount from "@hooks/useMount";
import { useNavigate } from "react-router";
import toast from "react-hot-toast";
import { ThrowError } from "@util/nav";
import {
  coerceResponseToBase64Url,
  coerceToArrayBuffer,
  coerceToBase64Url,
} from "@util/coerce";

import feishuLogo from "@/assets/img/login_methods/feishu.png";
import webpFeishuLogo from "@/assets/img/login_methods/feishu.webp";
import dingLogo from "@/assets/img/login_methods/ding.png";
import webpDingLogo from "@/assets/img/login_methods/ding.webp";
import passkeyLogo from "@/assets/img/login_methods/passkeys.svg";

import LoginItem from "@components/auth/LoginItem";
import { Stack, Box, Typography, List, Paper, Skeleton } from "@mui/material";

import { AxiosError } from "axios";
import { ErrNetwork, apiV1 } from "@api/base";
import { useApiV1 } from "@api/v1/hook";

import useUser from "@store/useUser";

export const LoginForm: FC = () => {
  const nav = useNavigate();
  const useQuery = createUseQuery();

  const [appCode] = useQuery("appCode", "");

  const token = useUser((state) => state.token);

  const { data: appInfo } = useApiV1<App.LoginInfo>(
    `public/app/?appCode=${appCode}`,
    {
      immutable: true,
      enableLoading: true,
      onError(err) {
        if (err.msg !== ErrNetwork) {
          ThrowError(nav, "登录对象异常", err.msg, appCode);
        }
      },
    },
  );

  async function onGoLogin(thirdParty: string) {
    try {
      const {
        data: {
          data: { url },
        },
      } = await apiV1.get(`public/login/${thirdParty}/link/${appCode}`);
      window.open(url, "_self");
    } catch (err) {
      if (err instanceof Error) toast.error(err.message);
    }
  }

  const onGoFeishuLogin = () => onGoLogin("feishu");
  const onGoDingTalkLogin = () => onGoLogin("dingTalk");
  const onPasskeyLogin = async () => {
    try {
      const {
        data: { data: options },
      } = await apiV1.get("public/login/passkey/");
      options.publicKey.challenge = coerceToArrayBuffer(
        options.publicKey.challenge,
      );
      const credential = await navigator.credentials.get(options);
      if (credential?.type !== "public-key") {
        toast.error(`获取凭据失败，凭据类型不正确: ${credential?.type}`);
        return;
      }
      const pubKeyCred = credential as any;
      const {
        data: { data },
      } = await apiV1.post<{ data: User.Login.Verified }>(
        "public/login/passkey/",
        {
          app_code: appCode,
          credential: {
            id: pubKeyCred.id,
            rawId: coerceToBase64Url(pubKeyCred.rawId),
            response: coerceResponseToBase64Url(pubKeyCred.response),
            type: pubKeyCred.type,
          },
        },
      );
      window.open(data.callback, "_self");
    } catch (err: any) {
      if (err instanceof AxiosError) {
        err = err as ApiError<void>;
        if (err.msg) toast.error(err.msg);
      } else {
        if (err.name != "NotAllowedError") toast.error(`创建凭据失败: ${err}`);
      }
    }
  };

  useMount(() => {
    if (token && !appCode) nav("/user", { replace: true });
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
                  sx={{
                    textWrap: "nowrap",
                    "&,&>b": {
                      letterSpacing: "1px",
                    },
                    "&>b": {
                      textWrap: "balance",
                    },
                  }}
                >
                  <b>{appInfo.name}</b>
                  {` 请求登录你的帐号`}
                </Typography>
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
            webpLogo={webpFeishuLogo}
            text={"飞书"}
            onClick={onGoFeishuLogin}
          />
          <LoginItem
            logo={dingLogo}
            webpLogo={webpDingLogo}
            text={"钉钉"}
            onClick={onGoDingTalkLogin}
          />
          <LoginItem
            logo={passkeyLogo}
            text={"通行密钥"}
            onClick={onPasskeyLogin}
            sx={{
              "& img": {
                filter: "invert(1)",
              },
            }}
            disableDivider
          />
        </List>
      </Stack>
    </Box>
  );
};
export default LoginForm;
