import { useCallback, useMemo, useRef, useState } from "react";
import toast from "react-hot-toast";

import { useUserApiV1 } from "@api/v1/user/hook";

import { shallow } from "zustand/shallow";
import useU2fDialog from "@store/useU2fDialog";

interface U2F {
  available: boolean;
  isLoading: boolean;
  setPrefer: (method: User.U2F.Methods) => Promise<void>;
  refreshToken: () => Promise<User.U2F.Result>;
}

export const useU2F = (): U2F => {
  const [loadData, setLoadData] = useState(false);
  const dataPromise = useRef<null | Promise<void>>(null);
  const dataResolver = useRef<() => void>(() => {});
  const dataLoaded = useRef(false);

  const { phone, mfa, passkey } = useU2fDialog(
    (state) => ({
      phone: state.phone,
      mfa: state.mfa,
      passkey: state.passkey,
    }),
    shallow
  );

  const { data, isLoading } = useUserApiV1<User.U2F.Status>(
    loadData ? "u2f/" : null,
    {
      revalidateIfStale: false,
      revalidateOnFocus: false,
      onSuccess: (data) => {
        useU2fDialog.getState().setStatus(data);
        dataLoaded.current = true;
        dataResolver.current();
      },
      onError: (err) => {
        toast.error(`载入 U2F 状态失败: ${err}`);
      },
    }
  );

  const available = useMemo(
    // 未加载时固定为 true
    () => !data || phone || passkey || mfa,
    [phone, passkey, mfa]
  );

  const setPrefer = async (method: User.U2F.Methods) => {
    // todo 接口未编写
    // useU2fDialog.getState().setPrefer()
  };

  const refreshToken = useCallback(async () => {
    const u2fStatus = useU2fDialog.getState().u2f;
    if (u2fStatus && u2fStatus.valid_before > new Date().getTime() / 1000) {
      return u2fStatus;
    }

    if (dataPromise.current) await dataPromise.current;
    if (!dataLoaded.current) {
      dataPromise.current = new Promise((resolve) => {
        dataResolver.current = resolve;
      });
      setLoadData(true);
      await dataPromise.current;
    }

    const result = await useU2fDialog.getState().openDialog();
    useU2fDialog.getState().setToken(result);
    return result;
  }, []);

  return {
    available,
    isLoading,
    setPrefer,
    refreshToken,
  };
};

export default useU2F;
