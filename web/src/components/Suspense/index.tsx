import { FC, Suspense as ReactSuspense, PropsWithChildren } from "react";

import { LoadingFullContainer } from "@/components";

export const Suspense: FC<PropsWithChildren> = ({ children }) => {
  return (
    <ReactSuspense fallback={<LoadingFullContainer />}>
      {children}
    </ReactSuspense>
  );
};
export default Suspense;
