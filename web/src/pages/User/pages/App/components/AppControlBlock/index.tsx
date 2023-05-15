import { FC, useState } from "react";
import { useLoadingToast, useMount, useInterval } from "@hooks";

import { Block } from "@/pages/User/components";
import { TipIconButton } from "@components";
import {
  Paper,
  TableContainer,
  Table,
  TableHead,
  TableRow,
  TableCell,
  TableBody,
  Stack,
} from "@mui/material";

import { DeleteForeverOutlined } from "@mui/icons-material";

import { GetOwnedAppList } from "@api/v1/user/app";

import { useUser } from "@store";

export const AppControlBlock: FC = () => {
  const apps = useUser((state) => state.apps);
  const setApps = useUser((state) => state.setState("apps"));

  const [onRequestApps, setOnRequestApps] = useState(true);
  const [loadAppsToast, closeAppsToast] = useLoadingToast();

  async function loadApps() {
    setOnRequestApps(true);
    try {
      const data = await GetOwnedAppList();
      setApps(data);
      closeAppsToast();
    } catch ({ msg }) {
      if (msg) loadAppsToast(msg as string);
    }
    setOnRequestApps(false);
  }

  useInterval(loadApps, !apps && !onRequestApps ? 2000 : null);
  useMount(() => {
    if (!apps) loadApps();
    else setOnRequestApps(false);
  });

  return (
    <Block title={"App"}>
      <Paper sx={{ width: "100%", overflowX: "auto", marginTop: "1.3rem" }}>
        <TableContainer sx={{ maxHeight: 440 }}>
          <Table stickyHeader aria-label="sticky table">
            <TableHead>
              <TableRow>
                <TableCell>名称</TableCell>
                <TableCell>AppCode</TableCell>
                <TableCell>授权</TableCell>
                <TableCell></TableCell>
              </TableRow>
            </TableHead>
            <TableBody>
              {(apps || []).map((app) => (
                <TableRow hover role="checkbox" tabIndex={-1} key={app.id}>
                  <TableCell>{app.name}</TableCell>
                  <TableCell>{app.appCode}</TableCell>
                  <TableCell>
                    {app.permitAllGroup
                      ? "ALL"
                      : app.groups.length > 0
                      ? app.groups.map((group) => group.name).join(",")
                      : "NONE"}
                  </TableCell>
                  <TableCell>
                    <Stack flexDirection={"row"}>
                      <TipIconButton title={"删除"}>
                        <DeleteForeverOutlined />
                      </TipIconButton>
                    </Stack>
                  </TableCell>
                </TableRow>
              ))}
            </TableBody>
          </Table>
        </TableContainer>
      </Paper>
    </Block>
  );
};
export default AppControlBlock;
