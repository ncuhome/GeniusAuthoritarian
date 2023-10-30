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
  const submitHandler = useRef(() => {});

  const [smsCode, setSmsCode] = useState("");
  const [isSendingSms, setIsSendingSms] = useState(false);
  const [smsCoolDown, setSmsCoolDown] = useState(0);
  useInterval(
    () => setSmsCoolDown((num) => num - 1),
    smsCoolDown > 0 ? 1000 : null
  );

  const [mfaCode, setMfaCode] = useState("");

  const onSwitchTab = (target: User.U2F.Methods) => {
    switch (target) {
      case "phone":
        submitHandler.current = () => {
          if (smsCode.length != 6) {
            toast.error("短信验证码有误");
            return;
          }
          return onSubmit("phone", { code: smsCode });
        };
        break;
      case "mfa":
        //todo
        break;
      case "passkey":
        //todo
        break;
    }
    setTabValue(target);
  };

  const onCancel = () => {
    const states = useU2fDialog.getState();
    states.closeDialog();
    if (states.reject) states.reject("user canceled");
  };

  const onSubmit = async (method: User.U2F.Methods, data: any) => {
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

  const renderTabPanel = () => {
    switch (tabValue) {
      case "phone":
        return (
          <Stack flexDirection={"row"} mt={3}>
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
      default:
        return (
          <Stack alignItems={"center"} spacing={1}>
            <ReportProblemOutlined fontSize={"large"} color={"warning"} />
            <Typography>没有可用的 U2F 身份验证方法</Typography>
          </Stack>
        );
    }
  };

  useEffect(() => {
    if (open) {
      setMfaCode("");
      setSmsCode("");
    }
  }, [open]);
  useEffect(() => {
    if (u2fStatus.prefer != "" && u2fStatus[u2fStatus.prefer])
      onSwitchTab(u2fStatus.prefer);
    else if (u2fStatus.passkey) onSwitchTab("passkey");
    else if (u2fStatus.mfa) onSwitchTab("mfa");
    else if (u2fStatus.phone) onSwitchTab("phone");
  }, [u2fStatus]);

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
            onSwitchTab(value as User.U2F.Methods)
          }
          variant="fullWidth"
          sx={{
            my: 2,
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
        <LoadingButton
          disabled={tabValue === ""}
          loading={isLoading}
          onClick={() => submitHandler.current()}
        >
          确定
        </LoadingButton>
        <Button onClick={onCancel}>取消</Button>
      </DialogActions>
    </Dialog>
  );
};
export default U2fDialog;
