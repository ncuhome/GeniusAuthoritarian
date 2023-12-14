import { FC, useMemo } from "react";
import { Typography } from "@mui/material";

interface Props {
  ip: string;
}

export const Ip: FC<Props> = ({ ip }) => {
  const notSchoolNet = useMemo(() => ip.indexOf("10.") !== 0, [ip]);

  return notSchoolNet ? (
    <Typography variant={"body2"} component={"span"} color={"warning.light"}>
      {ip}
    </Typography>
  ) : (
    <span>校园网</span>
  );
};
export default Ip;
