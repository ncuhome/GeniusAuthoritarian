import { FC, ReactNode, useState } from "react";
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

const UserRouters: Array<{
  name: string;
  path: string;
  element: ReactNode;
}> = [
  { name: "导航", path: "", element: <Navigation /> },
  { name: "个人资料", path: "profile", element: <Profile /> },
];

export const User: FC = () => {
  const [currentTab, setCurrentTab] = useState<number>(() => {
    const matchPath = window.location.pathname.replace("/user/", "");
    for (let i = 0; i < UserRouters.length; i++) {
      if (matchPath === UserRouters[i].path) {
        return i;
      }
    }
    return 0;
  });
  const [dialog, openDialog] = useUser(
    (state) => [state.dialog, state.openDialog],
    shallow
  );

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
