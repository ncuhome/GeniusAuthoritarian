import { FC } from "react";

import { Button, ButtonGroup, Collapse, List } from "@mui/material";
import { Add } from "@mui/icons-material";

import { useUserApiV1 } from "@api/v1/user/hook";
import { apiV1User } from "@api/v1/user/base";

interface Props {
  mfaEnabled?: boolean;
}

export const Passkey: FC<Props> = ({ mfaEnabled }) => {
  const { data, mutate } = useUserApiV1(mfaEnabled ? "passkey/" : null, {
    enableLoading: true,
  });

  const onRegister = async () => {};

  return (
    <>
      <ButtonGroup variant={"outlined"}>
        <Button startIcon={<Add />} onClick={onRegister} disabled={!mfaEnabled}>
          {mfaEnabled === false ? "需要启用双因素认证" : "添加通行密钥"}
        </Button>
      </ButtonGroup>
    </>
  );
};

export default Passkey;
