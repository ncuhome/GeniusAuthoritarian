import { FC } from "react";
import { Typography, TypographyProps } from "@mui/material";

export const BlockTitle: FC<TypographyProps> = ({ children, ...props }) => {
  return (
    <Typography
      variant={"h5"}
      fontWeight={"bold"}
      {...props}
    >
      {children}
    </Typography>
  );
};
export default BlockTitle;
