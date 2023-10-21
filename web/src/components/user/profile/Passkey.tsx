import { FC } from "react";

import { Button, ButtonGroup, Collapse, List } from "@mui/material";
import { Add } from "@mui/icons-material";

import { useUserApiV1 } from "@api/v1/user/hook";
import { apiV1User } from "@api/v1/user/base";

export const Passkey: FC = () => {
  const { data, mutate } = useUserApiV1("passkey/", {
    enableLoading: true,
  });

  const onRegister = async () => {};

  return (
    <>
      <ButtonGroup variant={"outlined"}>
        <Button startIcon={<Add />} onClick={onRegister}>
          添加通行密钥
        </Button>
      </ButtonGroup>
    </>
  );
};

export default Passkey;
