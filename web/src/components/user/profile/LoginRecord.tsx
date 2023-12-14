import { FC } from "react";

import LoginRecordItem from "./LoginRecordItem";
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
            <TableCell>登录时间</TableCell>
            <TableCell>应用</TableCell>
            <TableCell>地址</TableCell>
            <TableCell>设备</TableCell>
          </TableRow>
        </TableHead>
        <TableBody>
          {records.map((record) => (
            <LoginRecordItem key={record.id} {...record} />
          ))}
        </TableBody>
      </Table>
    </Box>
  );
};
export default LoginRecord;
