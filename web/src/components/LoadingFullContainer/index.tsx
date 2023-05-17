import { FC } from "react";

import { CircularProgress, Stack } from "@mui/material";

export const LoadingFullContainer: FC = () => {
  return (
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
  );
};
export default LoadingFullContainer;
