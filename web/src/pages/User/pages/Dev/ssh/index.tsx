import { FC, useState } from "react";
import toast from "react-hot-toast";

import {
  Stack,
  TextField,
  Select,
  MenuItem,
  FormControl,
  InputLabel,
  ButtonGroup,
  Fade,
  Box,
} from "@mui/material";
import { LoadingButton } from "@mui/lab";
import { LockPerson, Link, Lock, RestartAlt } from "@mui/icons-material";
import Block from "@components/user/Block";
import TipButton from "@components/TipButton";
import TipIconButton from "@components/TipIconButton";

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

  const [isUnlockLoading, setIsUnlockLoading] = useState(false);

  const { isLoading: isMfaStatusLoading } = useUserApiV1("profile/mfa", {
    immutable: sshKey !== null,
    enableLoading: true,
    onSuccess: (data: any) => setMfaEnabled(data.mfa),
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

  async function onResetSshKey(code: string) {
    try {
      const {
        data: { data },
      } = await apiV1User.put("dev/ssh/", {
        code,
      });
      setMfaCodeCallback(null);
      setSshKey(data);
      toast.success("SSH 密钥已重新生成");
    } catch ({ msg }) {
      if (msg) toast.error(msg as string);
    }
  }

  return (
    <Block>
      <Stack spacing={2}>
        <Stack flexDirection={"row"}>
          <ButtonGroup variant="outlined">
            <TipButton
              title={"重置 SSH 密钥"}
              onClick={() => setMfaCodeCallback(onResetSshKey)}
            >
              <RestartAlt />
            </TipButton>
          </ButtonGroup>

          <Fade in={Boolean(sshKey)}>
            <Box
              sx={{
                marginLeft: "0.6rem",
              }}
            >
              <TipIconButton
                title={"锁定"}
                color="primary"
                onClick={() => setSshKey(null)}
              >
                <Lock />
              </TipIconButton>
            </Box>
          </Fade>
        </Stack>

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
          <Stack
            alignItems={"center"}
            justifyContent={"center"}
            height={"10rem"}
          >
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
              loading={isUnlockLoading || isMfaStatusLoading}
              onClick={() => setMfaCodeCallback(onShowSshKeys)}
            >
              {mfaEnabled === false ? "请在个人资料页开启 MFA" : "解锁"}
            </LoadingButton>
          </Stack>
        )}
      </Stack>
    </Block>
  );
};
export default Ssh;
