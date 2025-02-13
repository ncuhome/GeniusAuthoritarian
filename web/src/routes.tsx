import { lazy } from "react";
import { RouteObject } from "react-router";

const Authoritarian = lazy(() => import("@/pages/Authoritarian"));
const Home = lazy(() => import("@/pages/Authoritarian/Home"));
const Feishu = lazy(() => import("@/pages/Authoritarian/Feishu"));
const DingTalk = lazy(() => import("@/pages/Authoritarian/DingTalk"));
const Login = lazy(() => import("@/pages/Authoritarian/Login"));

const User = lazy(() => import("@/pages/User"));
const Navigation = lazy(() => import("@/pages/User/Navigation"));
const Profile = lazy(() => import("@/pages/User/Profile"));
const Dev = lazy(() => import("@/pages/User/Dev"));
const Admin = lazy(() => import("@/pages/User/Admin"));

const PageNotFound = lazy(() => import("@components/PageNotFound"));
const Error = lazy(() => import("@/pages/Error"));

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
