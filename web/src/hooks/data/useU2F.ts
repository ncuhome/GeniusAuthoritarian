import { useCallback, useEffect, useRef, useState } from "react";

import { useUserApiV1 } from "@api/v1/user/hook";

import useU2fDialog from "@store/useU2fDialog";

interface U2F {
  isLoading: boolean;
  refreshToken: (tip?: string) => Promise<string>;
}

export const useU2F = (): U2F => {
  const [loadData, setLoadData] = useState(false);
  const dataPromise = useRef<null | Promise<void>>(null);
  const dataResolver = useRef<() => void>(() => {});
  const dataLoaded = useRef(false);

  const onDataLoaded = useCallback((data: User.U2F.Status) => {
    useU2fDialog.getState().setStatus(data);
    dataLoaded.current = true;
    dataResolver.current();
  }, []);

  const { data, isLoading } = useUserApiV1<User.U2F.Status>(
    loadData ? "u2f/" : null,
    {
      enableLoading: true,
      revalidateIfStale: false,
      revalidateOnFocus: false,
      onSuccess: onDataLoaded,
    },
  );

  useEffect(() => {
    if (data) onDataLoaded(data);
  }, [data]);

  const refreshToken = useCallback(async (tip?: string) => {
    if (dataPromise.current) await dataPromise.current;
    if (!dataLoaded.current) {
      dataPromise.current = new Promise((resolve) => {
        dataResolver.current = resolve;
      });
      setLoadData(true);
      await dataPromise.current;
      dataPromise.current = null;
    }

    const result = await useU2fDialog.getState().openDialog(tip);
    useU2fDialog.getState().setToken(result);
    return result.token;
  }, []);

  return {
    isLoading,
    refreshToken,
  };
};

export default useU2F;
