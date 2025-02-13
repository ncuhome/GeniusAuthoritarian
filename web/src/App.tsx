import useMount from "@hooks/useMount";
import { createBrowserRouter, RouterProvider } from "react-router";
import routes from "@/routes";

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
        <Suspense>
          <RouterProvider router={createBrowserRouter(routes)} />
        </Suspense>
      </Box>
    </ThemeProvider>
  );
}
