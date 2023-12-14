import { memo } from "react";
import { unix } from "dayjs";

import Ip from "./Ip";
import UserAgent from "./UserAgent";
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
        <TableCell>
          {record.useragent ? (
            <UserAgent useragent={record.useragent} />
          ) : undefined}
        </TableCell>
      </TableRow>
    );
  },
  () => true,
);
export default LoginRecordItem;
