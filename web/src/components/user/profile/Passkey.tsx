import { FC, useMemo, useRef } from "react";
import toast from "react-hot-toast";
import {
  coerceToArrayBuffer,
  coerceToBase64Url,
  coerceResponseToBase64Url,
} from "@util/coerce";

import { TransitionGroup } from "react-transition-group";
import PasskeyItem from "./PasskeyItem";
import { List, Box, Collapse, Alert } from "@mui/material";
import { LoadingButton } from "@mui/lab";
import { Add } from "@mui/icons-material";

import useU2F from "@hooks/useU2F";

import { AxiosError } from "axios";
import { apiV1User } from "@api/v1/user/base";

import { useUserApiV1 } from "@api/v1/user/hook";

export const Passkey: FC = () => {
  const { data, mutate } = useUserApiV1<User.Passkey.Cred[]>("passkey/", {
    enableLoading: true,
  });
  const registeredItem = useRef(-1);

  const passkeyAvailable = useMemo(() => !!window.PublicKeyCredential, []);

  const { isLoading: isU2fLoading, refreshToken } = useU2F();

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
    const token = await refreshToken();
    try {
      await apiV1User.delete("passkey/", {
        params: {
          token,
          id: item.id,
        },
      });
      mutate((data) => [...data!.filter((el) => el.id !== item.id)]);
      toast.success("删除成功");
    } catch (err) {
      console.log(err instanceof AxiosError);
      // if (msg) toast.error(msg as any);
    }
  };
  const onRegister = async () => {
    const token = await refreshToken();
    try {
      const {
        data: { data: options },
      } = await apiV1User.get("passkey/register/", {
        params: {
          token,
        },
      });
      options.publicKey.challenge = coerceToArrayBuffer(
        options.publicKey.challenge,
      );
      options.publicKey.user.id = coerceToArrayBuffer(
        options.publicKey.user.id,
      );
      const credential = await navigator.credentials.create(options);
      if (
        Object.prototype.toString.call(credential) !==
        "[object PublicKeyCredential]"
      ) {
        toast.error(`创建凭据失败，凭据类型不正确`);
        return;
      }
      const pubKeyCred = credential as PublicKeyCredential;
      const {
        data: { data: newItem },
      } = await apiV1User.post("passkey/register/", {
        id: pubKeyCred.id,
        authenticatorAttachment: pubKeyCred.authenticatorAttachment,
        rawId: coerceToBase64Url(pubKeyCred.rawId),
        response: coerceResponseToBase64Url(pubKeyCred.response),
        type: pubKeyCred.type,
      });
      registeredItem.current = newItem.id;
      mutate((data) => {
        if (
          !data ||
          data.length === 0 ||
          data.findIndex((item) => item.id === newItem.id) != -1
        )
          return data;
        return [newItem, ...data];
      });
    } catch (err: any) {
      if (err instanceof AxiosError) {
        err = err as ApiError<void>;
        if (err.msg) toast.error(err.msg);
      } else {
        if (err.name != "NotAllowedError") toast.error(`创建凭据失败: ${err}`);
      }
    }
  };

  return (
    <>
      {passkeyAvailable ? (
        <LoadingButton
          variant="outlined"
          loading={isU2fLoading}
          startIcon={<Add />}
          onClick={onRegister}
        >
          添加通行密钥
        </LoadingButton>
      ) : (
        <Alert severity="warning">此设备不支持通行密钥认证</Alert>
      )}

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
              {data.map((item) => (
                <Collapse key={item.id}>
                  <PasskeyItem
                    item={item}
                    onRename={(name) => onRename(item.id, name)}
                    onDelete={() => onDelete(item)}
                    enableRename={item.id === registeredItem.current}
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
