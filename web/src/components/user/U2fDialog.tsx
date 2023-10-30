import { FC, useEffect, useState } from "react";
import { useInterval } from "@hooks/useInterval";
import toast from "react-hot-toast";

import {
  Dialog,
  DialogActions,
  DialogTitle,
  DialogContent,
  DialogContentText,
  Button,
  Tabs,
  Tab,
  Typography,
  Stack,
  TextField,
  Alert,
  AlertTitle,
} from "@mui/material";
import { LoadingButton } from "@mui/lab";
import {
  ReportProblemOutlined,
  SensorOccupiedOutlined,
} from "@mui/icons-material";

import { apiV1User } from "@api/v1/user/base";

import { shallow } from "zustand/shallow";
import useU2fDialog from "@store/useU2fDialog";
import {
  coerceResponseToBase64Url,
  coerceToArrayBuffer,
  coerceToBase64Url,
} from "@util/coerce";

const U2fDialog: FC = () => {
  const open = useU2fDialog((state) => state.open);
  const tip = useU2fDialog((state) => state.tip);
  const u2fStatus = useU2fDialog(
    (state) => ({
      prefer: state.prefer,
      mfa: state.mfa,
      phone: state.phone,
      passkey: state.passkey,
    }),
    shallow
  );
  const u2fToken = useU2fDialog((state) => state.u2f);

  const [tokenAvailable, setTokenAvailable] = useState(false);

  const [tabValue, setTabValue] = useState<User.U2F.Methods>("");
  const [isLoading, setIsLoading] = useState(false);

  const [smsCode, setSmsCode] = useState("");
  const [isSendingSms, setIsSendingSms] = useState(false);
  const [smsCoolDown, setSmsCoolDown] = useState(0);
  useInterval(
    () => setSmsCoolDown((num) => num - 1),
    smsCoolDown > 0 ? 1000 : null
  );

  const [mfaCode, setMfaCode] = useState("");

  const isTokenAvailable = () => {
    return !!u2fToken && u2fToken.valid_before > Date.now() / 1000;
  };
  useInterval(
    () => {
      if (!isTokenAvailable()) setTokenAvailable(false);
    },
    tokenAvailable ? 1000 : null
  );

  const onCancel = () => {
    const states = useU2fDialog.getState();
    states.closeDialog();
    if (states.reject) states.reject("user canceled");
  };

  const onSubmit = async (method: string = tabValue) => {
    if (tokenAvailable) {
      const states = useU2fDialog.getState();
      if (states.resolve) states.resolve(u2fToken!);
      states.closeDialog();
      return;
    }
    let data: any;
    switch (method) {
      case "phone":
        if (smsCode === "") {
          toast.error("验证码不能为空");
          return;
        }
        if (smsCode.length != 5) {
          toast.error("短信验证码有误");
          return;
        }
        data = { code: smsCode };
        break;
      case "mfa":
        if (mfaCode === "") {
          toast.error("校验码不能为空");
          return;
        }
        if (mfaCode.length != 6) {
          toast.error("校验码有误");
          return;
        }
        data = { code: mfaCode };
        break;
      case "passkey":
        const {
          data: { data: options },
        } = await apiV1User.get("passkey/options");
        options.publicKey.challenge = coerceToArrayBuffer(
          options.publicKey.challenge
        );
        options.publicKey.allowCredentials =
          options.publicKey.allowCredentials.map((cred: any) => {
            cred.id = coerceToArrayBuffer(cred.id);
            return cred;
          });
        const credential = await navigator.credentials.get(options);
        if (!(credential instanceof PublicKeyCredential)) {
          toast.error(`获取凭据失败，凭据类型不正确`);
          return;
        }
        const pubKeyCred = credential as any;
        data = {
          id: pubKeyCred.id,
          rawId: coerceToBase64Url(pubKeyCred.rawId),
          response: coerceResponseToBase64Url(pubKeyCred.response),
          type: pubKeyCred.type,
        };
        break;
    }
    setIsLoading(true);
    try {
      const {
        data: { data: result },
      } = await apiV1User.post<{ data: User.U2F.Result }>(
        `u2f/${method}`,
        data
      );
      const states = useU2fDialog.getState();
      states.setToken(result);
      if (states.resolve) states.resolve(result);
      states.closeDialog();
    } catch ({ msg }) {
      if (msg) toast.error(msg as any);
    }
    setIsLoading(false);
  };

  async function onSendSmsCode() {
    setIsSendingSms(true);
    try {
      await apiV1User.post("identity/sms");
      setSmsCoolDown(60);
    } catch ({ msg }) {
      if (msg) toast.error(msg as string);
    }
    setIsSendingSms(false);
  }

  useEffect(() => {
    if (open) {
      setMfaCode("");
      setSmsCode("");
      if (isTokenAvailable()) setTokenAvailable(true);
      else if (tabValue === "passkey") onSubmit();
    }
  }, [open]);
  useEffect(() => {
    if (u2fStatus.prefer != "" && u2fStatus[u2fStatus.prefer])
      setTabValue(u2fStatus.prefer);
    else if (u2fStatus.passkey) setTabValue("passkey");
    else if (u2fStatus.mfa) setTabValue("mfa");
    else if (u2fStatus.phone) setTabValue("phone");
  }, [u2fStatus]);

  const renderTabPanel = () => {
    switch (tabValue) {
      case "phone":
        return (
          <Stack alignItems={"center"}>
            <Typography mb={2.5}>
              每天上限五条，出现异常情况请过会儿再试
            </Typography>
            <Stack flexDirection={"row"} width={"100%"} maxWidth={"21rem"}>
              <TextField
                variant={"outlined"}
                label={"验证码"}
                sx={{ width: "60%" }}
                value={smsCode}
                onChange={(e) => {
                  if (!isNaN(Number(e.target.value)))
                    setSmsCode(e.target.value);
                }}
              />
              <Stack
                flexDirection={"row"}
                flexGrow={1}
                paddingLeft={"5%"}
                boxSizing={"border-box"}
              >
                <LoadingButton
                  variant={"contained"}
                  fullWidth
                  disabled={!!smsCoolDown}
                  loading={isSendingSms}
                  onClick={onSendSmsCode}
                >
                  {smsCoolDown ? smsCoolDown + "s" : "发送"}
                </LoadingButton>
              </Stack>
            </Stack>
          </Stack>
        );
      case "mfa":
        return (
          <Stack alignItems={"center"}>
            <TextField
              variant={"outlined"}
              label={"校验码"}
              name={"twofactor_token"}
              value={mfaCode}
              onChange={(e) => {
                if (!isNaN(Number(e.target.value))) setMfaCode(e.target.value);
              }}
            />
          </Stack>
        );
      case "passkey":
        return (
          <Stack alignItems={"center"}>
            <SensorOccupiedOutlined
              sx={{
                fontSize: "6rem",
                mt: 2,
                mb: 4,
              }}
            />
            <Button variant={"outlined"} onClick={() => onSubmit()}>
              重试
            </Button>
          </Stack>
        );
      default:
        return (
          <Stack alignItems={"center"} spacing={1}>
            <ReportProblemOutlined fontSize={"large"} color={"warning"} />
            <Typography>没有可用的 U2F 身份验证方法</Typography>
          </Stack>
        );
    }
  };

  return (
    <Dialog fullWidth open={open} onClose={onCancel}>
      <DialogTitle>U2F 身份校验</DialogTitle>
      <DialogContent>
        {!tip && tokenAvailable ? undefined : (
          <DialogContentText>
            {tip ? tip : "你正在进行敏感操作，需要额外的身份校验"}
          </DialogContentText>
        )}

        {tokenAvailable ? (
          <Alert severity="success">
            <AlertTitle>已认证</AlertTitle>
            最近 5 分钟已通过验证，无需再次校验。你可以通过刷新提前移除校验状态
          </Alert>
        ) : (
          <>
            <Tabs
              value={tabValue}
              onChange={(e, value: string) => {
                setTabValue(value as User.U2F.Methods);
                if (value === "passkey") return onSubmit("passkey");
              }}
              variant="fullWidth"
              sx={{
                mt: 2,
                mb: 3,
              }}
            >
              <Tab label={"短信"} value={"phone"} disabled={!u2fStatus.phone} />
              <Tab label={"MFA"} value={"mfa"} disabled={!u2fStatus.mfa} />
              <Tab
                label={"通行密钥"}
                value={"passkey"}
                disabled={!u2fStatus.passkey}
              />
            </Tabs>
            {renderTabPanel()}
          </>
        )}
      </DialogContent>
      <DialogActions>
        <Button onClick={onCancel}>取消</Button>
        <LoadingButton
          disabled={tabValue === ""}
          loading={isLoading}
          onClick={() => onSubmit()}
        >
          确定
        </LoadingButton>
      </DialogActions>
    </Dialog>
  );
};
export default U2fDialog;
