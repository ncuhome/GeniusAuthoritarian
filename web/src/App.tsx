import { useMount } from "@hooks";
import { BrowserRouter, Route, Routes } from "react-router-dom";
import { Toaster } from "react-hot-toast";

import { Home, Error, Feishu, DingTalk, Login, User } from "./pages";

import { ThemeProvider, createTheme } from "@mui/material/styles";
const darkTheme = createTheme({
  palette: {
    mode: "dark",
  },
});

export default function App() {
  useMount(() => {
    // Ncuhome
    console.log(
      atob(
        "IF8gIF8gIF9fXyBfICAgXyBfICBfICBfX18gIF9fICBfXyBfX18gCnwgXHwgfC8gX198IHwgfCB8IHx8IHwvIF8gXHwgIFwvICB8IF9ffAp8IC5gIHwgKF9ffCB8X3wgfCBfXyB8IChfKSB8IHxcL3wgfCBffCAKfF98XF98XF9fX3xcX19fL3xffHxffFxfX18vfF98ICB8X3xfX198"
      )
    );
    // Mmx
    console.log(
      atob(
        "ICAgX19fX18gICAgICAgICAgICAgICAgICAKICAvICAgICBcICAgX19fX18gX19fICBfX18KIC8gIFwgLyAgXCAvICAgICBcXCAgXC8gIC8KLyAgICBZICAgIFwgIFkgWSAgXD4gICAgPCAKXF9fX198X18gIC9fX3xffCAgL19fL1xfIFwKICAgICAgICBcLyAgICAgIFwvICAgICAgXC8="
      )
    );
  });
  return (
    <>
      <Toaster
        toastOptions={{
          style: {
            borderRadius: "20px",
            background: "#2f2f2f",
            color: "#fff",
          },
        }}
      />
      <ThemeProvider theme={darkTheme}>
        <BrowserRouter>
          <Routes>
            <Route index element={<Home />} />
            <Route path={"error"} element={<Error />} />
            <Route path={"feishu"} element={<Feishu />} />
            <Route path={"dingTalk"} element={<DingTalk />} />
            <Route path={"login"} element={<Login />} />
            <Route path={"user/*"} element={<User />} />
          </Routes>
        </BrowserRouter>
      </ThemeProvider>
    </>
  );
}
