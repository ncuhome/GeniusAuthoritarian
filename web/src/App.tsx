import { BrowserRouter, Route, Routes } from "react-router-dom";
import { Toaster } from "react-hot-toast";

import {Home, Error, Feishu, DingTalk} from "./pages";

import { ThemeProvider, createTheme } from "@mui/material/styles";
const darkTheme = createTheme({
  palette: {
    mode: "dark",
  },
});

export default function App() {
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
              <Route path={"feishu"} element={<Feishu/>}/>
              <Route path={"dingTalk"} element={<DingTalk/>}/>
          </Routes>
        </BrowserRouter>
      </ThemeProvider>
    </>
  );
}
