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
import KeyPair from "@components/user/dev/ssh/KeyPair";

import { apiV1User } from "@api/v1/user/base";

import useU2F from "@hooks/useU2F";

const Ssh: FC = () => {
  const { isLoading: u2fLoading, refreshToken } = useU2F();

  const [sshKey, setSshKey] = useState<User.SSH.Keys | null>(null);
  const [keyMode, setKeyMode] = useState<User.SSH.KeyMode>("ssh");

  const [isUnlockLoading, setIsUnlockLoading] = useState(false);

  async function onShowSshKeys() {
    setIsUnlockLoading(true);
    try {
      const token = await refreshToken();
      try {
        const {
          data: { data },
        } = await apiV1User.get("dev/ssh/", {
          params: {
            token,
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
    const token = await refreshToken("重置密钥不会断开已连接终端");
    try {
      const {
        data: { data },
      } = await apiV1User.put("dev/ssh/", {
        token,
      });
      setSshKey(data);
      toast.success("SSH 密钥已重新生成");
    } catch ({ msg }) {
      if (msg) toast.error(msg as string);
    }
  }

  async function onKillAllProcess() {
    const token = await refreshToken();
    try {
      await apiV1User.post("dev/ssh/killall", {
        token,
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
          <KeyPair
            spacing={3}
            py={2}
            keyMode={keyMode}
            keys={sshKey}
            onSetKeyMode={setKeyMode}
          />
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
              loading={isUnlockLoading}
              onClick={onShowSshKeys}
            >
              解锁
            </LoadingButton>
          </Stack>
        )}
      </Stack>
    </Block>
  );
};
export default Ssh;
