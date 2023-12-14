import { FC } from "react";
import { unix } from "dayjs";

import Ip from "./Ip";
import {
  Box,
  Table,
  TableBody,
  TableCell,
  TableHead,
  TableRow,
} from "@mui/material";

interface Props {
  records: User.LoginRecord[];
}

export const LoginRecord: FC<Props> = ({ records }) => {
  return (
    <Box
      sx={{
        marginTop: "0.5rem",
        width: "100%",
        overflowY: "auto",
        whiteSpace: "nowrap",
      }}
    >
      <Table>
        <TableHead>
          <TableRow>
            <TableCell>时间</TableCell>
            <TableCell>站点</TableCell>
            <TableCell>IP 地址</TableCell>
          </TableRow>
        </TableHead>
        <TableBody>
          {records.map((record) => (
            <TableRow key={record.id}>
              <TableCell>
                {unix(record.createdAt).format("YYYY/MM/DD HH:mm")}
              </TableCell>
              <TableCell>{record.target}</TableCell>
              <TableCell>
                <Ip ip={record.ip} />
              </TableCell>
            </TableRow>
          ))}
        </TableBody>
      </Table>
    </Box>
  );
};
export default LoginRecord;
