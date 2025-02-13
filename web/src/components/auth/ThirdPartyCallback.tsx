import { memo, useRef, useState } from "react";
import { createUseQuery } from "@hooks/useQuery";
import useMount from "@hooks/useMount";
import { useNavigate } from "react-router";
import { ThrowError } from "@util/nav";
import toast from "react-hot-toast";

import LoginLoading from "@components/auth/LoginLoading";
import { Stack, Typography, TextField } from "@mui/material";
import { LoadingButton } from "@mui/lab";

import { apiV1 } from "@api/base";

interface Props {
  keyCode: string;
  keyAppCode: string;
  thirdParty: string;
}

export const ThirdPartyCallback = memo<Props>(
  ({ keyCode, keyAppCode, thirdParty }) => {
    const nav = useNavigate();
    const useQuery = createUseQuery();
    const [code] = useQuery(keyCode, "");
    const [appCode] = useQuery(keyAppCode, "");

    const [mfaToken, setMfaToken] = useState("");
    const [mfaCode, setMfaCode] = useState("");
    const [mfaLoading, setMfaLoading] = useState(false);
    const mfaInput = useRef<HTMLInputElement | null>(null);

    async function login() {
      try {
        const {
          data: { data },
        } = await apiV1.post<{ data: User.Login.ThirdParty }>(
          `public/login/${thirdParty}/${appCode}`,
          {
            code,
          },
        );
        if (!data.mfa) window.open(data.callback!, "_self");
        else setMfaToken(data.token);
      } catch (err) {
        if (err instanceof Error && err.message)
          ThrowError(nav, "登录失败", err.message, appCode);
      }
    }

    async function mfa() {
      if (!mfaCode) {
        toast.error("双因素校验码不能为空");
        mfaInput.current?.focus();
        return;
      }
      if (mfaCode.length != 6) {
        toast.error("双因素校验码错误");
        mfaInput.current?.focus();
        return;
      }

      setMfaLoading(true);
      try {
        const {
          data: { data },
        } = await apiV1.post<{ data: User.Login.Verified }>(
          "public/login/mfa",
          {
            token: mfaToken,
            code: mfaCode,
          },
        );
        window.open(data.callback, "_self");
      } catch (err) {
        if (err instanceof Error) toast.error(err.message);
      }
      setMfaLoading(false);
    }

    useMount(() => {
      if (!code) {
        ThrowError(nav, "登录失败", "参数缺失", appCode);
        return;
      }
      login();
    });

    return mfaToken ? (
      <Stack
        justifyContent={"center"}
        alignItems={"center"}
        sx={{
          height: "100%",
        }}
      >
        <Stack
          spacing={2.5}
          sx={{
            width: "20rem",
            maxWidth: "85%",
          }}
        >
          <Typography variant={"h5"} fontWeight={"bold"} letterSpacing={"2px"}>
            双因素认证
          </Typography>
          <TextField
            label={"校验码"}
            name={"twofactor_token"}
            fullWidth
            value={mfaCode}
            onChange={(e) => {
              if (isNaN(Number(e.target.value)) && e.target.value != "") return;
              setMfaCode(e.target.value);
            }}
            inputRef={mfaInput}
          />
          <LoadingButton
            variant={"outlined"}
            loading={mfaLoading}
            onClick={mfa}
            sx={{
              py: 1,
            }}
          >
            确认
          </LoadingButton>
        </Stack>
      </Stack>
    ) : (
      <LoginLoading />
    );
  },
);
export default ThirdPartyCallback;
