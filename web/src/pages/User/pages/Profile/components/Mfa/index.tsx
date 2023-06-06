import { FC, useEffect, useRef, useState } from "react";
import toast from "react-hot-toast";

import {
  Chip,
  Stack,
  StackProps,
  Dialog,
  DialogTitle,
  DialogContent,
  DialogActions,
  ButtonGroup,
  Button,
  TextField,
} from "@mui/material";
import { LoadingButton } from "@mui/lab";
import { Done, Remove } from "@mui/icons-material";

import { apiV1User } from "@api/v1/user/base";

interface Props extends StackProps {
  enabled: boolean;
  setEnabled: (enabled: boolean) => void;
}

export const Mfa: FC<Props> = ({ enabled, setEnabled, ...props }) => {
  const [checkMfaCallback, setCheckMfaCallback] = useState<
    ((code: string) => void) | null
  >(null);
  const [checkMfaCode, setCheckMfaCode] = useState("");
  const checkMfaInputEl = useRef<HTMLInputElement | null>(null);
  const [isMfaCheckLoading, setIsMfaCheckLoading] = useState(false);

  const [showNewMfa, setShowNewMfa] = useState(false);
  const [mfaNew, setMfaNew] = useState<User.Mfa.New | null>(null);
  const [isLoadingMfaNew, setIsLoadingMfaNew] = useState(false);

  async function onEnableMfa() {
    setIsLoadingMfaNew(true);
    try {
      const {
        data: { data },
      } = await apiV1User.get("mfa/");
      setMfaNew(data);
      setShowNewMfa(true);
    } catch ({ msg }) {
      if (msg) toast.error(msg as string);
    }
    setIsLoadingMfaNew(false);
  }

  async function onCheckMfaEnable(code: string) {
    setIsMfaCheckLoading(true);
    try {
      await apiV1User.post("mfa/", {
        code,
      });
      setCheckMfaCallback(null);
      setEnabled(true);
      toast.success("已启用双因素认证");
    } catch ({ msg }) {
      if (msg) toast.error(msg as string);
    }
    setIsMfaCheckLoading(false);
  }

  async function onDisableMfa(code: string) {
    setIsMfaCheckLoading(true);
    try {
      await apiV1User.delete("mfa/", {
        params: {
          code,
        },
      });
      setCheckMfaCallback(null);
      setEnabled(false);
      toast.success("已关闭双因素认证");
    } catch ({ msg }) {
      if (msg) toast.error(msg as string);
    }
    setIsMfaCheckLoading(false);
  }

  useEffect(() => {
    if (checkMfaCallback) setCheckMfaCode("");
  }, [checkMfaCallback]);

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

        <ButtonGroup
          sx={{
            ml: 2,
          }}
        >
          {enabled ? (
            <>
              <Button
                color={"warning"}
                onClick={() => setCheckMfaCallback(() => onDisableMfa)}
              >
                关闭
              </Button>
            </>
          ) : (
            <>
              <LoadingButton
                variant={"outlined"}
                loading={isLoadingMfaNew}
                onClick={onEnableMfa}
              >
                开启
              </LoadingButton>
            </>
          )}
        </ButtonGroup>
      </Stack>

      <Dialog open={showNewMfa} onClose={() => setShowNewMfa(false)}>
        <DialogTitle>双因素认证</DialogTitle>
        <DialogContent>
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
        </DialogContent>
        <DialogActions>
          <Button onClick={() => setShowNewMfa(false)}>取消</Button>
          <Button
            onClick={() => {
              setShowNewMfa(false);
              setCheckMfaCallback(() => onCheckMfaEnable);
            }}
          >
            下一步
          </Button>
        </DialogActions>
      </Dialog>

      <Dialog
        open={Boolean(checkMfaCallback)}
        onClose={() => setCheckMfaCallback(null)}
      >
        <DialogTitle>双因素认证校验</DialogTitle>
        <DialogContent
          sx={{
            width: "20rem",
          }}
        >
          <TextField
            autoFocus
            fullWidth
            margin="dense"
            inputRef={checkMfaInputEl}
            label={"校验码"}
            name={"twofactor_token"}
            value={checkMfaCode}
            onChange={(e) => {
              if (!Number(e.target.value) && e.target.value != "") return;
              setCheckMfaCode(e.target.value);
            }}
          />
        </DialogContent>
        <DialogActions>
          <Button onClick={() => setCheckMfaCallback(null)}>取消</Button>
          <LoadingButton
            loading={isMfaCheckLoading}
            onClick={() => {
              if (!checkMfaCode) {
                toast.error("请输入校验码");
                checkMfaInputEl.current?.focus();
                return;
              }
              if (checkMfaCode.length != 6) {
                toast.error("校验码错误");
                checkMfaInputEl.current?.focus();
                return;
              }
              checkMfaCallback!(checkMfaCode);
            }}
          >
            确认
          </LoadingButton>
        </DialogActions>
      </Dialog>
    </>
  );
};
export default Mfa;
