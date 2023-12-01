import { lazy, useMemo } from "react";
import useMount from "@hooks/useMount";
import { BrowserRouter, Route, Routes } from "react-router-dom";

const User = lazy(() => import("./pages/User"));
const Authoritarian = lazy(() => import("./pages/Authoritarian"));
import Error from "@/pages/Error";

import Suspense from "@components/Suspense";
import { Box } from "@mui/material";
import { ThemeProvider, createTheme } from "@mui/material/styles";

import useTheme from "@store/useTheme";

import dayjs from "dayjs";
import duration from "dayjs/plugin/duration";

dayjs.extend(duration);

export default function App() {
  const isDarkTheme = useTheme((state) => state.dark);

  const theme = useMemo(
    () =>
      createTheme({
        palette: {
          mode: isDarkTheme ? "dark" : undefined,
        },
      }),
    [isDarkTheme],
  );

  useMount(() => {
    if (import.meta.env.MODE === "production") {
      console.log(
        atob(
          "ICAgX19fX18gICAgICAgICAgICAgICAgICAKICAvICAgICBcICAgX19fX18gX19fICBfX18KIC8gIFwgLyAgXCAvICAgICBcXCAgXC8gIC8KLyAgICBZICAgIFwgIFkgWSAgXD4gICAgPCAKXF9fX198X18gIC9fX3xffCAgL19fL1xfIFwKICAgICAgICBcLyAgICAgIFwvICAgICAgXC8=",
        ),
      );
      console.log(
        "%c可恶，看板娘都要我自己画！！！",
        "margin-top: 8px;font-size: 13.5px",
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
            <Route path={"error"} element={<Error />} />
            <Route
              path={"user/*"}
              element={
                <ThemeProvider theme={theme}>
                  <Suspense>
                    <User />
                  </Suspense>
                </ThemeProvider>
              }
            />
            <Route
              path={"*"}
              element={
                <Suspense>
                  <Authoritarian />
                </Suspense>
              }
            />
          </Routes>
        </BrowserRouter>
      </Box>
    </ThemeProvider>
  );
}
