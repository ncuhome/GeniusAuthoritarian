import { FC, PropsWithChildren } from "react";

import BlockTitle from "@components/user/BlockTitle";
import { Paper, Box, Typography, SxProps } from "@mui/material";

interface Props extends PropsWithChildren {
  title?: string;
  subtitle?: string;
  sx?: SxProps;
}

export const Block: FC<Props> = ({ title, subtitle, children, sx }) => {
  return (
    <Box component={Paper} elevation={5} sx={sx}>
      {title ? <BlockTitle>{title}</BlockTitle> : null}
      {subtitle ? (
        <Typography variant={"subtitle2"} color={"text.secondary"}>
          {subtitle}
        </Typography>
      ) : null}

      {children}
    </Box>
  );
};
export default Block;
