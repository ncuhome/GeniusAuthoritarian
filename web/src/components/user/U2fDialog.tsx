import { FC, useEffect, useRef, useState } from "react";
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
} from "@mui/material";
import { LoadingButton } from "@mui/lab";
import { ReportProblemOutlined } from "@mui/icons-material";

import { apiV1User } from "@api/v1/user/base";

import { shallow } from "zustand/shallow";
import useU2fDialog from "@store/useU2fDialog";

const U2fDialog: FC = () => {
  const open = useU2fDialog((state) => state.open);
  const u2fStatus = useU2fDialog(
    (state) => ({
      prefer: state.prefer,
      mfa: state.mfa,
      phone: state.phone,
      passkey: state.passkey,
    }),
    shallow
  );

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

  const onCancel = () => {
    const states = useU2fDialog.getState();
    states.closeDialog();
    if (states.reject) states.reject("user canceled");
  };

  const onSubmit = async () => {
    setIsLoading(true);
    let data: any;
    switch (tabValue) {
      case "phone":
        if (smsCode === "") {
          toast.error("验证码不能为空");
          return;
        }
        if (smsCode.length != 6) {
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
        //todo
        break;
    }
    try {
      const {
        data: { data: result },
      } = await apiV1User.post<{ data: User.U2F.Result }>(
        `u2f/${tabValue}`,
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
          <Stack flexDirection={"row"}>
            <TextField
              variant={"outlined"}
              label={"验证码"}
              sx={{ width: "60%" }}
              value={smsCode}
              onChange={(e) => {
                if (!isNaN(Number(e.target.value))) setSmsCode(e.target.value);
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
        );
      case "mfa":
        return (
            <Stack alignItems={"center"}>
              <TextField
                  variant={"outlined"}
                  label={"校验码"}
                  value={mfaCode}
                  onChange={(e) => {
                    if (!isNaN(Number(e.target.value))) setMfaCode(e.target.value);
                  }}
              />
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
    <Dialog open={open} onClose={onCancel}>
      <DialogTitle>U2F 身份校验</DialogTitle>
      <DialogContent>
        <DialogContentText>
          你正在进行敏感操作，需要额外的身份校验
        </DialogContentText>

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
      </DialogContent>
      <DialogActions>
        <Button onClick={onCancel}>取消</Button>
        <LoadingButton
          disabled={tabValue === ""}
          loading={isLoading}
          onClick={onSubmit}
        >
          确定
        </LoadingButton>
      </DialogActions>
    </Dialog>
  );
};
export default U2fDialog;
