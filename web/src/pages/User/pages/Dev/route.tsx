import { ReactNode } from "react";

import App from "./app";
import Ssh from "./ssh";

interface Route {
  label: string;
  description: string;
  element: ReactNode;
}

const routes: Route[] = [
  {
    label: "应用管理",
    description:
      "应用一经创建，就会显示在导航页面，请谨慎操作。应用密钥仅在创建时显示，请妥善保管。对接状态仅作用于显示，不影响系统逻辑",
    element: <App />,
  },
  {
    label: "SSH 密钥",
    description: "请通过 28 节点 222 端口连接，该账号仅用于端口转发，严禁分享",
    element: <Ssh />,
  },
];

export default routes;
