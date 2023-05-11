import { FC, Suspense as ReactSuspense, PropsWithChildren } from "react";

import { CircularProgress, Stack } from "@mui/material";

export const Suspense: FC<PropsWithChildren> = ({ children }) => {
  return (
    <ReactSuspense
      fallback={
        <Stack
          sx={{
            height: "100%",
            width: "100%",
          }}
          justifyContent={"center"}
          alignItems={"center"}
        >
          <CircularProgress size={50} />
        </Stack>
      }
    >
      {children}
    </ReactSuspense>
  );
};
export default Suspense;
