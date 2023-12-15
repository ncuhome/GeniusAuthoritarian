import { FC } from "react";

import OnlineDeviceItem from "./OnlineDeviceItem";
import {
  Box,
  Table,
  TableBody,
  TableCell,
  TableHead,
  TableRow,
} from "@mui/material";

interface Props {
  records: User.LoginRecordOnline[];
}

export const OnlineDevice: FC<Props> = ({ records }) => {
  return (
    <Box
      sx={{
        marginTop: "0.5rem",
        width: "100%",
        whiteSpace: "nowrap",
      }}
    >
      <Table>
        <TableHead>
          <TableRow>
            <TableCell>令牌有效</TableCell>
            <TableCell>应用</TableCell>
            <TableCell>设备</TableCell>
            <TableCell>操作</TableCell>
          </TableRow>
        </TableHead>
        <TableBody>
          {records.map((record) => (
            <OnlineDeviceItem key={record.id} {...record} />
          ))}
        </TableBody>
      </Table>
    </Box>
  );
};
export default OnlineDevice;
