import { FC, PropsWithChildren } from "react";
import { useInterval } from "@hooks";
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
} from "@mui/material";

import { GetUserProfile } from "@api/v1/user/profile";

import { shallow } from "zustand/shallow";
import { useUser } from "@store";

export const Profile: FC = () => {
  const [profile] = useUser((state) => [state.profile], shallow);
  const [setProfile] = useUser((state) => [state.setState("profile")], shallow);

  async function loadProfile() {
    try {
      const data = await GetUserProfile();
      setProfile(data);
    } catch ({ msg }) {
      if (msg) toast.error(msg as string);
    }
  }

  const GridTextField: FC<PropsWithChildren & TextFieldProps> = ({
    children,
    ...props
  }) => {
    return (
      <Grid item xs={12} sm={6} /*md={4}*/>
        <TextField
          variant={"outlined"}
          inputProps={{
            readOnly: true,
          }}
          fullWidth
          {...props}
        >
          {children}
        </TextField>
      </Grid>
    );
  };

  useInterval(loadProfile, profile ? null : 2000);

  return (
    <Container>
      <Box component={Paper} elevation={5}>
        <Typography variant={"h5"} fontWeight={"bold"} marginBottom={"1rem"}>
          Profile
        </Typography>
        <Grid container spacing={2}>
          <GridTextField label={"姓名"} value={profile?.user.name} />
          <GridTextField label={"电话"} value={profile?.user.phone} />
        </Grid>
      </Box>

      {profile && profile.loginRecord.length ? (
        <Box component={Paper} elevation={5}>
          <Typography variant={"h5"} fontWeight={"bold"} marginBottom={"1rem"}>
            Record
          </Typography>
          <Table>
            <TableHead>
              <TableRow>
                <TableCell>登录时间</TableCell>
                <TableCell>目标</TableCell>
                <TableCell>IP</TableCell>
              </TableRow>
            </TableHead>
            <TableBody>
              {profile.loginRecord.reverse().map((record) => (
                <TableRow key={record.id}>
                  <TableCell>
                    {moment(record.createdAt * 1000).format("YYYY/MM/DD hh:mm")}
                  </TableCell>
                  <TableCell>{record.target}</TableCell>
                  <TableCell>{record.ip}</TableCell>
                </TableRow>
              ))}
            </TableBody>
          </Table>
        </Box>
      ) : undefined}
    </Container>
  );
};
export default Profile;
