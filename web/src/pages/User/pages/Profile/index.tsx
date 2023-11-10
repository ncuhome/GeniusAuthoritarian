import { FC, useMemo } from "react";
import { unix } from "dayjs";
import toast from "react-hot-toast";

import Block from "@components/user/Block";
import ChildBlock from "@components/user/ChildBlock";
import Ip from "@components/user/profile/Ip";
import Mfa from "@components/user/profile/Mfa";
import Passkey from "@components/user/profile/Passkey";
import {
  Container,
  Box,
  Table,
  TableBody,
  TableCell,
  TableRow,
  TableHead,
  Skeleton,
  Checkbox,
  FormControlLabel,
  Stack,
  Tooltip,
  Avatar,
  Typography,
  Paper,
} from "@mui/material";
import { PermIdentity } from "@mui/icons-material";

import { useUserApiV1 } from "@api/v1/user/hook";

import useUser from "@store/useUser";
import useU2fDialog from "@store/useU2fDialog";

const u2fMethods: {
  label: string;
  value: User.U2F.Methods;
}[] = [
  { label: "短信", value: "phone" },
  { label: "双因素认证", value: "mfa" },
  { label: "通行密钥", value: "passkey" },
];

export const Profile: FC = () => {
  const profile = useUser((state) => state.profile);
  const setProfile = useUser((state) => state.setState("profile"));

  const userGroups: string = useMemo(() => {
    if (!profile) return "";
    return profile.user.groups.map((group) => group.name).join("、");
  }, [profile]);

  useUserApiV1<User.Profile>("profile/", {
    enableLoading: true,
    onSuccess: (data) => setProfile(data),
  });

  const { data: u2fStatus, mutate: mutateU2f } = useUserApiV1<User.U2F.Status>(
    "u2f/",
    {
      onError: (err) => {
        toast.error(`载入 U2F 状态失败: ${err}`);
      },
    }
  );
  const setPrefer = useU2fDialog((state) => state.setPrefer);

  const onChangePrefer = async (target: User.U2F.Methods) => {
    try {
      await setPrefer(target);
      mutateU2f({
        ...u2fStatus!,
        prefer: target,
      });
    } catch ({ msg }) {
      if (msg) toast.error(msg as string);
    }
  };

  return (
    <Container>
      <Stack flexDirection={"row"} mt={"3rem"} mb={"2.3rem"}>
        {profile ? (
          <Avatar
            component={Paper}
            elevation={6}
            src={profile.user.avatar_url}
            sx={{
              height: 90,
              width: 90,
            }}
          />
        ) : (
          <Skeleton variant={"circular"} height={90} width={103} />
        )}
        <Stack ml={3} width={"100%"} justifyContent={"space-between"}>
          {profile ? (
            <Typography
              variant={"h5"}
              sx={{
                fontWeight: 600,
              }}
            >
              {profile.user.name}
            </Typography>
          ) : (
            <Skeleton width={75} />
          )}
          <Stack>
            <Typography variant={"body2"} color={"text.secondary"}>
              当前身份状态
            </Typography>
            <Stack flexDirection={"row"} mt={0.5}>
              <PermIdentity sx={{ marginRight: 1 }} />
              {profile ? (
                <Typography>{`${userGroups}组成员`}</Typography>
              ) : (
                <Skeleton width={120} />
              )}
            </Stack>
          </Stack>
        </Stack>
      </Stack>

      <Block title={"Security"}>
        <ChildBlock
          title={"首选校验方式"}
          desc={"更改 U2F 校验时默认选择的校验方式"}
        >
          {u2fStatus ? (
            <Stack flexDirection={"row"}>
              {u2fMethods.map((m) => (
                <Tooltip
                  key={m.value}
                  title={!(u2fStatus as any)[m.value] ? "未启用" : undefined}
                  placement={"top"}
                  arrow
                >
                  <FormControlLabel
                    label={m.label}
                    control={
                      <Checkbox
                        checked={u2fStatus?.prefer === m.value}
                        onChange={(e) =>
                          onChangePrefer(e.target.checked ? m.value : "")
                        }
                        color={
                          !(u2fStatus as any)[m.value] ? "warning" : undefined
                        }
                      />
                    }
                  />
                </Tooltip>
              ))}
            </Stack>
          ) : (
            <Skeleton height={42} width={300} />
          )}
        </ChildBlock>

        <ChildBlock
          title={"双因素认证"}
          desc={
            "两步验证在第三方登录时增加一道额外的身份认证，可以预防飞书、钉钉账号被盗用的情况。启用此功能需要使用 Google Authenticator APP 或密码保险库如 1password 等工具保存密钥与生成一次性密码"
          }
        >
          {profile ? (
            <Mfa
              enabled={profile.user.mfa}
              setEnabled={(enabled) =>
                setProfile({
                  user: {
                    ...profile.user,
                    mfa: enabled,
                  },
                  loginRecord: profile.loginRecord,
                })
              }
            />
          ) : (
            <Skeleton
              variant="rounded"
              height={35}
              sx={{
                maxWidth: "13rem",
              }}
            />
          )}
        </ChildBlock>

        <ChildBlock
          title={"通行密钥"}
          desc={
            "通行密钥可以是支持生物识别的手机电脑，可以是硬件密钥，也可以存入密码保险库跨设备同步。使用通行密钥可以免账户密码进行身份验证且自带双因素，是一种安全便捷的认证方式"
          }
        >
          <Passkey />
        </ChildBlock>
      </Block>

      {profile && profile.loginRecord.length ? (
        <Block title={"Record"} subtitle={"最近十次登录记录"}>
          <Box
            sx={{
              marginTop: "0.5rem",
              width: "100%",
              overflowY: "auto",
              whiteSpace: "nowrap",
            }}
          >
            <Table>
              <TableHead>
                <TableRow>
                  <TableCell>登录时间</TableCell>
                  <TableCell>站点</TableCell>
                  <TableCell>IP</TableCell>
                </TableRow>
              </TableHead>
              <TableBody>
                {profile.loginRecord.map((record) => (
                  <TableRow key={record.id}>
                    <TableCell>
                      {unix(record.createdAt).format("YYYY/MM/DD HH:mm")}
                    </TableCell>
                    <TableCell>{record.target}</TableCell>
                    <TableCell>
                      <Ip ip={record.ip} />
                    </TableCell>
                  </TableRow>
                ))}
              </TableBody>
            </Table>
          </Box>
        </Block>
      ) : undefined}
    </Container>
  );
};
export default Profile;
