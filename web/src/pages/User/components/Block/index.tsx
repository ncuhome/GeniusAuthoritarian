import { FC, PropsWithChildren } from "react";

import { BlockTitle } from "@/pages/User/components";
import { Paper, Box, Typography } from "@mui/material";

interface Props extends PropsWithChildren {
  title?: string;
  subtitle?: string;
}

export const Block: FC<Props> = ({ title, subtitle, children }) => {
  return (
    <Box component={Paper} elevation={5}>
      {title ? <BlockTitle>{title}</BlockTitle> : null}
      {subtitle ? (
        <Typography variant={"subtitle2"}>{subtitle}</Typography>
      ) : null}

      {children}
    </Box>
  );
};
export default Block;
