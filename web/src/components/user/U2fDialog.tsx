import { FC, useCallback, useEffect, useState } from "react";
import useInterval from "@hooks/useInterval";
import useTimeout from "@hooks/useTimeout";
import useKeyDown from "@hooks/useKeyDown";
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
  IconButton,
  CircularProgress,
  Box,
} from "@mui/material";
import { LoadingButton } from "@mui/lab";
import {
  ReportProblemOutlined,
  SensorOccupiedOutlined,
} from "@mui/icons-material";

import { apiV1User } from "@api/v1/user/base";

import { useShallow } from "zustand/react/shallow";
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
    useShallow((state) => ({
      prefer: state.prefer,
      mfa: state.mfa,
      phone: state.phone,
      passkey: state.passkey,
    })),
  );
  const u2fToken = useU2fDialog((state) => state.u2f);
  const tabValue = useU2fDialog((state) => state.tabValue);
  const setTabValue = useU2fDialog((state) => state.setTabValue);

  const [tokenAvailable, setTokenAvailable] = useState(false);

  const [isLoading, setIsLoading] = useState(false);
  const [autoConfirm, setAutoConfirm] = useState(3);
  useTimeout(() => setAutoConfirm(3), autoConfirm != 3 && !open ? 100 : null);
  useInterval(
    () => {
      setAutoConfirm((num) => {
        const next = num - 0.1;
        if (next <= 0) onSubmit();
        return next;
      });
    },
    autoConfirm > 0 && tokenAvailable && open ? 100 : null,
  );

  const [smsCode, setSmsCode] = useState("");
  const [isSendingSms, setIsSendingSms] = useState(false);
  const [smsCoolDown, setSmsCoolDown] = useState(0);
  useInterval(
    () => setSmsCoolDown((num) => num - 1),
    smsCoolDown > 0 ? 1000 : null,
  );

  const [mfaCode, setMfaCode] = useState("");

  const isTokenAvailable = useCallback(() => {
    const u2fTokenRef = useU2fDialog.getState().u2f;
    return !!u2fTokenRef && u2fTokenRef.valid_before > Date.now() / 1000;
  }, []);
  useInterval(
    () => {
      if (!isTokenAvailable()) setTokenAvailable(false);
    },
    tokenAvailable ? 1000 : null,
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
          options.publicKey.challenge,
        );
        options.publicKey.allowCredentials =
          options.publicKey.allowCredentials.map((cred: any) => {
            cred.id = coerceToArrayBuffer(cred.id);
            return cred;
          });
        const credential = await navigator.credentials.get(options);
        if (credential?.type !== "public-key") {
          toast.error(`获取凭据失败，凭据类型不正确: ${credential?.type}`);
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
        data,
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
  useKeyDown("Enter", onSubmit);

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
      if (isTokenAvailable()) {
        setAutoConfirm(3);
        setTokenAvailable(true);
      }
    }
  }, [open]);
  useEffect(() => {
    if (open && !isTokenAvailable() && tabValue === "passkey")
      onSubmit("passkey");
  }, [open, tabValue]);

  const renderTabPanel = () => {
    switch (tabValue) {
      case "phone":
        return (
          <Stack alignItems={"center"}>
            <Stack
              sx={{
                width: "100%",
                maxWidth: "21rem",
              }}
            >
              <Alert
                variant={"outlined"}
                severity={"info"}
                sx={{ mb: 2.5, width: "100%", boxSizing: "border-box" }}
              >
                每天上限五条，请避免频繁重试
              </Alert>
              <Stack flexDirection={"row"} width={"100%"}>
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
                    variant={"outlined"}
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
          </Stack>
        );
      case "mfa":
        return (
          <Stack alignItems={"center"} pt={1.5}>
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
            <IconButton onClick={() => onSubmit()} sx={{ padding: "1.5rem" }}>
              <SensorOccupiedOutlined
                sx={{
                  fontSize: "6rem",
                }}
              />
            </IconButton>
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
        {tokenAvailable ? undefined : (
          <DialogContentText>
            你正在进行敏感操作，需要额外的身份校验
          </DialogContentText>
        )}
        {tip ? (
          <Alert
            severity={"info"}
            sx={{
              mt: 2,
              "& + .MuiAlert-root": {
                mt: 2,
              },
            }}
          >
            {tip}
          </Alert>
        ) : undefined}

        {tokenAvailable ? (
          <Alert
            severity="success"
            icon={
              <Box>
                <Stack position={"relative"}>
                  <CircularProgress
                    thickness={4}
                    size={23.5}
                    color="success"
                    variant="determinate"
                    value={(autoConfirm * 100) / 3}
                  />
                  <Box
                    sx={{
                      top: 0,
                      left: 0,
                      bottom: 0,
                      right: 0,
                      position: "absolute",
                      display: "flex",
                      alignItems: "center",
                      justifyContent: "center",
                    }}
                  >
                    <Typography
                      variant="caption"
                      component="div"
                      color="success"
                    >{`${Math.round(autoConfirm)}`}</Typography>
                  </Box>
                </Stack>
              </Box>
            }
          >
            <AlertTitle>已认证</AlertTitle>
            最近 5 分钟已通过验证，无需再次校验。你可以通过刷新提前移除校验状态
          </Alert>
        ) : (
          <>
            <Tabs
              value={tabValue}
              onChange={(e, value: string) =>
                setTabValue(value as User.U2F.Methods)
              }
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
