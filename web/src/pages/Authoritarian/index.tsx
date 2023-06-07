import { FC } from "react";
import { Route, Routes } from "react-router-dom";

import Home from "./pages/Home";
import Feishu from "./pages/Feishu";
import DingTalk from "./pages/DingTalk";
import Login from "./pages/Login";

import PageNotFound from "@components/PageNotFound";

export const Authoritarian: FC = () => {
  return (
    <Routes>
      <Route index element={<Home />} />
      <Route path={"feishu"} element={<Feishu />} />
      <Route path={"dingTalk"} element={<DingTalk />} />
      <Route path={"login"} element={<Login />} />
      <Route path={"*"} element={<PageNotFound />} />
    </Routes>
  );
};
export default Authoritarian;
