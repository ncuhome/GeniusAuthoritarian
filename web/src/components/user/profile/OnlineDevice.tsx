import { FC, useCallback } from "react";

import OnlineDeviceItem from "./OnlineDeviceItem";
import { Flipper } from "react-flip-toolkit";
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
        overflowY: "auto",
        whiteSpace: "nowrap",
      }}
    >
      <Flipper flipKey={records}>
        <Table>
          <TableHead>
            <TableRow>
              <TableCell>上线时间</TableCell>
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
      </Flipper>
    </Box>
  );
};
export default OnlineDevice;
