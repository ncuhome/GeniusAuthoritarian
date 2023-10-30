import { useCallback, useMemo, useRef, useState } from "react";
import toast from "react-hot-toast";

import { useUserApiV1 } from "@api/v1/user/hook";
import { apiV1User } from "@api/v1/user/base";

import { shallow } from "zustand/shallow";
import useU2fDialog from "@store/useU2fDialog";

interface U2F {
  isLoading: boolean;
  setPrefer: (method: User.U2F.Methods) => Promise<void>;
  refreshToken: (tip?: string) => Promise<string>;
}

export const useU2F = (): U2F => {
  const [loadData, setLoadData] = useState(false);
  const dataPromise = useRef<null | Promise<void>>(null);
  const dataResolver = useRef<() => void>(() => {});
  const dataLoaded = useRef(false);

  const { isLoading } = useUserApiV1<User.U2F.Status>(
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

  const setPrefer = async (method: User.U2F.Methods) => {
    await apiV1User.put("u2f/prefer", {
      prefer: method,
    });
    useU2fDialog.getState().setPrefer(method);
  };

  const refreshToken = useCallback(async (tip?: string) => {
    if (dataPromise.current) await dataPromise.current;
    if (!dataLoaded.current) {
      dataPromise.current = new Promise((resolve) => {
        dataResolver.current = resolve;
      });
      setLoadData(true);
      await dataPromise.current;
    }

    const result = await useU2fDialog.getState().openDialog(tip);
    useU2fDialog.getState().setToken(result);
    return result.token;
  }, []);

  return {
    isLoading,
    setPrefer,
    refreshToken,
  };
};

export default useU2F;
