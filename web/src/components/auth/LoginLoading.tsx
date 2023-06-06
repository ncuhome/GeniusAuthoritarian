import { FC } from "react";

import { Stack, CircularProgress, Typography } from "@mui/material";

export const LoginLoading: FC = () => {
  return (
    <Stack
      sx={{
        height: "100%",
        width: "100%",
      }}
      justifyContent={"center"}
      alignItems={"center"}
      spacing={5}
    >
      <CircularProgress size={"6rem"} />

      <Typography
        variant={"h6"}
        color={"text.primary"}
        sx={{
          letterSpacing: "0.15em",
        }}
      >
        正在跳转
      </Typography>
    </Stack>
  );
};
export default LoginLoading;
