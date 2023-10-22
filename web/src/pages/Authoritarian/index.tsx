import { FC } from "react";
import { Route, Routes } from "react-router-dom";

import Home from "./pages/Home";
import Feishu from "./pages/Feishu";
import DingTalk from "./pages/DingTalk";
import Login from "./pages/Login";

import { Toaster } from "react-hot-toast";
import PageNotFound from "@components/PageNotFound";

export const Authoritarian: FC = () => {
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
      <Routes>
        <Route index element={<Home />} />
        <Route path={"feishu"} element={<Feishu />} />
        <Route path={"dingTalk"} element={<DingTalk />} />
        <Route path={"login"} element={<Login />} />
        <Route path={"*"} element={<PageNotFound />} />
      </Routes>
    </>
  );
};
export default Authoritarian;
