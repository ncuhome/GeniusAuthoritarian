import { FC } from "react";
import { useNavigate } from "react-router-dom";
import "./NavHeader.css";

import logo_white from "@/assets/img/logo-white.png";
import logo_white_webp from "@/assets/img/logo-white.webp";
import logo_dark from "@/assets/img/logo-dark.png";
import logo_dark_webp from "@/assets/img/logo-dark.webp";

import {
  Box,
  Stack,
  Paper,
  Tabs,
  Tab,
  Typography,
  IconButton,
} from "@mui/material";
import { LogoutRounded, DarkMode, LightMode } from "@mui/icons-material";

import { Logout } from "@api/v1/user/base";

import useUser from "@store/useUser";
import useTheme from "@store/useTheme";

interface Props {
  routers: Array<{
    name: string;
    path: string;
  }>;
  currentTab: number;
  onChangeTab: (index: number) => void;
}

export const NavHeader: FC<Props> = ({
  routers,
  currentTab,
  onChangeTab: handleChangeTab,
}) => {
  const nav = useNavigate();

  const setDialog = useUser((state) => state.setDialog);

  const darkTheme = useTheme((state) => state.dark);
  const setDarkTheme = useTheme((state) => state.setState("dark"));

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
          "&>picture": {
            height: "60%",
            "&>img": {
              height: "100%",
            },
          },
        }}
      >
        <picture>
          <source
            type="image/webp"
            srcSet={darkTheme ? logo_white_webp : logo_dark_webp}
          />
          <img src={darkTheme ? logo_white : logo_dark} alt={"NCUHOME"} />
        </picture>
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

        <Stack
          flexDirection={"row"}
          alignItems={"center"}
          sx={{
            "&>*": {
              marginLeft: "0.5rem",
            },
          }}
        >
          <IconButton onClick={() => setDarkTheme(!darkTheme)}>
            {darkTheme ? <DarkMode /> : <LightMode />}
          </IconButton>
          <IconButton
            onClick={async () => {
              const ok = await setDialog({
                title: "注销登录",
              });
              if (ok) Logout();
            }}
          >
            <LogoutRounded />
          </IconButton>
        </Stack>
      </Stack>
    </Stack>
  );
};
export default NavHeader;
