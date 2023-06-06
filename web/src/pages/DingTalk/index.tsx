import { FC } from "react";

import { ThirdPartyCallback } from "@components";

export const DingTalk: FC = () => {
  return (
    <ThirdPartyCallback
      keyCode={"authCode"}
      keyAppCode={"state"}
      thirdParty={"dingTalk"}
    />
  );
};
