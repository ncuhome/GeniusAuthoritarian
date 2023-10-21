import { FC, useState } from "react";
import useMfaCode from "@hooks/useMfaCode";
import toast from "react-hot-toast";

import {
  Button,
  ButtonGroup,
  Dialog,
  DialogContent,
  DialogActions,
  List,
} from "@mui/material";
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

  const openMfaDialog = useMfaCode();

  const onRegister = async () => {
    const code = await openMfaDialog();
    try {
      const {
        data: { data: options },
      } = await apiV1User.get("passkey/register/", {
        params: {
          code,
        },
      });
      const encoder = new TextEncoder();
      options.publicKey.challenge = encoder.encode(options.publicKey.challenge);
      options.publicKey.user.id = encoder.encode(options.publicKey.user.id);

      try {
        const credential = await navigator.credentials.create(options);
        if (Credential === null) {
          toast.error(`创建凭据失败，凭据为 null`);
          return;
        }
        try {
          await apiV1User.post("passkey/register/", credential);
        } catch ({ msg }) {
          if (msg) toast.error(msg as any);
        }
      } catch (err) {
        toast.error(`创建凭据失败: ${err}`);
      }
    } catch ({ msg }) {
      if (msg) toast.error(msg as any);
    }
  };

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
