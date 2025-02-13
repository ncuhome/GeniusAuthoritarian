import { FC } from "react";
import { Route, Routes, Outlet } from "react-router";

import Home from "./Home";
import Feishu from "./Feishu";
import DingTalk from "./DingTalk";
import Login from "./Login";

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
      <Outlet />
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
