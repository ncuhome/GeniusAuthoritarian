import { FC, useCallback, useMemo, useState } from "react";
import { useInterval, useMount } from "@hooks";
import toast from "react-hot-toast";
import moment from "moment";

import {
  Container,
  Box,
  Paper,
  Grid,
  TextField,
  TextFieldProps,
  Typography,
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

  const loadProfile = useCallback(async () => {
    setOnRequest(true);
    try {
      const data = await GetUserProfile();
      setProfile(data);
    } catch ({ msg }) {
      if (msg) toast.error(msg as string);
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
      <Box component={Paper} elevation={5}>
        <Typography variant={"h5"} fontWeight={"bold"} marginBottom={"1rem"}>
          Profile
        </Typography>
        <Grid container spacing={2}>
          <GridTextField label={"姓名"} value={profile?.user.name} />
          <GridTextField label={"电话"} value={profile?.user.phone} />
          <GridTextField label={"身份组"} value={userGroups} />
        </Grid>
      </Box>

      {profile && profile.loginRecord.length ? (
        <Box component={Paper} elevation={5}>
          <Typography variant={"h5"} fontWeight={"bold"}>
            Record
          </Typography>
          <Typography variant={"subtitle2"} marginBottom={"1rem"}>
            最近十次登录记录
          </Typography>
          <Box
            sx={{
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
        </Box>
      ) : undefined}
    </Container>
  );
};
export default Profile;
