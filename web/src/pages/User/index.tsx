import { FC, useMemo } from "react";
import { Routes, useMatch } from "react-router-dom";

import { Header } from "./components";
import { Box } from "@mui/material";

const routerTabs = {
  导航: "",
  个人资料: "profile",
} as { [key: string]: string };

export const User: FC = () => {
  const match = useMatch("/user/:path");
  const currentTab = useMemo(() => {
    const keys = Object.keys(routerTabs);
    if (!match) return keys[0];
    for (let i = 0; i < keys.length; i++) {
      if (match.params.path === routerTabs[keys[i]]) {
        return keys[i];
      }
    }
    return keys[0];
  }, [match]);

  return (
    <Box>
      <Box
        sx={{
          width: "100%",
          position: "sticky",
          height: "3.5rem",
        }}
      >
        <Header routers={routerTabs} currentTab={currentTab} />
      </Box>
    </Box>
  );
};
export default User;
