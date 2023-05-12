import { FC, PropsWithChildren } from "react";

import { Paper, Box } from "@mui/material";

export const Block: FC<PropsWithChildren> = ({ children }) => {
  return (
    <Box component={Paper} elevation={5}>
      {children}
    </Box>
  );
};
export default Block;
