import { useMemo } from "react";
import { useSearchParams } from "react-router-dom";

export const useQuery = <
  V extends {
    toString(): string;
  }
>(
  key: string,
  initialValue: V
): [V, (v: V) => void] => {
  const [searchParams, setSearchParams] = useSearchParams();
  return [
    useMemo(
      () => (searchParams.get(key) as unknown as V) || initialValue,
      [searchParams]
    ),
    (v: V) =>
      setSearchParams((prev) => {
        prev.set(key, v.toString());
        return prev;
      }),
  ];
};
