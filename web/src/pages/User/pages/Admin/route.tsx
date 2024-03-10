import { ReactNode } from "react";
import Status from "./Status";

interface Route {
  label: string;
  element: ReactNode;
}

const routes: Route[] = [
  {
    label: "统计数据",
    element: <Status />,
  },
  {
    label: "应用管理",
    element: null,
  },
];

export default routes;
