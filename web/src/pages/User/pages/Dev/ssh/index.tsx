import { FC, useMemo, useState } from "react";
import toast from "react-hot-toast";

import {
  Stack,
  TextField,
  Select,
  MenuItem,
  FormControl,
  InputLabel,
} from "@mui/material";
import { LoadingButton } from "@mui/lab";
import { LockPerson, Link } from "@mui/icons-material";
import Block from "@components/user/Block";

import { useUserApiV1 } from "@api/v1/user/hook";
import { apiV1User } from "@api/v1/user/base";

import useMfaCodeDialog from "@store/useMfaCodeDialog";

const Ssh: FC = () => {
  const setMfaCodeCallback = useMfaCodeDialog((state) =>
    state.setState("callback")
  );

  const [sshKey, setSshKey] = useState<User.SSH.Keys | null>(null);
  const [keyMode, setKeyMode] = useState<"pem" | "ssh">("ssh");
  const [mfaEnabled, setMfaEnabled] = useState<boolean | null>(null);

  const [isUnlockLoading, setIsUnlockLoading] = useState(true);

  useUserApiV1("profile/mfa", {
    immutable: sshKey !== null,
    enableLoading: true,
    onSuccess: (data: any) => {
      setMfaEnabled(data.mfa);
      setIsUnlockLoading(false);
    },
  });

  async function onShowSshKeys(code: string) {
    setIsUnlockLoading(true);
    try {
      const {
        data: { data },
      } = await apiV1User.get("dev/ssh/", {
        params: {
          code,
        },
      });
      setMfaCodeCallback(null);
      setSshKey(data);
    } catch ({ msg }) {
      if (msg) toast.error(msg as string);
    }
    setIsUnlockLoading(false);
  }

  return (
    <Block>
      {sshKey ? (
        <Stack spacing={3} py={2}>
          <FormControl variant={"outlined"} fullWidth>
            <InputLabel id={"key-mode-select"}>格式</InputLabel>
            <Select
              labelId={"key-mode-select"}
              label={"格式"}
              value={keyMode}
              defaultValue={"ssh"}
              onChange={(e) => setKeyMode(e.target.value as "pem" | "ssh")}
            >
              <MenuItem value={"ssh"}>SSH</MenuItem>
              <MenuItem value={"pem"}>PEM</MenuItem>
            </Select>
          </FormControl>

          <TextField
            label={"用户名"}
            fullWidth
            value={sshKey.username}
            InputProps={{
              readOnly: true,
            }}
            onClick={(e: any) => {
              e.target.select();
            }}
          />
          <TextField
            label={"公钥"}
            fullWidth
            multiline
            value={sshKey[keyMode].public.trimEnd()}
            InputProps={{
              readOnly: true,
            }}
            onClick={(e: any) => {
              e.target.select();
            }}
          />
          <TextField
            label={"密钥"}
            fullWidth
            multiline
            value={sshKey[keyMode].private.trimEnd()}
            InputProps={{
              readOnly: true,
            }}
            onClick={(e: any) => {
              e.target.select();
            }}
          />
        </Stack>
      ) : (
        <Stack alignItems={"center"} justifyContent={"center"} height={"10rem"}>
          <Stack
            justifyContent={"center"}
            alignItems={"center"}
            flexDirection={"row"}
            mb={"1rem"}
          >
            <Link color={"disabled"} />
            <LockPerson
              sx={{
                fontSize: "3rem",
                paddingBottom: "0.55rem",
                paddingX: "0.5rem",
              }}
            />
            <Link color={"disabled"} />
          </Stack>

          <LoadingButton
            variant={"outlined"}
            disabled={mfaEnabled === false}
            loading={isUnlockLoading}
            onClick={() => setMfaCodeCallback(onShowSshKeys)}
          >
            {mfaEnabled === false ? "请在个人资料页开启 MFA" : "解锁"}
          </LoadingButton>
        </Stack>
      )}
    </Block>
  );
};
export default Ssh;
