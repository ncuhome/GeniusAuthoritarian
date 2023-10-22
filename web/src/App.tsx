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
    [isDarkTheme]
  );

  useMount(() => {
    if (import.meta.env.MODE === "production") {
      console.log(
        atob(
          "ICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICBfX19fXyAgICAgICAgICAgICAgICAgCiBfICBfICBfX18gXyAgIF8gXyAgXyAgX19fICBfXyAgX18gX19fICAgICAvICAgICBcICAgX19fX18gX19fICBfX18KfCBcfCB8LyBfX3wgfCB8IHwgfHwgfC8gXyBcfCAgXC8gIHwgX198ICAgLyAgXCAvICBcIC8gICAgIFxcICBcLyAgLwp8IC5gIHwgKF9ffCB8X3wgfCBfXyB8IChfKSB8IHxcL3wgfCBffCAgIC8gICAgWSAgICBcICBZIFkgIFw+ICAgIDwgCnxffFxffFxfX198XF9fXy98X3x8X3xcX19fL3xffCAgfF98X19ffCAgXF9fX198X18gIC9fX3xffCAgL19fL1xfIFwKICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgIFwvICAgICAgXC8gICAgICBcLw=="
        )
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
