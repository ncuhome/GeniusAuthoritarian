import { FC, PropsWithChildren } from "react";

import { BlockTitle } from "@/pages/User/components";
import { Box } from "@mui/material";

interface Props extends PropsWithChildren {
  title: string;
}

export const BlockArea: FC<Props> = ({ title, children }) => {
  return (
    <Box>
      <BlockTitle>{title}</BlockTitle>
      {children}
    </Box>
  );
};
export default BlockArea;
