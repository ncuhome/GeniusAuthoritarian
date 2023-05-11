import { FC } from "react";
import { useNavigate } from "react-router-dom";
import logo from "@/assets/img/logo-lg.png";
import "./styles.css";

import {
  Box,
  Stack,
  Paper,
  Tabs,
  Tab,
  Typography,
  IconButton,
} from "@mui/material";
import { LogoutRounded } from "@mui/icons-material";

import { Logout } from "@api/v1/user/base";

interface Props {
  routers: Array<{
    name: string;
    path: string;
  }>;
  currentTab: number;
  onChangeTab: (index: number) => void;
}

export const Header: FC<Props> = ({
  routers,
  currentTab,
  onChangeTab: handleChangeTab,
}) => {
  const nav = useNavigate();

  return (
    <Stack
      id={"user-nav"}
      flexDirection={"row"}
      sx={{
        px: "min(5%, 3rem)",
        height: "inherit",
      }}
      component={Paper}
      elevation={6}
      square
    >
      <Box
        sx={{
          height: "100%",
          display: { xs: "none", sm: "flex" },
          alignItems: "center",
          marginRight: "1rem",
          "&>img": {
            height: "60%",
          },
        }}
      >
        <img src={logo} alt={"NCUHOME"} />
      </Box>

      <Stack
        flexDirection={"row"}
        flexGrow={1}
        justifyContent={"space-between"}
      >
        <Tabs value={currentTab} textColor="inherit">
          {routers.map((r, index) => (
            <Tab
              key={r.name}
              label={<Typography variant={"subtitle1"}>{r.name}</Typography>}
              value={index}
              onClick={() => {
                handleChangeTab(index);
                nav(r.path);
              }}
              disableRipple
            />
          ))}
        </Tabs>

        <Stack flexDirection={"row"} alignItems={"center"}>
          <IconButton onClick={Logout}>
            <LogoutRounded />
          </IconButton>
        </Stack>
      </Stack>
    </Stack>
  );
};
export default Header;
