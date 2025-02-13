import { FC, PropsWithChildren } from "react";
import "./Block.scss";

import BlockTitle from "@components/user/BlockTitle";
import { Paper, Box, Typography, SxProps } from "@mui/material";

interface Props extends PropsWithChildren {
  title?: string;
  subtitle?: string;
  sx?: SxProps;
  disablePadding?: boolean;
}

export const Block: FC<Props> = ({
  title,
  subtitle,
  children,
  sx,
  disablePadding,
}) => {
  if (disablePadding)
    sx = {
      ...sx,
      padding: "unset!important",
    };

  return (
    <Box component={Paper} elevation={3} sx={sx} className={"user-block"}>
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
