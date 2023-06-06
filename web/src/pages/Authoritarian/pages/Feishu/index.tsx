import { FC } from "react";

import ThirdPartyCallback from "@components/auth/ThirdPartyCallback";
export const Feishu: FC = () => {
  return (
    <ThirdPartyCallback
      keyCode={"code"}
      keyAppCode={"state"}
      thirdParty={"feishu"}
    />
  );
};
export default Feishu;
