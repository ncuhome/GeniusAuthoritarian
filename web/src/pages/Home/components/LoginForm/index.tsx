import { FC } from "react";
import { useNavigate } from "react-router-dom";
import { useQuery, useMount } from "@hooks";
import toast from "react-hot-toast";
import feishuLogo from "@/assets/img/login/feishu.png";
import dingLogo from "@/assets/img/login/ding.png";

import { ThrowError } from "@util/nav";

import { Stack, Box, Typography, List, Paper } from "@mui/material";
import { LoginItem } from "./components";

import {GetFeishuLoginUrl, GetDingTalkLoginUrl} from "@api/v1/login";

export const LoginForm: FC = () => {
  const nav = useNavigate();
    const [target] = useQuery("target", "");

    async function goFeishuLogin() {
        try {
            const url = await GetFeishuLoginUrl(target);
            window.open(url, "_self");
        } catch ({msg}) {
            if (msg) toast.error(msg as string);
        }
    }

    async function goDingTalkLogin() {
        try {
            const url = await GetDingTalkLoginUrl(target);
            window.open(url, "_self");
        } catch ({msg}) {
            if (msg) toast.error(msg as string);
        }
    }

    useMount(() => {
        if (!target) {
            ThrowError(nav, "请求不合法");
            return;
        }

        switch (true) {
            case navigator.userAgent.indexOf("Feishu") !== -1:
                goFeishuLogin();
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
        <Typography
          variant={"h5"}
          sx={{
            marginBottom: "2rem",
          }}
        >
          登录
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
