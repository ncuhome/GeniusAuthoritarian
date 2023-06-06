import { FC, useMemo } from "react";
import toast from "react-hot-toast";
import moment from "moment";

import { Block, BlockTitle } from "@/pages/User/components";
import { Ip } from "./components";
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
  Chip,
  ButtonGroup,
  Button,
} from "@mui/material";
import { LoadingButton } from "@mui/lab";
import { Done, Remove } from "@mui/icons-material";

import { useUserApiV1 } from "@api/v1/user/hook";

import { useUser } from "@store";

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

  const GridTextField: FC<TextFieldProps> = ({ ...props }) => {
    return (
      <Grid item xs={12} sm={6} position={"relative"}>
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
        ) : (
          <Skeleton variant={"rounded"} height={56} />
        )}
      </Grid>
    );
  };

  return (
    <Container>
      <Block title={"Profile"}>
        <Grid container spacing={2} marginTop={0} marginBottom={3}>
          <GridTextField label={"姓名"} value={profile?.user.name} />
          <GridTextField label={"电话"} value={profile?.user.phone} />
          <GridTextField label={"身份组"} value={userGroups} />
        </Grid>

        <BlockTitle>MFA</BlockTitle>

        <Stack flexDirection={"row"} marginTop={1}>
          <Chip
            label={profile?.user.mfa ? "已开启" : "未启用"}
            variant={"outlined"}
            icon={
              profile?.user.mfa ? (
                <Done color={"success"} fontSize="small" />
              ) : (
                <Remove />
              )
            }
          />
        </Stack>
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
