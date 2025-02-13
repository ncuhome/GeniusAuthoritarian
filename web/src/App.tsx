import useMount from "@hooks/useMount";
import { BrowserRouter, Route, Routes } from "react-router";

import Authoritarian from "@/pages/Authoritarian";
import Home from "@/pages/Authoritarian/Home";
import Feishu from "@/pages/Authoritarian/Feishu";
import DingTalk from "@/pages/Authoritarian/DingTalk";
import Login from "@/pages/Authoritarian/Login";
import PageNotFound from "@components/PageNotFound";

import User from "@/pages/User";
import Navigation from "@/pages/User/Navigation";
import Profile from "@/pages/User/Profile";
import Dev from "@/pages/User/Dev";
import Admin from "@/pages/User/Admin";

import Error from "@/pages/Error";

import Suspense from "@components/Suspense";
import { Box } from "@mui/material";
import { ThemeProvider, createTheme } from "@mui/material/styles";

import dayjs from "dayjs";
import duration from "dayjs/plugin/duration";

dayjs.extend(duration);

export default function App() {
  useMount(() => {
    if (import.meta.env.MODE === "production") {
      console.log(
        atob(
          "ICAgX19fX18gICAgICAgICAgICAgICAgICAKICAvICAgICBcICAgX19fX18gX19fICBfX18KIC8gIFwgLyAgXCAvICAgICBcXCAgXC8gIC8KLyAgICBZICAgIFwgIFkgWSAgXD4gICAgPCAKXF9fX198X18gIC9fX3xffCAgL19fL1xfIFwKICAgICAgICBcLyAgICAgIFwvICAgICAgXC8=",
        ),
      );
    }
  });

  return (
    <ThemeProvider
      theme={createTheme({
        palette: {
          mode: "dark",
        },
      })}
    >
      <Box
        sx={{
          height: "100vh",
          backgroundColor: "#242424",
          colorScheme: "dark",
          color: "text.primary",
        }}
      >
        <BrowserRouter>
          <Routes>
            <Route path={"/"} element={<Authoritarian />}>
              <Route index element={<Home />} />
              <Route path={"feishu"} element={<Feishu />} />
              <Route path={"dingTalk"} element={<DingTalk />} />
              <Route path={"login"} element={<Login />} />
              <Route path={"*"} element={<PageNotFound />} />
            </Route>
            <Route path={"error"} element={<Error />} />
            <Route path={"user"} element={<User />}>
              <Route index element={<Navigation />} />
              <Route path={"profile"} element={<Profile />} />
              <Route path={"dev"} element={<Dev />} />
              <Route path={"admin"} element={<Admin />} />
              <Route path={"*"} element={<PageNotFound />} />
            </Route>
          </Routes>
        </BrowserRouter>
      </Box>
    </ThemeProvider>
  );
}
