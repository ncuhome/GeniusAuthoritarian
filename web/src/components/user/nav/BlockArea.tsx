import { FC, PropsWithChildren } from "react";

import BlockTitle from "@components/user/BlockTitle";
import { Box } from "@mui/material";

interface Props extends PropsWithChildren {
  title: string;
}

export const BlockArea: FC<Props> = ({ title, children }) => {
  return (
    <Box className={"block-area"}>
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
