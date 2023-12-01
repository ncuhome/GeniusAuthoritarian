import { FC, PropsWithChildren } from "react";

import BlockTitle from "@components/user/BlockTitle";
import { Box, SxProps } from "@mui/material";

interface Props extends PropsWithChildren {
  title: string;
  sx?: SxProps;
}

export const BlockArea: FC<Props> = ({ title, children, sx }) => {
  return (
    <Box className={"block-area"} sx={sx}>
      <BlockTitle
        sx={{
          marginBottom: "0.8rem",
        }}
      >
        {title}
      </BlockTitle>
      {children}
    </Box>
  );
};
export default BlockArea;
