import { FC } from "react";

import ThirdPartyCallback from "@components/auth/ThirdPartyCallback";

export const DingTalk: FC = () => {
  return (
    <ThirdPartyCallback
      keyCode={"authCode"}
      keyAppCode={"state"}
      thirdParty={"dingTalk"}
    />
  );
};
export default DingTalk;
