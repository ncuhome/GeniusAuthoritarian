import { FC, useCallback, useMemo, useState } from "react";
import { useInterval, useMount, useLoadingToast } from "@hooks";
import toast from "react-hot-toast";
import moment from "moment";

import { Block } from "@/pages/User/components";
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
} from "@mui/material";

import { GetUserProfile } from "@api/v1/user/profile";

import { useUser } from "@store";

export const Profile: FC = () => {
  const profile = useUser((state) => state.profile);
  const setProfile = useUser((state) => state.setState("profile"));

  const [onRequest, setOnRequest] = useState(true);

  const userGroups: string = useMemo(() => {
    if (!profile) return "";
    return profile.user.groups.map((group) => group.name).join("、");
  }, [profile]);

  const [loadProfileFailedToast, closeLoadProfileFailedToast] =
    useLoadingToast();

  const loadProfile = useCallback(async () => {
    setOnRequest(true);
    try {
      const data = await GetUserProfile();
      setProfile(data);
      closeLoadProfileFailedToast("Profile Loaded");
    } catch ({ msg }) {
      if (msg) loadProfileFailedToast(msg as string);
    }
    setOnRequest(false);
  }, []);

  const GridTextField: FC<TextFieldProps> = ({ ...props }) => {
    return (
      <Grid item xs={12} sm={6} position={"relative"}>
        {props.value ? (
          <TextField
            variant={"outlined"}
            inputProps={{
              readOnly: true,
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

  useInterval(loadProfile, profile || onRequest ? null : 2000);
  useMount(() => {
    if (!profile) loadProfile();
    else setOnRequest(false);
  });

  return (
    <Container>
      <Block title={"Profile"}>
        <Grid container spacing={2} marginTop={"0"}>
          <GridTextField label={"姓名"} value={profile?.user.name} />
          <GridTextField label={"电话"} value={profile?.user.phone} />
          <GridTextField label={"身份组"} value={userGroups} />
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
                    <TableCell>{record.ip}</TableCell>
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
