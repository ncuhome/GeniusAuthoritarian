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
            <TableCell>时间</TableCell>
            <TableCell>站点</TableCell>
            <TableCell>IP 地址</TableCell>
          </TableRow>
        </TableHead>
        <TableBody>
          {records.map((record) => (
            <LoginRecordItem {...record} />
          ))}
        </TableBody>
      </Table>
    </Box>
  );
};
export default LoginRecord;
