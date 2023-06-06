import { FC, useMemo } from "react";
import { Tooltip, Typography } from "@mui/material";

interface Props {
  ip: string;
}

export const Ip: FC<Props> = ({ ip }) => {
  const highLight = useMemo(() => ip.indexOf("10.") !== 0, [ip]);

  return highLight ? (
    <Tooltip title={"非校园网登录"} placement={"top"} arrow>
      <Typography variant={"body2"} component={"span"} color={"warning.light"}>
        {ip}
      </Typography>
    </Tooltip>
  ) : (
    <span>{ip}</span>
  );
};
export default Ip
