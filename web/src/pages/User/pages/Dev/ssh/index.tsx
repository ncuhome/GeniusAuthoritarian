import { FC, useMemo, useState } from "react";
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
import {
  LockPerson,
  Link,
  Lock,
  RestartAlt,
  PersonOff,
} from "@mui/icons-material";
import Block from "@components/user/Block";
import TipButton from "@components/TipButton";
import TipIconButton from "@components/TipIconButton";

import { useUserApiV1 } from "@api/v1/user/hook";
import { apiV1User } from "@api/v1/user/base";

import useMfaCode from "@hooks/useMfaCode";

const Ssh: FC = () => {
  const onMfaCode = useMfaCode();

  const [sshKey, setSshKey] = useState<User.SSH.Keys | null>(null);
  const [keyMode, setKeyMode] = useState<"pem" | "ssh">("ssh");

  const [isUnlockLoading, setIsUnlockLoading] = useState(false);

  const { isLoading: isMfaStatusLoading, data: mfaData } =
    useUserApiV1<User.Mfa.Status>("profile/mfa", {
      immutable: sshKey !== null,
      enableLoading: true,
    });
  const mfaEnabled = useMemo(() => mfaData?.mfa, [mfaData]);

  async function onShowSshKeys() {
    setIsUnlockLoading(true);
    try {
      const code = await onMfaCode();
      try {
        const {
          data: { data },
        } = await apiV1User.get("dev/ssh/", {
          params: {
            code,
          },
        });
        setSshKey(data);
      } catch ({ msg }) {
        if (msg) toast.error(msg as string);
      }
    } catch (err) {}
    setIsUnlockLoading(false);
  }

  async function onResetSshKey() {
    const code = await onMfaCode("重置密钥不会断开已连接终端");
    try {
      const {
        data: { data },
      } = await apiV1User.put("dev/ssh/", {
        code,
      });
      setSshKey(data);
      toast.success("SSH 密钥已重新生成");
    } catch ({ msg }) {
      if (msg) toast.error(msg as string);
    }
  }

  async function onKillAllProcess() {
    const code = await onMfaCode();
    try {
      await apiV1User.post("dev/ssh/killall", {
        code,
      });
      toast.success("已发送 KILLALL 指令");
    } catch ({ msg }) {
      if (msg) toast.error(msg as string);
    }
  }

  return (
    <Block>
      <Stack spacing={2}>
        <Stack flexDirection={"row"}>
          <ButtonGroup variant="outlined">
            <TipButton title={"重置 SSH 密钥"} onClick={onResetSshKey}>
              <RestartAlt />
            </TipButton>
            <TipButton title={"结束进程 (终端)"} onClick={onKillAllProcess}>
              <PersonOff />
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
            height={"15rem"}
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
              onClick={onShowSshKeys}
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
