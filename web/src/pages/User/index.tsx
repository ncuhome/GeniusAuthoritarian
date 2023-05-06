import { FC, ReactNode, useMemo, useState } from "react";
import { Routes, Route } from "react-router-dom";

import { PageNotFound } from "@components";
import { Header } from "./components";
import { Box, Stack } from "@mui/material";

const routerTabs = [
  { name: "导航", path: "", element: undefined },
  { name: "个人资料", path: "profile", element: undefined },
] as Array<{
  name: string;
  path: string;
  element: ReactNode;
}>;

export const User: FC = () => {
  const [currentTab, setCurrentTab] = useState<number>(() => {
    const matchPath = window.location.pathname.replace("/user/", "");
    for (let i = 0; i < routerTabs.length; i++) {
      if (matchPath === routerTabs[i].path) {
        return i;
      }
    }
    return 0;
  });
  const routeElements = useMemo(
    () =>
      routerTabs.map((tab) => (
        <Route key={tab.path} path={tab.path} element={tab.element} />
      )),
    []
  );

  return (
    <Stack
      sx={{
        width: "100%",
        height: "100%",
      }}
    >
      <Box
        sx={{
          width: "100%",
          position: "sticky",
          height: "3.5rem",
        }}
      >
        <Header
          routers={routerTabs}
          currentTab={currentTab}
          onChangeTab={setCurrentTab}
        />
      </Box>
      <Box
        sx={{
          minHeight: "calc(100% - 3.5rem)",
        }}
      >
        <Routes>
          {routeElements}
          <Route path={"*"} element={<PageNotFound />} />
        </Routes>
      </Box>
    </Stack>
  );
};
export default User;
