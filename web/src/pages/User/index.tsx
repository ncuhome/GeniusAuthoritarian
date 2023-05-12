import { FC, ReactNode, useMemo, useState } from "react";
import { Routes, Route } from "react-router-dom";
import "./style.css";

import { Navigation, Profile } from "./pages";
import { PageNotFound } from "@components";
import { Header } from "./components";
import {
  Box,
  Stack,
  Dialog,
  DialogTitle,
  DialogContent,
  DialogActions,
  Button,
} from "@mui/material";

import { shallow } from "zustand/shallow";
import { useUser } from "@store";

type RouterElement = {
  name: string;
  path: string;
  element: ReactNode;
};

const BaseUserRouters: RouterElement[] = [
  { name: "导航", path: "", element: <Navigation /> },
  { name: "个人资料", path: "profile", element: <Profile /> },
];

const UserRoutersExtra: {
  [name: string]: RouterElement[];
} = {
  研发: [{ name: "应用管理", path: "app", element: <PageNotFound /> }],
};

export const User: FC = () => {
  const [dialog, openDialog] = useUser(
    (state) => [state.dialog, state.openDialog],
    shallow
  );
  const groups = useUser((state) => state.groups);

  const UserRouters = useMemo<RouterElement[]>(() => {
    let routers = BaseUserRouters;
    for (let i = 0; i < groups.length; i++) {
      if (UserRoutersExtra[groups[i]])
        routers = routers.concat(UserRoutersExtra[groups[i]]);
    }
    return routers;
  }, [groups]);
  const [currentTab, setCurrentTab] = useState<number>(() => {
    const matchPath = window.location.pathname.replace("/user/", "");
    for (let i = 0; i < UserRouters.length; i++) {
      if (matchPath === UserRouters[i].path) {
        return i;
      }
    }
    return 0;
  });

  return (
    <Stack id={"user"}>
      <Box
        sx={{
          width: "100%",
          position: "sticky",
          height: "3.5rem",
        }}
      >
        <Header
          routers={UserRouters}
          currentTab={currentTab}
          onChangeTab={setCurrentTab}
        />
      </Box>
      <Box
        sx={{
          overflowY: "overlay",
          minHeight: "calc(100% - 3.5rem)",
        }}
      >
        <Routes>
          {UserRouters.map((tab) => (
            <Route key={tab.path} {...tab} />
          ))}
          <Route path={"*"} element={<PageNotFound />} />
        </Routes>
      </Box>

      <Dialog
        sx={{ "& .MuiDialog-paper": { width: "60%", maxHeight: 435 } }}
        maxWidth="xs"
        open={openDialog}
        onClose={() => dialog.callback(false)}
      >
        <DialogTitle>{dialog.title}</DialogTitle>
        {dialog.content ? (
          <DialogContent>{dialog.content}</DialogContent>
        ) : null}
        <DialogActions>
          <Button autoFocus onClick={() => dialog.callback(false)}>
            取消
          </Button>
          <Button onClick={() => dialog.callback(true)}>确认</Button>
        </DialogActions>
      </Dialog>
    </Stack>
  );
};
export default User;
