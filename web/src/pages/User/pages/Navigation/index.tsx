import { FC } from "react";

import { Stack, Typography } from "@mui/material";

export const Navigation: FC = () => {
  return (
    <Stack justifyContent={"center"} alignItems={"center"} height={"100%"}>
      <Typography variant={"h5"} fontWeight={"bold"} sx={{ opacity: 0.5 }}>
        别急，马上写
      </Typography>
    </Stack>
  );
};
export default Navigation;
