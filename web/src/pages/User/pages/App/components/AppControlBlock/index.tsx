import { FC, useState } from "react";
import { useLoadingToast, useMount, useInterval } from "@hooks";
import toast from "react-hot-toast";

import { AppForm } from "@/pages/User/pages/App/components";
import { Block } from "@/pages/User/components";
import { AppTableRow } from "./components";
import {
  Paper,
  TableContainer,
  Table,
  TableHead,
  TableRow,
  TableCell,
  TableBody,
  Stack,
  CircularProgress,
  Typography,
  Dialog,
  DialogTitle,
  DialogContent,
} from "@mui/material";

import { GetOwnedAppList, DeleteApp, ModifyApp, App } from "@api/v1/user/app";

import { shallow } from "zustand/shallow";
import { useUser, useAppModifyForm } from "@store";

export const AppControlBlock: FC = () => {
  const apps = useUser((state) => state.apps);
  const setApps = useUser((state) => state.setState("apps"));
  const setDialog = useUser((state) => state.setDialog);

  const [onRequestApps, setOnRequestApps] = useState(true);
  const [loadAppsToast, closeAppsToast] = useLoadingToast();

  const [name, callback, permitAll, permitGroups] = useAppModifyForm(
    (state) => [
      state.name,
      state.callback,
      state.permitAll,
      state.permitGroups,
    ],
    shallow
  );
  const [setName, setCallback, setPermitAll, setPermitGroups] =
    useAppModifyForm(
      (state) => [
        state.setState("name"),
        state.setState("callback"),
        state.setState("permitAll"),
        state.setState("permitGroups"),
      ],
      shallow
    );
  const [onModifyApp, setOnModifyApp] = useState<App | null>(null);
  const [modifyingApp, setModifyingApp] = useState(false);

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

  async function handleDeleteApp(app: App) {
    try {
      const ok = await setDialog({
        title: "确认删除",
        content: `正在删除名称为 ${app.name} ，appCode 为 ${app.appCode} 的应用`,
      });
      if (!ok) return;
      await DeleteApp(app.id);
      setApps((apps || []).filter((a) => a.id !== app.id));
      toast.success("删除成功");
    } catch ({ msg }) {
      if (msg) toast.error(msg as string);
    }
  }

  async function handleModifyApp() {
    if (!onModifyApp) return;
    setModifyingApp(true);
    try {
      await ModifyApp(
        onModifyApp.id,
        name,
        callback,
        permitAll,
        permitGroups?.map((g) => g.id)
      );
      toast.success("修改成功");
      setApps(
        (apps || []).map((app) =>
          app.id === onModifyApp.id
            ? ({
                id: onModifyApp.id,
                name: name,
                appCode: onModifyApp.appCode,
                callback: callback,
                permitAllGroup: permitAll,
                groups: permitGroups || [],
              } as App)
            : app
        )
      );
      setOnModifyApp(null);
    } catch ({ msg }) {
      if (msg) toast.error(msg as string);
    }
    setModifyingApp(false);
  }

  function showModifyAppDialog(app: App) {
    setName(app.name);
    setCallback(app.callback);
    setPermitAll(app.permitAllGroup);
    setPermitGroups(app.groups);
    setOnModifyApp(app);
  }

  useInterval(loadApps, !apps && !onRequestApps ? 2000 : null);
  useMount(() => {
    if (!apps) loadApps();
    else setOnRequestApps(false);
  });

  return (
    <Block title={"App"}>
      <Paper
        sx={{
          width: "100%",
          overflowX: "auto",
          marginTop: "1.3rem",
        }}
      >
        <TableContainer
          sx={{ height: 440, display: "flex", flexDirection: "column" }}
        >
          <Table stickyHeader aria-label="sticky table">
            <TableHead>
              <TableRow>
                <TableCell>名称</TableCell>
                <TableCell>AppCode</TableCell>
                <TableCell>授权</TableCell>
                <TableCell>回调</TableCell>
                <TableCell></TableCell>
              </TableRow>
            </TableHead>
            <TableBody>
              {(apps || []).map((app) => (
                <AppTableRow
                  key={app.id}
                  app={app}
                  onModify={showModifyAppDialog}
                  onDelete={handleDeleteApp}
                />
              ))}
            </TableBody>
          </Table>
          <Stack
            sx={{
              flexGrow: 1,
              width: "100%",
              display: apps && apps.length > 0 ? "none" : null,
            }}
            justifyContent={"center"}
            alignItems={"center"}
          >
            {apps ? (
              <Typography
                variant={"h5"}
                fontWeight={"bold"}
                sx={{ opacity: 0.5 }}
              >
                NO DATA
              </Typography>
            ) : (
              <CircularProgress />
            )}
          </Stack>
        </TableContainer>
      </Paper>

      <Dialog open={Boolean(onModifyApp)} onClose={() => setOnModifyApp(null)}>
        <DialogTitle>修改 App</DialogTitle>
        <DialogContent>
          <AppForm
            useForm={useAppModifyForm}
            loading={modifyingApp}
            submitText={"确认"}
            onSubmit={handleModifyApp}
            cancelText={"取消"}
            onCancel={() => setOnModifyApp(null)}
          />
        </DialogContent>
      </Dialog>
    </Block>
  );
};
export default AppControlBlock;
