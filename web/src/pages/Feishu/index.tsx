import { FC } from "react";

import { ThirdPartyCallback } from "@components";

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
