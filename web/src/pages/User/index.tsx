import { FC, ReactNode, useMemo, useState, lazy } from "react";
import { Routes, Route } from "react-router-dom";
import useKeyDown from "@hooks/useKeyDown";
import "./style.css";

const Dev = lazy(() => import("./pages/Dev"));
import Navigation from "./pages/Navigation";
import Profile from "./pages/Profile";

import U2fDialog from "@components/user/U2fDialog";
import Suspense from "@components/Suspense";
import PageNotFound from "@components/PageNotFound";
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
} from "@mui/material";

import { useShallow } from "zustand/react/shallow";
import useUser from "@store/useUser";
import useTheme from "@store/useTheme";

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
  研发: [
    {
      name: "研发中控",
      path: "dev",
      element: (
        <Suspense>
          <Dev />
        </Suspense>
      ),
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
    <>
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
        sx={{
          backgroundColor: isDarkTheme ? undefined : "#fff",
          colorScheme: isDarkTheme ? undefined : "light",
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
          <Routes>
            {UserRouters.map((tab) => (
              <Route key={tab.path} {...tab} />
            ))}
            <Route path={"*"} element={<PageNotFound />} />
          </Routes>
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
    </>
  );
};
export default User;
