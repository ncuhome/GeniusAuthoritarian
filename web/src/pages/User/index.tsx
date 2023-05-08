import { FC, ReactNode, useState } from "react";
import { Routes, Route } from "react-router-dom";
import "./style.css";

import { Navigation, Profile } from "./pages";
import { PageNotFound } from "@components";
import { Header } from "./components";
import { Box, Stack } from "@mui/material";

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

  return (
    <Stack
      id={"user"}
    >
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
          overflowY: 'overlay',
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
    </Stack>
  );
};
export default User;
