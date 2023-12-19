import { FC, useState, CSSProperties, Fragment } from "react";
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
  Drawer,
  List,
  ListItem,
  ListItemIcon,
  ListItemButton,
  ListItemText,
  useTheme as useMuiTheme,
  useMediaQuery,
} from "@mui/material";
import {
  LogoutRounded,
  Menu,
  KeyboardArrowRightRounded,
} from "@mui/icons-material";

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
  const isSmallScreen = useMediaQuery(muiTheme.breakpoints.down("sm"));

  const [menuOpen, setMenuOpen] = useState(false);

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

  const onRoute = (index: number, path: string) => {
    handleChangeTab(index);
    nav(path);
  };

  const renderLogo = (style?: CSSProperties) => {
    return (
      <Picture
        name={`logo-${darkTheme ? "white" : "dark"}`}
        alt={"NCUHOME"}
        imgStyle={style}
      />
    );
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
      {isSmallScreen ? undefined : (
        <Box
          sx={{
            height: "100%",
            display: "flex",
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
          {renderLogo()}
        </Box>
      )}

      <Stack
        flexDirection={"row"}
        flexGrow={1}
        justifyContent={"space-between"}
      >
        {isSmallScreen ? (
          <IconButton
            color="inherit"
            aria-label="open drawer"
            edge="start"
            onClick={() => setMenuOpen((open) => !open)}
          >
            <Menu />
          </IconButton>
        ) : (
          <Tabs value={currentTab} textColor="inherit">
            {routers.map((r, index) => (
              <Tab
                key={r.name}
                label={<Typography variant={"subtitle1"}>{r.name}</Typography>}
                value={index}
                onClick={() => onRoute(index, r.path)}
                disableRipple
              />
            ))}
          </Tabs>
        )}

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

      {isSmallScreen ? (
        <Drawer
          anchor={"left"}
          open={menuOpen}
          onClose={() => setMenuOpen(false)}
        >
          <Stack
            alignItems={"center"}
            sx={{
              boxSizing: "border-box",
              pt: 2,
              px: 3,
              minWidth: "100%",
            }}
          >
            {renderLogo({
              maxWidth: "100%",
              width: "7.5rem",
            })}
          </Stack>
          <List>
            {routers.map((r, index) => (
              <Fragment key={r.name}>
                <Divider />
                <ListItem disablePadding>
                  <ListItemButton
                    onClick={() => {
                      onRoute(index, r.path);
                      setMenuOpen(false);
                    }}
                    sx={{
                      paddingRight: "1.5rem",
                    }}
                  >
                    <ListItemIcon
                      sx={{
                        minWidth: "2.2rem",
                      }}
                    >
                      <KeyboardArrowRightRounded
                        sx={{
                          color: "text.disabled",
                        }}
                      />
                    </ListItemIcon>
                    <ListItemText primary={r.name} />
                  </ListItemButton>
                </ListItem>
              </Fragment>
            ))}
          </List>
        </Drawer>
      ) : undefined}
    </Stack>
  );
};
export default NavHeader;
