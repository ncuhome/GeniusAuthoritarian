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
  Grid,
  GridProps,
  TextField,
  TextFieldProps,
  Table,
  TableBody,
  TableCell,
  TableRow,
  TableHead,
  Skeleton,
  Checkbox,
  FormControlLabel,
  Stack,
} from "@mui/material";

import { useUserApiV1 } from "@api/v1/user/hook";

import useUser from "@store/useUser";
import useU2fDialog from "@store/useU2fDialog";

const GridItem: FC<GridProps> = ({ children, ...props }) => (
  <Grid item xs={12} sm={6} {...props}>
    {children ? children : <Skeleton variant={"rounded"} height={56} />}
  </Grid>
);

const GridTextField: FC<TextFieldProps> = ({ ...props }) => {
  return (
    <GridItem>
      {props.value ? (
        <TextField
          variant={"outlined"}
          inputProps={{
            readOnly: true,
            style: {
              cursor: "default",
            },
          }}
          fullWidth
          onClick={async () => {
            try {
              await navigator.clipboard.writeText(props.value as string);
              toast.success("已复制");
            } catch (e) {
              console.log(e);
              toast.error(`复制失败: ${e}`);
            }
          }}
          {...props}
        />
      ) : undefined}
    </GridItem>
  );
};

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

  const prefer = useU2fDialog((state) => state.prefer);
  const setPrefer = useU2fDialog((state) => state.setPrefer);

  const onChangePrefer = async (target: User.U2F.Methods) => {
    try {
      await setPrefer(target);
    } catch ({ msg }) {
      if (msg) toast.error(msg as string);
    }
  };

  return (
    <Container>
      <Block title={"Profile"}>
        <Grid container spacing={2} marginTop={0} marginBottom={3}>
          <GridTextField label={"姓名"} value={profile?.user.name} />
          <GridTextField label={"电话"} value={profile?.user.phone} />
          <GridTextField label={"身份组"} value={userGroups} />
        </Grid>
      </Block>

      <Block title={"Security"}>
        <ChildBlock
          title={"首选校验方式"}
          desc={"更改 U2F 校验时默认选择的校验方式"}
        >
          <Stack flexDirection={"row"}>
            <FormControlLabel
              label={"电话"}
              control={
                <Checkbox
                  checked={prefer === "phone"}
                  onClick={(e) => onChangePrefer("phone")}
                />
              }
            />
            <FormControlLabel
              label={"双因素认证"}
              control={
                <Checkbox
                  checked={prefer === "mfa"}
                  onClick={(e) => onChangePrefer("mfa")}
                />
              }
            />
            <FormControlLabel
              label={"通行密钥"}
              control={
                <Checkbox
                  checked={prefer === "passkey"}
                  onClick={(e) => onChangePrefer("passkey")}
                />
              }
            />
          </Stack>
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
