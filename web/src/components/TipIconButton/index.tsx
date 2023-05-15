import { FC } from "react";

import { IconButton, IconButtonProps, Tooltip } from "@mui/material";

interface Props {
  title: string;
}

export const TipIconButton: FC<IconButtonProps & Props> = ({
  title,
  children,
  ...rest
}) => {
  return (
    <Tooltip title={title} placement={"top"} arrow>
      <IconButton {...rest}>{children}</IconButton>
    </Tooltip>
  );
};
export default TipIconButton;
