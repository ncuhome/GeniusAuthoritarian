import { FC } from "react";

import { Button, Tooltip, ButtonProps } from "@mui/material";

interface Props {
  title: string;
}

export const TipButton: FC<ButtonProps & Props> = ({
  title,
  children,
  ...rest
}) => {
  return (
    <Tooltip title={title} placement={"top"} arrow>
      <Button {...rest}>{children}</Button>
    </Tooltip>
  );
};
export default TipButton;
