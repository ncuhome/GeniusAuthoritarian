import { FC } from "react";
import { useNavigate } from "react-router-dom";
import toast from "react-hot-toast";
import "./NavHeader.css";

import DarkModeSwitch from "@components/DarkModeSwitch";
import Picture from "@components/Picture";
import {
  Box,
  Stack,
  Paper,
  Tabs,
  Tab,
  Typography,
  IconButton,
  Divider,
  useTheme as useMuiTheme,
} from "@mui/material";
import { LogoutRounded } from "@mui/icons-material";

import { Logout, apiV1User } from "@api/v1/user/base";

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
  const muiTheme = useMuiTheme();

  const setDialog = useUser((state) => state.setDialog);

  const darkTheme = useTheme((state) => state.dark);
  const setDarkTheme = useTheme((state) => state.setState("dark"));

  const onLogout = async () => {
    try {
      await apiV1User.post("logout");
      Logout();
    } catch ({ msg }) {
      if (msg) toast.error(msg as string);
    }
  };

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
        <Picture
          name={`logo-${darkTheme ? "white" : "dark"}`}
          alt={"NCUHOME"}
        />
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
          divider={<Divider orientation={"vertical"} variant="middle" />}
          sx={{
            "& hr": {
              mx: 0.6,
              height: "1rem",
            },
          }}
        >
          <DarkModeSwitch
            onChange={() => setDarkTheme(!darkTheme)}
            checked={darkTheme}
            style={{
              paddingTop: "1px",
              marginLeft: "5px",
              marginRight: "5px",
            }}
            sunColor={muiTheme.palette.action.active}
            size={22}
          />
          <IconButton
            onClick={async () => {
              const ok = await setDialog({
                title: "注销登录",
                content: "使用的身份令牌将被自动销毁",
              });
              if (ok) onLogout();
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
