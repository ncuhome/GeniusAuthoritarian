import { FC, useMemo, PropsWithChildren } from "react";
import toast from "react-hot-toast";
import moment from "moment";

import { Block } from "@/pages/User/components";
import { Ip, Mfa } from "./components";
import {
  Container,
  Box,
  Grid,
  TextField,
  TextFieldProps,
  Table,
  TableBody,
  TableCell,
  TableRow,
  TableHead,
  Skeleton,
  Stack,
  Typography,
} from "@mui/material";

import { useUserApiV1 } from "@api/v1/user/hook";

import { useUser } from "@store";

const GridItem: FC<PropsWithChildren> = ({ children }) => (
  <Grid item xs={12} sm={6} position={"relative"}>
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

  return (
    <Container>
      <Block title={"Profile"}>
        <Grid container spacing={2} marginTop={0} marginBottom={3}>
          <GridTextField label={"姓名"} value={profile?.user.name} />
          <GridTextField label={"电话"} value={profile?.user.phone} />
          <GridTextField label={"身份组"} value={userGroups} />
          <GridItem>
            {profile ? (
              <Stack
                alignItems={"center"}
                height={"100%"}
                flexDirection={"row"}
              >
                <Typography>MFA:</Typography>
                <Mfa
                  ml={1.8}
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
              </Stack>
            ) : undefined}
          </GridItem>
        </Grid>
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
                      {moment(record.createdAt * 1000).format(
                        "YYYY/MM/DD HH:mm"
                      )}
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
