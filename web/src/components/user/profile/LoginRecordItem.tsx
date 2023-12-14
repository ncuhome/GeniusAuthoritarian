import { memo } from "react";
import { unix } from "dayjs";

import Ip from "./Ip";
import { TableCell, TableRow } from "@mui/material";

interface Props extends User.LoginRecord {}

export const LoginRecordItem = memo<Props>(
  ({ ...record }) => {
    return (
      <TableRow key={record.id}>
        <TableCell>
          {unix(record.createdAt).format("YYYY/MM/DD HH:mm")}
        </TableCell>
        <TableCell>{record.target}</TableCell>
        <TableCell>
          <Ip ip={record.ip} />
        </TableCell>
      </TableRow>
    );
  },
  (prev, next) => prev.id === next.id,
);
export default LoginRecordItem;
