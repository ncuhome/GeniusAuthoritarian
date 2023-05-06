import { FC, ReactNode, useMemo } from "react";
import { Routes, Route, useMatch, PathRouteProps } from "react-router-dom";

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
  const match = useMatch("/user/:path");
  const currentTab = useMemo(() => {
    if (!match) return 0;
    for (let i = 0; i < routerTabs.length; i++) {
      if (match.params.path === routerTabs[i].path) {
        return i;
      }
    }
    return 0;
  }, [match]);
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
        <Header routers={routerTabs} currentTab={currentTab} />
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
