import { FC } from "react";

import { CircularProgress, Stack, SxProps } from "@mui/material";

interface Props {
  sx?: SxProps;
}

export const LoadingFullContainer: FC<Props> = ({ sx }) => {
  return (
    <Stack
      sx={{
        height: "100%",
        width: "100%",
        ...sx,
      }}
      justifyContent={"center"}
      alignItems={"center"}
    >
      <CircularProgress size={50} />
    </Stack>
  );
};
export default LoadingFullContainer;
