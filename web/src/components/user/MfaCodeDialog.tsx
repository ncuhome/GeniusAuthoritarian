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

import useMfaCodeDialog from "@store/useMfaCodeDialog";

export const MfaCodeDialog: FC = () => {
  const callback = useMfaCodeDialog((state) => state.callback);
  const resetDialog = useMfaCodeDialog((state) => state.resetDialog);

  const [code, setCode] = useState("");
  const inputEl = useRef<HTMLInputElement | null>(null);

  const handleCancel = () => {
    if(callback) callback(null)
    resetDialog()
  }

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
    callback!(code);
  }

  useEffect(() => {
    if (callback) {
      setCode("");
    }
  }, [callback]);

  return (
    <Dialog
      open={Boolean(callback)}
      onAnimationStart={() => {
        if (callback) inputEl.current?.focus();
      }}
      onClose={handleCancel}
    >
      <DialogTitle>双因素认证校验</DialogTitle>
      <DialogContent>
        <TextField
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
          sx={{
            width: "16rem",
            maxWidth: "100%",
          }}
        />
      </DialogContent>
      <DialogActions>
        <Button onClick={handleCancel}>取消</Button>
        <Button onClick={handleSubmit}>确认</Button>
      </DialogActions>
    </Dialog>
  );
};
export default MfaCodeDialog;
