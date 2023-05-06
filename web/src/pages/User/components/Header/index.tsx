import { FC } from "react";
import { useNavigate } from "react-router-dom";
import logo from "@/assets/img/logo-lg.png";
import "./styles.css";

import { Box, Stack, Paper, Tabs, Tab, Typography } from "@mui/material";

interface Props {
  routers: { [name: string]: string };
  currentTab?: string;
}

export const Header: FC<Props> = ({ routers, currentTab }) => {
  const nav = useNavigate();

  const handleGoHome = () => nav("/user/");

  return (
    <Stack
      id={"user-nav"}
      flexDirection={"row"}
      sx={{
        px: "3rem",
        height: "inherit",
      }}
      component={Paper}
      elevation={6}
    >
      <Box
        sx={{
          height: "100%",
          display: "flex",
          alignItems: "center",
          marginRight: "3rem",
          "&>img": {
            height: "60%",
          },
        }}
        onClick={handleGoHome}
      >
        <img src={logo} alt={""} />
      </Box>

      <Tabs
        value={currentTab}
        textColor="inherit"
      >
        {Object.keys(routers).map((name) => (
          <Tab
            key={name}
            label={<Typography variant={"subtitle1"}>{name}</Typography>}
            value={name}
            onClick={() => nav(routers[name])}
            disableRipple
          />
        ))}
      </Tabs>
    </Stack>
  );
};
export default Header;
