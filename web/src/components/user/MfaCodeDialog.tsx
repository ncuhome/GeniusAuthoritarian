import { FC, useEffect, useRef, useState } from "react";
import toast from "react-hot-toast";

import {
  Button,
  Dialog,
  DialogActions,
  DialogContent,
  DialogTitle,
  TextField,
} from "@mui/material";
import { LoadingButton } from "@mui/lab";

import useMfaCodeDialog from "@store/useMfaCodeDialog";

export const MfaCodeDialog: FC = () => {
  const callback = useMfaCodeDialog((state) => state.callback);
  const setCallback = useMfaCodeDialog((state) => state.setState("callback"));

  const [code, setCode] = useState("");
  const [isLoading, setIsLoading] = useState(false);
  const inputEl = useRef<HTMLInputElement | null>(null);

  async function handleSubmit() {
    if (!code) {
      toast.error("请输入校验码");
      inputEl.current?.focus();
      return;
    }
    if (code.length != 6) {
      toast.error("校验码错误");
      inputEl.current?.focus();
      return;
    }
    setIsLoading(true);
    await callback!(code);
    setIsLoading(false);
  }

  useEffect(() => {
    if (callback) setCode("");
  }, [callback]);

  return (
    <Dialog open={Boolean(callback)} onClose={() => setCallback(null)}>
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
          inputRef={inputEl}
          label={"校验码"}
          name={"twofactor_token"}
          value={code}
          onChange={(e) => {
            if (!Number(e.target.value) && e.target.value != "") return;
            setCode(e.target.value);
          }}
        />
      </DialogContent>
      <DialogActions>
        <Button onClick={() => setCallback(null)}>取消</Button>
        <LoadingButton loading={isLoading} onClick={handleSubmit}>
          确认
        </LoadingButton>
      </DialogActions>
    </Dialog>
  );
};
export default MfaCodeDialog;
