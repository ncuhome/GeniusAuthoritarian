import { FC } from "react";
import useMfaCode from "@hooks/useMfaCode";
import toast from "react-hot-toast";
import {
  coerceToArrayBuffer,
  coerceToBase64Url,
  coerceResponseToBase64Url,
} from "@util/coerce";

import { TransitionGroup } from "react-transition-group";
import PasskeyItem from "./PasskeyItem";
import { Button, ButtonGroup, List, Box, Collapse } from "@mui/material";
import { Add } from "@mui/icons-material";

import { apiV1User } from "@api/v1/user/base";

import { useUserApiV1 } from "@api/v1/user/hook";

interface Props {
  mfaEnabled?: boolean;
}

export const Passkey: FC<Props> = ({ mfaEnabled }) => {
  const { data, mutate } = useUserApiV1<User.Passkey.Cred[]>(
    mfaEnabled ? "passkey/" : null,
    {
      enableLoading: true,
    }
  );

  const openMfaDialog = useMfaCode();

  const onRename = async (id: number, name: string) => {
    if (name.length === 0) {
      toast.error("新名称不能为空");
      return;
    }
    if (name.length > 15) {
      toast.error("名称最高 15 字");
      return;
    }
    try {
      await apiV1User.patch("passkey/", {
        id,
        name,
      });
      const index = data!.findIndex((el) => el.id === id);
      data![index].name = name;
      mutate([...data!]);
    } catch ({ msg }) {
      if (msg) toast.error(msg as any);
    }
  };
  const onDelete = async (item: User.Passkey.Cred) => {
    const code = await openMfaDialog();
    try {
      await apiV1User.delete("passkey/", {
        params: {
          code,
          id: item.id,
        },
      });
      mutate((data) => [...data!.filter((el) => el.id !== item.id)]);
      toast.success("删除成功");
    } catch ({ msg }) {
      if (msg) toast.error(msg as any);
    }
  };
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
      options.publicKey.challenge = coerceToArrayBuffer(
        options.publicKey.challenge
      );
      options.publicKey.user.id = coerceToArrayBuffer(
        options.publicKey.user.id
      );

      try {
        const credential = await navigator.credentials.create(options);
        if (!(credential instanceof PublicKeyCredential)) {
          toast.error(`创建凭据失败，凭据类型不正确`);
          return;
        }
        const pubKeyCred = credential as PublicKeyCredential;
        try {
          const {
            data: { data: newItem },
          } = await apiV1User.post("passkey/register/", {
            id: pubKeyCred.id,
            authenticatorAttachment: pubKeyCred.authenticatorAttachment,
            rawId: coerceToBase64Url(pubKeyCred.rawId),
            response: coerceResponseToBase64Url(pubKeyCred.response),
            type: pubKeyCred.type,
          });
          mutate((data) => {
            if (
              !data ||
              data.length === 0 ||
              data.findIndex((item) => item.id === newItem.id) != -1
            )
              return data;
            return [newItem, ...data];
          });
        } catch ({ msg }) {
          if (msg) toast.error(msg as any);
        }
      } catch (err: any) {
        if (err.name != "NotAllowedError") toast.error(`创建凭据失败: ${err}`);
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

      {data ? (
        <Box
          sx={{
            maxWidth: "100%",
            overflowY: "auto",
          }}
        >
          <List
            sx={{
              minWidth: "27rem",
            }}
          >
            <TransitionGroup>
              {data.map((item, index) => (
                <Collapse key={item.id}>
                  <PasskeyItem
                    item={item}
                    onRename={(name) => onRename(item.id, name)}
                    onDelete={() => onDelete(item)}
                  />
                </Collapse>
              ))}
            </TransitionGroup>
          </List>
        </Box>
      ) : undefined}
    </>
  );
};

export default Passkey;
