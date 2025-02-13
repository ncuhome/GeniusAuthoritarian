import { RouteObject } from "react-router";

import Authoritarian from "@/pages/Authoritarian";
import Home from "@/pages/Authoritarian/Home";
import Feishu from "@/pages/Authoritarian/Feishu";
import DingTalk from "@/pages/Authoritarian/DingTalk";
import Login from "@/pages/Authoritarian/Login";

import User from "@/pages/User";
import Navigation from "@/pages/User/Navigation";
import Profile from "@/pages/User/Profile";
import Dev from "@/pages/User/Dev";
import Admin from "@/pages/User/Admin";

import PageNotFound from "@components/PageNotFound";
import Error from "@/pages/Error";

const routes: RouteObject[] = [
  {
    path: "/",
    element: <Authoritarian />,
    children: [
      {
        index: true,
        element: <Home />,
      },
      {
        path: "feishu",
        element: <Feishu />,
      },
      {
        path: "dingTalk",
        element: <DingTalk />,
      },
      {
        path: "login",
        element: <Login />,
      },
      {
        path: "*",
        element: <PageNotFound />,
      },
    ],
  },
  {
    path: "user",
    element: <User />,
    children: [
      {
        index: true,
        element: <Navigation />,
      },
      {
        path: "profile",
        element: <Profile />,
      },
      {
        path: "dev",
        element: <Dev />,
      },
      {
        path: "admin",
        element: <Admin />,
      },
      {
        path: "*",
        element: <PageNotFound />,
      },
    ],
  },
  {
    path: "error",
    element: <Error />,
  },
];

export default routes;
