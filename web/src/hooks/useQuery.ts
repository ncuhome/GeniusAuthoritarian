import { useCallback, useMemo } from "react";
import { useSearchParams } from "react-router";

export const createUseQuery = () => {
  const [searchParams, setSearchParams] = useSearchParams();
  return useCallback(
    <
      V extends {
        toString(): string;
      }
    >(
      key: string,
      initialValue: V
    ): [V, (v: V) => void] => {
      return [
        useMemo(
          () => (searchParams.get(key) as unknown as V) || initialValue,
          [searchParams]
        ),
        useCallback(
          (v: V) =>
            setSearchParams((prev) => {
              prev.set(key, v.toString());
              return prev;
            }),
          [setSearchParams]
        ),
      ];
    },
    [searchParams, setSearchParams]
  );
};
