import { FC, PropsWithChildren } from "react";

import { Box, Divider, Typography } from "@mui/material";

interface Props extends PropsWithChildren {
  title: string;
  desc?: string;
}

export const ChildBlock: FC<Props> = ({ title, desc, children }) => {
  return (
    <Box className={"user-child-block"}>
      <Typography variant={"subtitle1"}>{title}</Typography>
      {desc ? (
        <>
          <Divider />
          <Typography variant={"body2"}>{desc}</Typography>
        </>
      ) : undefined}
      <Box mt={1.3}>{children}</Box>
    </Box>
  );
};

export default ChildBlock;
