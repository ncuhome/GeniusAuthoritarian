import { ReactNode } from "react";

interface Route {
  label: string;
  element: ReactNode;
}

const routes: Route[] = [
  {
    label: "统计数据",
    element: null,
  },
  {
    label: "应用管理",
    element: null,
  },
];

export default routes;
