import { FC, useState } from "react";
import { useInterval } from "@hooks/useInterval";
import toast from "react-hot-toast";

import {
  Chip,
  Stack,
  StackProps,
  Dialog,
  DialogContent,
  DialogActions,
  Button,
  Stepper,
  Step,
  StepLabel,
  Box,
  TextField,
  Typography,
} from "@mui/material";
import { LoadingButton } from "@mui/lab";
import { Done, Remove } from "@mui/icons-material";

import { apiV1User } from "@api/v1/user/base";

import useMfaCode from "@hooks/useMfaCode";
import { shallow } from "zustand/shallow";
import useNewMfaForm from "@store/useNewMfa";

interface Props extends StackProps {
  enabled: boolean;
  setEnabled: (enabled: boolean) => void;
}

export const Mfa: FC<Props> = ({ enabled, setEnabled, ...props }) => {
  const onMfaCode = useMfaCode();

  const [newMfaStep, newMfaSmsCode, newMfaCode] = useNewMfaForm(
    (state) => [state.step, state.smsCode, state.mfaCode],
    shallow
  );
  const setNewMfaStep = useNewMfaForm((state) => state.setState("step"));
  const setNewMfaSmsCode = useNewMfaForm((state) => state.setState("smsCode"));
  const setNewMfaCode = useNewMfaForm((state) => state.setState("mfaCode"));
  const resetNewMfaForm = useNewMfaForm((state) => state.reset);

  const [newMfaNextStepLoading, setNewMfaNextStepLoading] = useState(false);

  const [isSendingSms, setIsSendingSms] = useState(false);
  const [smsCoolDown, setSmsCoolDown] = useState(0);
  useInterval(
    () => setSmsCoolDown((num) => num - 1),
    smsCoolDown > 0 ? 1000 : null
  );

  const [showNewMfa, setShowNewMfa] = useState(false);
  const [mfaNew, setMfaNew] = useState<User.Mfa.New | null>(null);

  const [isCloseMfaLoading, setIsCloseMfaLoading] = useState(false);

  async function onEnableMfa() {
    resetNewMfaForm();
    setShowNewMfa(true);
  }

  async function onApplyNewMfa(smsCode: string, nextStep: number) {
    setNewMfaNextStepLoading(true);
    try {
      const {
        data: { data },
      } = await apiV1User.get("mfa/", {
        params: {
          code: smsCode,
        },
      });
      setMfaNew(data);
      setNewMfaStep(nextStep);
    } catch ({ msg }) {
      if (msg) toast.error(msg as string);
    }
    setNewMfaNextStepLoading(false);
  }

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

  async function onCheckMfaEnable(code: string) {
    try {
      await apiV1User.post("mfa/", {
        code,
      });
      setEnabled(true);
      setShowNewMfa(false);
      toast.success("已启用双因素认证");
    } catch ({ msg }) {
      if (msg) toast.error(msg as string);
    }
  }

  async function onDisableMfa() {
    setIsCloseMfaLoading(true)
    try {
      const code = await onMfaCode();
      try {
        await apiV1User.delete("mfa/", {
          params: {
            code,
          },
        });
        setEnabled(false);
        toast.success("已关闭双因素认证");
      } catch ({ msg }) {
        if (msg) toast.error(msg as string);
      }
    } catch (err) {

    }
    setIsCloseMfaLoading(false)
  }

  function renderNewMfaStep(step: number) {
    switch (step) {
      case 0:
        return (
          <>
            <Typography variant={"h6"} marginBottom={"1rem"}>
              短信身份校验
            </Typography>
            <Stack flexDirection={"row"}>
              <TextField
                variant={"outlined"}
                sx={{ width: "60%" }}
                inputProps={{
                  style: {
                    height: "1rem",
                  },
                }}
                value={newMfaSmsCode}
                onChange={(e) => setNewMfaSmsCode(e.target.value)}
              />
              <Stack
                flexDirection={"row"}
                flexGrow={1}
                paddingX={"1.3rem"}
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
          </>
        );
      case 1:
        return (
          <Stack
            alignItems={"center"}
            mb={2.5}
            sx={{
              width: "21rem",
              maxWidth: "100%",
            }}
          >
            <img
              alt={"totp qrcode"}
              src={mfaNew?.qrcode}
              style={{
                width: "17rem",
                maxWidth: "100%",
                boxSizing: "border-box",
                marginBottom: "1rem",
              }}
            />

            <Button
              onClick={async () => {
                try {
                  await navigator.clipboard.writeText(mfaNew?.url || "");
                  toast.success("已复制");
                } catch (err) {
                  toast.error(`复制失败: ${err}`);
                }
              }}
            >
              复制 totp url
            </Button>
          </Stack>
        );
      case 2:
        return (
          <>
            <Typography variant={"h6"} marginBottom={"1rem"}>
              双因素校验码
            </Typography>
            <TextField
              variant={"outlined"}
              value={newMfaCode}
              onChange={(e) => setNewMfaCode(e.target.value)}
            />
          </>
        );
    }
  }

  async function onNextNewMfaStep(step: number) {
    switch (step) {
      case 0:
        if (newMfaSmsCode == "") {
          toast.error("请输入短信校验码");
          return;
        }
        await onApplyNewMfa(newMfaSmsCode, step + 1);
        break;
      case 1:
        setNewMfaStep(step + 1);
        break;
      case 2:
        if (newMfaSmsCode == "") {
          toast.error("请输入双因素校验码");
          return;
        } else if (newMfaCode.length != 6 && !isNaN(Number(newMfaSmsCode))) {
          toast.error("双因素校验码错误");
          return;
        }
        await onCheckMfaEnable(newMfaCode);
        break;
    }
  }

  return (
    <>
      <Stack flexDirection={"row"} alignItems={"center"} {...props}>
        <Chip
          label={enabled ? "双因素认证已开启" : "双因素未启用"}
          variant={"outlined"}
          icon={
            enabled ? <Done color={"success"} fontSize="small" /> : <Remove />
          }
        />

        <Box
          sx={{
            ml: 2,
          }}
        >
          {enabled ? (
            <>
              <LoadingButton
                variant={"outlined"}
                color={"warning"}
                loading={isCloseMfaLoading}
                onClick={onDisableMfa}
              >
                关闭
              </LoadingButton>
            </>
          ) : (
            <>
              <Button variant={"outlined"} onClick={onEnableMfa}>
                开启
              </Button>
            </>
          )}
        </Box>
      </Stack>

      <Dialog open={showNewMfa} onClose={() => setShowNewMfa(false)}>
        <DialogContent>
          <Stack>
            <Stepper activeStep={newMfaStep}>
              <Step>
                <StepLabel>身份验证</StepLabel>
              </Step>
              <Step>
                <StepLabel>双因素校验绑定</StepLabel>
              </Step>
              <Step>
                <StepLabel>完成绑定</StepLabel>
              </Step>
            </Stepper>
            <Box
              sx={{
                marginTop: "1.5rem",
                paddingX: "1rem",
                boxSizing: "border-box",
              }}
            >
              <Stack>{renderNewMfaStep(newMfaStep)}</Stack>
            </Box>
          </Stack>
        </DialogContent>
        <DialogActions>
          <Button onClick={() => setShowNewMfa(false)}>取消</Button>
          <LoadingButton
            loading={newMfaNextStepLoading}
            onClick={() => onNextNewMfaStep(newMfaStep)}
          >
            {newMfaStep == 2 ? "完成" : "下一步"}
          </LoadingButton>
        </DialogActions>
      </Dialog>
    </>
  );
};
export default Mfa;
