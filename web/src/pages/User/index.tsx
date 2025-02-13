import { FC, useMemo, useState } from "react";
import { Outlet } from "react-router";
import useKeyDown from "@hooks/useKeyDown";
import "./style.css";

import { createTheme, ThemeProvider } from "@mui/material/styles";
import U2fDialog from "@components/user/U2fDialog";
import NavHeader from "@components/user/NavHeader";
import { Toaster } from "react-hot-toast";
import {
  Box,
  Stack,
  Dialog,
  DialogTitle,
  DialogContent,
  DialogActions,
  Button,
  Paper,
} from "@mui/material";

import { useShallow } from "zustand/react/shallow";
import useUser from "@store/useUser";
import useTheme from "@store/useTheme";

type RouterElement = {
  name: string;
  path: string;
};

const BaseUserRouters: RouterElement[] = [
  { name: "导航", path: "" },
  { name: "个人资料", path: "profile" },
];

const UserRoutersExtra: {
  [name: string]: RouterElement[];
} = {
  研发: [
    {
      name: "研发中控",
      path: "dev",
    },
  ],
  中心: [
    {
      name: "管理后台",
      path: "admin",
    },
  ],
};

export const User: FC = () => {
  const [dialog, openDialog, dialogResolver] = useUser(
    useShallow((state) => [
      state.dialog,
      state.openDialog,
      state.dialogResolver,
    ]),
  );
  useKeyDown("Enter", () => {
    if (openDialog) {
      dialogResolver?.(true);
    }
  });

  const groups = useUser((state) => state.groups);

  const isDarkTheme = useTheme((state) => state.dark);

  const theme = useMemo(
    () =>
      createTheme({
        palette: {
          mode: isDarkTheme ? "dark" : undefined,
        },
      }),
    [isDarkTheme],
  );

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
    <ThemeProvider theme={theme}>
      <Toaster
        toastOptions={
          isDarkTheme
            ? {
                style: {
                  borderRadius: "20px",
                  background: "#2f2f2f",
                  color: "#fff",
                },
              }
            : {
                style: {
                  borderRadius: "20px",
                },
              }
        }
      />
      <Stack
        id={"user"}
        component={Paper}
        elevation={0}
        sx={{
          color: "text.primary",
        }}
      >
        <Box
          sx={{
            width: "100%",
            position: "sticky",
            height: "3.5rem",
          }}
        >
          <NavHeader
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
          <Outlet />
        </Box>

        <U2fDialog />

        <Dialog
          sx={{ "& .MuiDialog-paper": { width: "60%", maxHeight: 435 } }}
          maxWidth="xs"
          open={openDialog}
          onClose={() => dialogResolver?.(false)}
        >
          <DialogTitle>{dialog.title}</DialogTitle>
          {dialog.content ? (
            <DialogContent>{dialog.content}</DialogContent>
          ) : null}
          <DialogActions>
            <Button onClick={() => dialogResolver?.(false)}>取消</Button>
            <Button onClick={() => dialogResolver?.(true)}>确认</Button>
          </DialogActions>
        </Dialog>
      </Stack>
    </ThemeProvider>
  );
};
export default User;
