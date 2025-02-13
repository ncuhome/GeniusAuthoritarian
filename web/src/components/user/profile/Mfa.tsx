import { FC, useState } from "react";
import useInterval from "@hooks/useInterval";
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
  CircularProgress,
} from "@mui/material";
import { LoadingButton } from "@mui/lab";
import { Done, Remove } from "@mui/icons-material";

import { apiV1User } from "@api/v1/user/base";

import useU2F from "@hooks/data/useU2F";
import useNewMfaForm from "@store/useNewMfa";

interface Props extends StackProps {
  enabled: boolean;
  setEnabled: (enabled: boolean) => void;
}

export const Mfa: FC<Props> = ({ enabled, setEnabled, ...props }) => {
  const { isLoading: isU2fLoading, refreshToken } = useU2F();

  const newMfaStep = useNewMfaForm((state) => state.step);
  const newMfaCode = useNewMfaForm((state) => state.mfaCode);
  const setNewMfaStep = (step: number) => {
    useNewMfaForm.setState({ step });
  };
  const setNewMfaCode = (code: string) => {
    useNewMfaForm.setState({ mfaCode: code });
  };
  const resetNewMfaForm = useNewMfaForm((state) => state.reset);

  const [newMfaNextStepLoading, setNewMfaNextStepLoading] = useState(false);

  const [smsCoolDown, setSmsCoolDown] = useState(0);
  useInterval(
    () => setSmsCoolDown((num) => num - 1),
    smsCoolDown > 0 ? 1000 : null,
  );

  const [showNewMfa, setShowNewMfa] = useState(false);
  const [mfaNew, setMfaNew] = useState<User.Mfa.New | null>(null);

  async function onEnableMfa() {
    resetNewMfaForm();
    setShowNewMfa(true);
    setNewMfaNextStepLoading(true);
    try {
      const token = await refreshToken();
      const {
        data: { data },
      } = await apiV1User.get("mfa/", {
        headers: {
          Authorization: token,
        },
      });
      setMfaNew(data);
      setNewMfaStep(1);
    } catch (err) {
      if (err instanceof Error) toast.error(err.message);
      setShowNewMfa(false);
    }
    setNewMfaNextStepLoading(false);
  }

  async function onCheckMfaEnable(code: string) {
    try {
      await apiV1User.post("mfa/", {
        code,
      });
      setEnabled(true);
      setShowNewMfa(false);
      toast.success("已启用双因素认证");
    } catch (err) {
      if (err instanceof Error) toast.error(err.message);
    }
  }

  async function onDisableMfa() {
    try {
      const token = await refreshToken();
      try {
        await apiV1User.delete("mfa/", {
          headers: {
            Authorization: token,
          },
        });
        setEnabled(false);
        toast.success("已关闭双因素认证");
      } catch (err) {
        if (err instanceof Error) toast.error(err.message);
      }
    } catch (err) {}
  }

  function renderNewMfaStep(step: number) {
    switch (step) {
      case 0:
        return (
          <Stack alignItems={"center"} mt={2.5}>
            <CircularProgress />
          </Stack>
        );
      case 1:
        return (
          <Stack
            alignItems={"center"}
            sx={{
              mt: 3,
              width: "100%",
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
          <Stack alignItems={"center"} mt={1}>
            <Box>
              <Typography variant={"h6"} marginBottom={"1rem"}>
                双因素校验码
              </Typography>
              <TextField
                variant={"outlined"}
                name={"twofactor_token"}
                value={newMfaCode}
                onChange={(e) => setNewMfaCode(e.target.value)}
              />
            </Box>
          </Stack>
        );
    }
  }

  async function onNextNewMfaStep(step: number) {
    switch (step) {
      case 0:
        return;
      case 1:
        setNewMfaStep(step + 1);
        break;
      case 2:
        if (newMfaCode == "") {
          toast.error("请输入双因素校验码");
          return;
        } else if (newMfaCode.length != 6 && !isNaN(Number(newMfaCode))) {
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
                loading={isU2fLoading}
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

      <Dialog open={showNewMfa} fullWidth onClose={() => setShowNewMfa(false)}>
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
