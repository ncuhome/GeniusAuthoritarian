import { FC, useState } from "react";
import toast from "react-hot-toast";

import { Stack, ButtonGroup, Fade, Box } from "@mui/material";
import { LoadingButton } from "@mui/lab";
import {
  LockPerson,
  Link,
  Lock,
  RestartAlt,
  PersonOff,
  Fingerprint,
} from "@mui/icons-material";
import Block from "@components/user/Block";
import TipButton from "@components/TipButton";
import TipIconButton from "@components/TipIconButton";
import KeyPair from "@components/user/dev/ssh/KeyPair";

import { apiV1User } from "@api/v1/user/base";

import useU2F from "@hooks/useU2F";

const invalidKey: User.SSH.Keys = {
  username: "114514",
  ssh: {
    public:
      "ssh-ed25519 F1WfOB4G90IdPFsJPr3XhZSupRSSJJ2MFfPxL081t5qEdLxrOJjdq70ou27PzazCCnwX",
    private:
      "-----BEGIN OPENSSH PRIVATE KEY-----\n" +
      "L0RtR+88cLywV7W4cMVvE4p/iIcIuyRx0rN0ViIRVcykH9U5DbU/+CEML+n0hPGh\n" +
      "fdA2hhLT7+0i4uh9VDmm9MnUrEzpkUNlRgOhlwZdOMyjnhDu8CQ6ARG9+6h0QPMr\n" +
      "UIO7na+qmJvR0CuY4GQj9Xxlla/68k0YXopE/tn+SvpVuLX79PzcvfMQh7Jn0Xyo\n" +
      "0sjCZUOWkKBBWTw7B/Lzl79FD3R+EAdz6uY7m9ulrk4Tsg9MQ0BjImgzfDA/qh4N\n" +
      "0YeWAr8y2qNrVE/ddsrrER6VR5zKRekp1jZcVFd9QFbwMZU0MaBYJ+1a\n" +
      "-----END OPENSSH PRIVATE KEY-----",
  },
  pem: {
    public: "",
    private: "",
  },
};

const Ssh: FC = () => {
  const { refreshToken } = useU2F();

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
          headers: {
            Authorization: token,
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
      } = await apiV1User.put("dev/ssh/", undefined, {
        headers: {
          Authorization: token,
        },
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
      await apiV1User.post("dev/ssh/killall", undefined, {
        headers: {
          Authorization: token,
        },
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

        <Box sx={{ position: "relative" }}>
          <KeyPair
            spacing={3}
            py={2}
            keyMode={sshKey ? keyMode : "ssh"}
            keys={sshKey ?? invalidKey}
            onSetKeyMode={setKeyMode}
            sx={
              sshKey
                ? { transition: "filter 0.25s" }
                : {
                    filter: "blur(6px)",
                  }
            }
          />

          <Fade in={!sshKey} appear={false}>
            <Stack
              alignItems={"center"}
              justifyContent={"center"}
              sx={{
                position: "absolute",
                left: 0,
                top: 0,
                bottom: 0,
                right: 0,
              }}
            >
              <Stack
                justifyContent={"center"}
                alignItems={"center"}
                flexDirection={"row"}
                mb={"1rem"}
              >
                <Link fontSize={"large"} />
                <LockPerson
                  sx={{
                    fontSize: "4rem",
                    paddingBottom: "0.75rem",
                    paddingX: "1rem",
                  }}
                />
                <Link fontSize={"large"} />
              </Stack>

              <LoadingButton
                variant={"outlined"}
                loading={isUnlockLoading}
                onClick={onShowSshKeys}
                startIcon={<Fingerprint />}
              >
                解锁
              </LoadingButton>
            </Stack>
          </Fade>
        </Box>
      </Stack>
    </Block>
  );
};
export default Ssh;
