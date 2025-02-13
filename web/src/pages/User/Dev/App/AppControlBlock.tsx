import { FC, useState } from "react";
import toast from "react-hot-toast";

import Block from "@components/user/Block";
import AppForm from "@components/user/dev/app/AppForm";
import AppTableRow from "@components/user/dev/app/AppTableRow";
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

import { useUserApiV1 } from "@api/v1/user/hook";
import { apiV1User } from "@api/v1/user/base";

import { useShallow } from "zustand/react/shallow";
import { useAppModifyForm } from "@store/useAppForm";
import useUser from "@store/useUser";

export const AppControlBlock: FC = () => {
  const apps = useUser((state) => state.apps);
  const setApps = useUser((state) => state.setApps);
  const setDialog = useUser((state) => state.setDialog);

  useUserApiV1<App.Detailed[]>("dev/app/", {
    immutable: true,
    enableLoading: true,
    onSuccess: (data) => setApps(data),
  });

  const [name, callback, permitAll, permitGroups] = useAppModifyForm(
    useShallow((state) => [
      state.name,
      state.callback,
      state.permitAll,
      state.permitGroups,
    ]),
  );
  const setAppModifyForm = useAppModifyForm((state) => state.setApp);
  const [onModifyApp, setOnModifyApp] = useState<App.Detailed | null>(null);
  const [modifyingApp, setModifyingApp] = useState(false);

  async function handleDeleteApp(app: App.Detailed) {
    try {
      const ok = await setDialog({
        title: "确认删除",
        content: `正在删除名称为 ${app.name} ，appCode 为 ${app.appCode} 的应用`,
      });
      if (!ok) return;
      await apiV1User.delete("dev/app/", {
        params: {
          id: app.id,
        },
      });
      setApps((apps || []).filter((a) => a.id !== app.id));
      toast.success("删除成功");
    } catch (err) {
      if (err instanceof Error) toast.error(err.message);
    }
  }

  async function handleModifyLinkOff(app: App.Detailed) {
    try {
      const ok = await setDialog({
        title: "确认修改对接状态",
        content: `正在将名称为 ${app.name} 的应用接入状态修改为：${
          !app.linkOff ? "未对接" : "已对接"
        }`,
      });
      if (!ok) return;
      await apiV1User.put("dev/app/linkOff", {
        id: app.id,
        linkOff: !app.linkOff,
      });
      setApps(
        (apps || []).map((a) =>
          a.id === app.id ? ({ ...a, linkOff: !a.linkOff } as App.Detailed) : a,
        ),
      );
      toast.success("修改成功");
    } catch (err) {
      if (err instanceof Error) toast.error(err.message);
    }
  }

  async function handleModifyApp() {
    if (!onModifyApp) return;
    setModifyingApp(true);
    try {
      await apiV1User.put("dev/app/", {
        id: onModifyApp.id,
        name,
        callback,
        permitAll,
        permitGroups: permitGroups?.map((g) => g.id),
      });
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
              } as App.Detailed)
            : app,
        ),
      );
      setOnModifyApp(null);
    } catch (err) {
      if (err instanceof Error) toast.error(err.message);
    }
    setModifyingApp(false);
  }

  function showModifyAppDialog(app: App.Detailed) {
    setAppModifyForm(app);
    setOnModifyApp(app);
  }

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
                <TableCell>对接</TableCell>
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
                  onModifyLinkOff={handleModifyLinkOff}
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
